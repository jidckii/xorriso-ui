package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"xorriso-ui/pkg/mkisofs"
	"xorriso-ui/pkg/models"
	"xorriso-ui/pkg/xorriso"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type BurnService struct {
	executor        *xorriso.Executor
	mkisofsExecutor *mkisofs.Executor
	mu              sync.Mutex
	currentJob      *models.BurnJob
	cancelFn        context.CancelFunc
}

func NewBurnService(executor *xorriso.Executor, mkisofsExecutor *mkisofs.Executor) *BurnService {
	return &BurnService{executor: executor, mkisofsExecutor: mkisofsExecutor}
}

// MkisofsAvailable returns true if mkisofs binary is available
func (s *BurnService) MkisofsAvailable() bool {
	return s.mkisofsExecutor != nil
}

// CheckDiskSpace checks if there's enough space on the filesystem containing path
func (s *BurnService) CheckDiskSpace(path string, requiredBytes int64) (bool, int64, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(filepath.Dir(path), &stat); err != nil {
		return false, 0, fmt.Errorf("failed to check disk space: %w", err)
	}
	available := int64(stat.Bavail) * int64(stat.Bsize)
	return available >= requiredBytes, available, nil
}

func (s *BurnService) ServiceName() string {
	return "BurnService"
}

func (s *BurnService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	return nil
}

func (s *BurnService) ServiceShutdown() error {
	if s.cancelFn != nil {
		s.cancelFn()
	}
	return nil
}

// StartBurn begins the disc burning process
func (s *BurnService) StartBurn(project *models.Project, devicePath string, opts models.BurnOptions) (string, error) {
	s.mu.Lock()
	if s.currentJob != nil && (s.currentJob.State == models.BurnStateWriting || s.currentJob.State == models.BurnStateVerifying) {
		s.mu.Unlock()
		return "", fmt.Errorf("burn already in progress")
	}

	jobID := uuid.New().String()
	s.currentJob = &models.BurnJob{
		ID:        jobID,
		State:     models.BurnStatePending,
		StartedAt: time.Now(),
	}
	s.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	s.cancelFn = cancel

	go s.runBurn(ctx, project, devicePath, opts, jobID)

	return jobID, nil
}

// buildISOCommand формирует общую часть команды xorriso для ISO-опций и файлов проекта
func (s *BurnService) buildISOCommand(cmd *xorriso.CommandBuilder, project *models.Project) {
	if project.VolumeID != "" {
		cmd.VolumeID(project.VolumeID)
	}

	// ISO level
	if project.ISOOptions.ISOLevel > 0 {
		cmd.ISOLevel(project.ISOOptions.ISOLevel)
	}

	cmd.RockRidge(project.ISOOptions.RockRidge)
	cmd.Joliet(project.ISOOptions.Joliet)

	if project.ISOOptions.HFSPlus {
		cmd.HFSPlus(true)
	}

	if project.ISOOptions.Zisofs {
		cmd.Zisofs(true)
	}

	if project.ISOOptions.MD5 {
		cmd.MD5("on")
	}

	if project.ISOOptions.BackupMode {
		cmd.ForBackup()
	}

	// Добавить файлы
	for _, entry := range project.Entries {
		cmd.Map(entry.SourcePath, entry.DestPath)
	}
}

func (s *BurnService) runBurn(ctx context.Context, project *models.Project, devicePath string, opts models.BurnOptions, jobID string) {
	startTime := time.Now()
	s.updateState(jobID, models.BurnStateWriting)

	// UDF mode: create ISO via mkisofs, then burn via xorriso cdrecord mode
	if project.ISOOptions.UDF {
		s.runBurnUDF(ctx, project, devicePath, opts, jobID, startTime)
		return
	}

	// Non-UDF: existing native mode logic
	// Формируем команду xorriso
	cmd := xorriso.NewCommand()
	cmd.Device(devicePath)

	s.buildISOCommand(cmd, project)

	// Опции записи
	if opts.Speed != "" && opts.Speed != "auto" {
		cmd.WriteSpeed(opts.Speed)
	}
	if opts.BurnMode != "" && opts.BurnMode != "auto" {
		cmd.WriteType(opts.BurnMode)
	}
	if opts.Padding > 0 {
		cmd.Padding(opts.Padding)
	}

	cmd.Dummy(opts.DummyMode)
	if opts.Multisession {
		cmd.Close(false) // диск остаётся открытым для дозаписи
	} else {
		cmd.Close(opts.CloseDisc)
	}
	cmd.StreamRecording(opts.StreamRecording)

	cmd.Commit()
	// Eject НЕ добавляем в основную команду — выполняем отдельно после верификации

	// Выполняем запись с отслеживанием прогресса
	var lastProgress models.BurnProgress
	result, err := s.executor.RunWithProgress(ctx, func(p xorriso.Progress) {
		progress := models.BurnProgress{
			Phase:        p.Phase,
			Percent:      p.Percent,
			Speed:        p.Speed,
			BytesWritten: p.BytesWritten,
			BytesTotal:   p.BytesTotal,
			ETA:          p.ETA,
			FIFOFill:     p.FIFOPercent,
		}
		lastProgress = progress

		s.mu.Lock()
		if s.currentJob != nil && s.currentJob.ID == jobID {
			s.currentJob.Progress = progress
		}
		s.mu.Unlock()

		if app := application.Get(); app != nil {
			app.Event.Emit(models.EventBurnProgress, progress)
		}
	}, cmd.Build()...)

	if err != nil {
		s.finishJob(jobID, models.BurnStateError, nil, err.Error())
		return
	}

	// Отправляем информационные строки в лог
	s.emitLogLines(result.InfoLines)

	if result.ExitCode != 0 {
		errMsg := "xorriso exited with code " + fmt.Sprintf("%d", result.ExitCode)
		if len(result.InfoLines) > 0 {
			errMsg = result.InfoLines[len(result.InfoLines)-1]
		}
		s.finishJob(jobID, models.BurnStateError, nil, errMsg)
		return
	}

	// Верификация
	var verifyErrors int
	var md5Match bool
	if opts.Verify {
		s.updateState(jobID, models.BurnStateVerifying)

		verifyCmd := xorriso.NewCommand()
		verifyCmd.InDevice(devicePath)
		if project.ISOOptions.MD5 {
			verifyCmd.MD5("on")
			verifyCmd.CheckMD5("FAILURE")
		}
		verifyCmd.CheckMedia(nil)

		verifyResult, verifyErr := s.executor.RunWithProgress(ctx, func(p xorriso.Progress) {
			progress := models.BurnProgress{
				Phase:   "verifying",
				Percent: p.Percent,
				Speed:   p.Speed,
			}

			s.mu.Lock()
			if s.currentJob != nil && s.currentJob.ID == jobID {
				s.currentJob.Progress = progress
			}
			s.mu.Unlock()

			if app := application.Get(); app != nil {
				app.Event.Emit(models.EventBurnProgress, progress)
			}
		}, verifyCmd.Build()...)

		if verifyErr != nil {
			s.finishJob(jobID, models.BurnStateError, nil, fmt.Sprintf("verification failed: %s", verifyErr))
			return
		}

		s.emitLogLines(verifyResult.InfoLines)

		readErrors, md5Mismatches := xorriso.ParseCheckMediaResult(verifyResult.ResultLines)
		verifyErrors = readErrors + md5Mismatches
		md5Match = md5Mismatches == 0

		if verifyResult.ExitCode != 0 {
			s.finishJob(jobID, models.BurnStateError, nil, fmt.Sprintf("verification reported errors (exit code %d)", verifyResult.ExitCode))
			return
		}
	}

	// Eject после всех операций
	if opts.Eject {
		ejectCmd := xorriso.NewCommand()
		ejectCmd.Device(devicePath)
		ejectCmd.Eject("all")
		ejectCtx, ejectCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer ejectCancel()
		_, _ = s.executor.Run(ejectCtx, ejectCmd.Build()...)
	}

	// Формируем итоговый результат
	duration := time.Since(startTime)
	bytesWritten := lastProgress.BytesWritten
	var avgSpeed string
	if duration.Seconds() > 0 && bytesWritten > 0 {
		mbPerSec := float64(bytesWritten) / 1024.0 / 1024.0 / duration.Seconds()
		avgSpeed = fmt.Sprintf("%.2f MB/s", mbPerSec)
	}

	burnResult := &models.BurnResult{
		Success:      true,
		BytesWritten: bytesWritten,
		Duration:     duration.String(),
		AverageSpeed: avgSpeed,
		MD5Match:     md5Match,
		VerifyErrors: verifyErrors,
	}

	s.finishJob(jobID, models.BurnStateDone, burnResult, "")
}

// CreateISO создаёт ISO-файл без записи на привод
func (s *BurnService) CreateISO(project *models.Project, outputPath string) (string, error) {
	s.mu.Lock()
	if s.currentJob != nil && (s.currentJob.State == models.BurnStateWriting ||
		s.currentJob.State == models.BurnStateCreatingISO) {
		s.mu.Unlock()
		return "", fmt.Errorf("operation already in progress")
	}

	jobID := uuid.New().String()
	s.currentJob = &models.BurnJob{
		ID:        jobID,
		State:     models.BurnStatePending,
		StartedAt: time.Now(),
	}
	s.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	s.cancelFn = cancel

	go s.runCreateISO(ctx, project, outputPath, jobID)

	return jobID, nil
}

func (s *BurnService) runCreateISO(ctx context.Context, project *models.Project, outputPath string, jobID string) {
	startTime := time.Now()
	s.updateState(jobID, models.BurnStateCreatingISO)

	// UDF mode: use mkisofs
	if project.ISOOptions.UDF {
		s.runCreateISOWithMkisofs(ctx, project, outputPath, jobID, startTime)
		return
	}

	// Non-UDF: existing xorriso native mode
	cmd := xorriso.NewCommand()
	cmd.StdioOutDevice(outputPath)

	s.buildISOCommand(cmd, project)

	cmd.Commit()

	var lastProgress models.BurnProgress
	result, err := s.executor.RunWithProgress(ctx, func(p xorriso.Progress) {
		progress := models.BurnProgress{
			Phase:        "creating_iso",
			Percent:      p.Percent,
			Speed:        p.Speed,
			BytesWritten: p.BytesWritten,
			BytesTotal:   p.BytesTotal,
			ETA:          p.ETA,
			FIFOFill:     p.FIFOPercent,
		}
		lastProgress = progress

		s.mu.Lock()
		if s.currentJob != nil && s.currentJob.ID == jobID {
			s.currentJob.Progress = progress
		}
		s.mu.Unlock()

		if app := application.Get(); app != nil {
			app.Event.Emit(models.EventBurnProgress, progress)
		}
	}, cmd.Build()...)

	if err != nil {
		s.finishJob(jobID, models.BurnStateError, nil, err.Error())
		return
	}

	s.emitLogLines(result.InfoLines)

	if result.ExitCode != 0 {
		errMsg := "xorriso exited with code " + fmt.Sprintf("%d", result.ExitCode)
		if len(result.InfoLines) > 0 {
			errMsg = result.InfoLines[len(result.InfoLines)-1]
		}
		s.finishJob(jobID, models.BurnStateError, nil, errMsg)
		return
	}

	duration := time.Since(startTime)
	bytesWritten := lastProgress.BytesWritten
	var avgSpeed string
	if duration.Seconds() > 0 && bytesWritten > 0 {
		mbPerSec := float64(bytesWritten) / 1024.0 / 1024.0 / duration.Seconds()
		avgSpeed = fmt.Sprintf("%.2f MB/s", mbPerSec)
	}

	burnResult := &models.BurnResult{
		Success:      true,
		BytesWritten: bytesWritten,
		Duration:     duration.String(),
		AverageSpeed: avgSpeed,
	}

	s.finishJob(jobID, models.BurnStateDone, burnResult, "")
}

func (s *BurnService) runBurnUDF(ctx context.Context, project *models.Project, devicePath string, opts models.BurnOptions, jobID string, startTime time.Time) {
	if s.mkisofsExecutor == nil {
		s.finishJob(jobID, models.BurnStateError, nil, "mkisofs not found: required for UDF support")
		return
	}

	// Create temp ISO
	tmpFile, err := os.CreateTemp("", "xorriso-ui-*.iso")
	if err != nil {
		s.finishJob(jobID, models.BurnStateError, nil, fmt.Sprintf("failed to create temp file: %s", err))
		return
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()

	// Check disk space
	totalSize := int64(0)
	for _, e := range project.Entries {
		totalSize += e.Size
	}
	requiredSize := totalSize + totalSize/20 // +5% overhead
	ok, available, err := s.CheckDiskSpace(tmpPath, requiredSize)
	if err != nil {
		os.Remove(tmpPath)
		s.finishJob(jobID, models.BurnStateError, nil, fmt.Sprintf("failed to check disk space: %s", err))
		return
	}
	if !ok {
		os.Remove(tmpPath)
		s.finishJob(jobID, models.BurnStateError, nil, fmt.Sprintf("not enough disk space: need %d MB, available %d MB", requiredSize/1024/1024, available/1024/1024))
		return
	}

	s.emitLog("Creating UDF ISO image via mkisofs...")

	// Phase 1: Create ISO via mkisofs
	s.updatePhase(jobID, "creating_iso")

	mkisofsOpts := mkisofs.BuildOpts{
		OutputPath: tmpPath,
		VolumeID:   project.VolumeID,
		UDF:        true,
		RockRidge:  project.ISOOptions.RockRidge,
		Joliet:     project.ISOOptions.Joliet,
		HFSPlus:    project.ISOOptions.HFSPlus,
		Zisofs:     project.ISOOptions.Zisofs,
		ISOLevel:   project.ISOOptions.ISOLevel,
		Files:      mkisofs.FileMappingsFromEntries(project.Entries),
	}

	err = s.mkisofsExecutor.BuildISO(ctx, mkisofsOpts, func(percent float64) {
		progress := models.BurnProgress{
			Phase:   "creating_iso",
			Percent: percent,
		}
		s.mu.Lock()
		if s.currentJob != nil && s.currentJob.ID == jobID {
			s.currentJob.Progress = progress
		}
		s.mu.Unlock()

		if app := application.Get(); app != nil {
			app.Event.Emit(models.EventBurnProgress, progress)
		}
	})

	if err != nil {
		os.Remove(tmpPath)
		s.finishJob(jobID, models.BurnStateError, nil, fmt.Sprintf("mkisofs failed: %s", err))
		return
	}

	s.emitLog("UDF ISO created, writing to disc...")

	// Phase 2: Burn ISO via xorriso cdrecord mode
	s.updatePhase(jobID, "writing")

	cmd := xorriso.NewCommand()
	cmd.CdrecordMode()
	cmd.CdrecordDev(devicePath)
	if opts.Speed != "" && opts.Speed != "auto" {
		cmd.CdrecordSpeed(opts.Speed)
	}
	cmd.Verbose()
	cmd.Arg(tmpPath)

	var lastProgress models.BurnProgress
	result, err := s.executor.RunWithProgress(ctx, func(p xorriso.Progress) {
		progress := models.BurnProgress{
			Phase:        "writing",
			Percent:      p.Percent,
			Speed:        p.Speed,
			BytesWritten: p.BytesWritten,
			BytesTotal:   p.BytesTotal,
			ETA:          p.ETA,
			FIFOFill:     p.FIFOPercent,
		}
		lastProgress = progress

		s.mu.Lock()
		if s.currentJob != nil && s.currentJob.ID == jobID {
			s.currentJob.Progress = progress
		}
		s.mu.Unlock()

		if app := application.Get(); app != nil {
			app.Event.Emit(models.EventBurnProgress, progress)
		}
	}, cmd.Build()...)

	if err != nil {
		s.finishJob(jobID, models.BurnStateError, nil, fmt.Sprintf("burn failed: %s", err))
		s.emitLog(fmt.Sprintf("Temporary ISO kept at: %s", tmpPath))
		return
	}

	s.emitLogLines(result.InfoLines)

	if result.ExitCode != 0 {
		errMsg := "xorriso exited with code " + fmt.Sprintf("%d", result.ExitCode)
		if len(result.InfoLines) > 0 {
			errMsg = result.InfoLines[len(result.InfoLines)-1]
		}
		s.finishJob(jobID, models.BurnStateError, nil, errMsg)
		s.emitLog(fmt.Sprintf("Temporary ISO kept at: %s", tmpPath))
		return
	}

	// Verification
	var verifyErrors int
	var md5Match bool
	if opts.Verify {
		s.updateState(jobID, models.BurnStateVerifying)

		verifyCmd := xorriso.NewCommand()
		verifyCmd.InDevice(devicePath)
		if project.ISOOptions.MD5 {
			verifyCmd.MD5("on")
			verifyCmd.CheckMD5("FAILURE")
		}
		verifyCmd.CheckMedia(nil)

		verifyResult, verifyErr := s.executor.RunWithProgress(ctx, func(p xorriso.Progress) {
			progress := models.BurnProgress{
				Phase:   "verifying",
				Percent: p.Percent,
				Speed:   p.Speed,
			}

			s.mu.Lock()
			if s.currentJob != nil && s.currentJob.ID == jobID {
				s.currentJob.Progress = progress
			}
			s.mu.Unlock()

			if app := application.Get(); app != nil {
				app.Event.Emit(models.EventBurnProgress, progress)
			}
		}, verifyCmd.Build()...)

		if verifyErr != nil {
			s.finishJob(jobID, models.BurnStateError, nil, fmt.Sprintf("verification failed: %s", verifyErr))
			return
		}

		s.emitLogLines(verifyResult.InfoLines)

		readErrors, md5Mismatches := xorriso.ParseCheckMediaResult(verifyResult.ResultLines)
		verifyErrors = readErrors + md5Mismatches
		md5Match = md5Mismatches == 0

		if verifyResult.ExitCode != 0 {
			s.finishJob(jobID, models.BurnStateError, nil, fmt.Sprintf("verification reported errors (exit code %d)", verifyResult.ExitCode))
			return
		}
	}

	// Cleanup temp ISO
	if opts.CleanupISO {
		os.Remove(tmpPath)
		s.emitLog("Temporary ISO removed")
	} else {
		s.emitLog(fmt.Sprintf("Temporary ISO kept at: %s", tmpPath))
	}

	// Eject
	if opts.Eject {
		ejectCmd := xorriso.NewCommand()
		ejectCmd.Device(devicePath)
		ejectCmd.Eject("all")
		ejectCtx, ejectCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer ejectCancel()
		_, _ = s.executor.Run(ejectCtx, ejectCmd.Build()...)
	}

	// Result
	duration := time.Since(startTime)
	bytesWritten := lastProgress.BytesWritten
	var avgSpeed string
	if duration.Seconds() > 0 && bytesWritten > 0 {
		mbPerSec := float64(bytesWritten) / 1024.0 / 1024.0 / duration.Seconds()
		avgSpeed = fmt.Sprintf("%.2f MB/s", mbPerSec)
	}

	burnResult := &models.BurnResult{
		Success:      true,
		BytesWritten: bytesWritten,
		Duration:     duration.String(),
		AverageSpeed: avgSpeed,
		MD5Match:     md5Match,
		VerifyErrors: verifyErrors,
	}

	s.finishJob(jobID, models.BurnStateDone, burnResult, "")
}

func (s *BurnService) runCreateISOWithMkisofs(ctx context.Context, project *models.Project, outputPath string, jobID string, startTime time.Time) {
	if s.mkisofsExecutor == nil {
		s.finishJob(jobID, models.BurnStateError, nil, "mkisofs not found: required for UDF support")
		return
	}

	// Check disk space
	totalSize := int64(0)
	for _, e := range project.Entries {
		totalSize += e.Size
	}
	requiredSize := totalSize + totalSize/20
	ok, available, err := s.CheckDiskSpace(outputPath, requiredSize)
	if err != nil {
		s.finishJob(jobID, models.BurnStateError, nil, fmt.Sprintf("failed to check disk space: %s", err))
		return
	}
	if !ok {
		s.finishJob(jobID, models.BurnStateError, nil, fmt.Sprintf("not enough disk space: need %d MB, available %d MB", requiredSize/1024/1024, available/1024/1024))
		return
	}

	mkisofsOpts := mkisofs.BuildOpts{
		OutputPath: outputPath,
		VolumeID:   project.VolumeID,
		UDF:        true,
		RockRidge:  project.ISOOptions.RockRidge,
		Joliet:     project.ISOOptions.Joliet,
		HFSPlus:    project.ISOOptions.HFSPlus,
		Zisofs:     project.ISOOptions.Zisofs,
		ISOLevel:   project.ISOOptions.ISOLevel,
		Files:      mkisofs.FileMappingsFromEntries(project.Entries),
	}

	err = s.mkisofsExecutor.BuildISO(ctx, mkisofsOpts, func(percent float64) {
		progress := models.BurnProgress{
			Phase:   "creating_iso",
			Percent: percent,
		}

		s.mu.Lock()
		if s.currentJob != nil && s.currentJob.ID == jobID {
			s.currentJob.Progress = progress
		}
		s.mu.Unlock()

		if app := application.Get(); app != nil {
			app.Event.Emit(models.EventBurnProgress, progress)
		}
	})

	if err != nil {
		s.finishJob(jobID, models.BurnStateError, nil, fmt.Sprintf("mkisofs failed: %s", err))
		return
	}

	// Get file size for result
	var bytesWritten int64
	if info, err := os.Stat(outputPath); err == nil {
		bytesWritten = info.Size()
	}

	duration := time.Since(startTime)
	var avgSpeed string
	if duration.Seconds() > 0 && bytesWritten > 0 {
		mbPerSec := float64(bytesWritten) / 1024.0 / 1024.0 / duration.Seconds()
		avgSpeed = fmt.Sprintf("%.2f MB/s", mbPerSec)
	}

	burnResult := &models.BurnResult{
		Success:      true,
		BytesWritten: bytesWritten,
		Duration:     duration.String(),
		AverageSpeed: avgSpeed,
	}

	s.finishJob(jobID, models.BurnStateDone, burnResult, "")
}

// emitLog sends a single log message via Wails event
func (s *BurnService) emitLog(msg string) {
	if app := application.Get(); app != nil {
		app.Event.Emit(models.EventBurnLogLine, msg)
	}
}

// updatePhase updates the current job's progress phase
func (s *BurnService) updatePhase(jobID string, phase string) {
	s.mu.Lock()
	if s.currentJob != nil && s.currentJob.ID == jobID {
		s.currentJob.Progress.Phase = phase
	}
	s.mu.Unlock()
}

// emitLogLines отправляет информационные строки через событие лога
func (s *BurnService) emitLogLines(lines []string) {
	app := application.Get()
	if app == nil {
		return
	}
	for _, line := range lines {
		app.Event.Emit(models.EventBurnLogLine, line)
	}
}

func (s *BurnService) updateState(jobID string, state models.BurnState) {
	s.mu.Lock()
	if s.currentJob != nil && s.currentJob.ID == jobID {
		s.currentJob.State = state
	}
	s.mu.Unlock()

	if app := application.Get(); app != nil {
		app.Event.Emit(models.EventBurnStateChanged, string(state))
	}
}

func (s *BurnService) finishJob(jobID string, state models.BurnState, result *models.BurnResult, errMsg string) {
	s.mu.Lock()
	if s.currentJob != nil && s.currentJob.ID == jobID {
		s.currentJob.State = state
		s.currentJob.Result = result
		s.currentJob.Error = errMsg
		s.currentJob.FinishedAt = time.Now()
	}
	s.mu.Unlock()

	if app := application.Get(); app != nil {
		if state == models.BurnStateDone {
			app.Event.Emit(models.EventBurnComplete, result)
		} else if state == models.BurnStateError {
			app.Event.Emit(models.EventBurnError, errMsg)
		}
	}
}

// CancelBurn cancels the current burn operation
func (s *BurnService) CancelBurn(jobID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.currentJob == nil || s.currentJob.ID != jobID {
		return fmt.Errorf("no matching burn job found")
	}

	if s.cancelFn != nil {
		s.cancelFn()
	}
	s.currentJob.State = models.BurnStateCancelled
	return nil
}

// GetJobStatus returns the current job status
func (s *BurnService) GetJobStatus(jobID string) (*models.BurnJob, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.currentJob == nil || s.currentJob.ID != jobID {
		return nil, fmt.Errorf("job not found")
	}
	return s.currentJob, nil
}

// BlankDisc blanks a rewritable disc
func (s *BurnService) BlankDisc(devicePath string, mode string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	_, err := s.executor.RunWithProgress(ctx, func(p xorriso.Progress) {
		if app := application.Get(); app != nil {
			app.Event.Emit(models.EventBurnProgress, models.BurnProgress{
				Phase:   "blanking",
				Percent: p.Percent,
			})
		}
	}, "-dev", devicePath, "-blank", mode)

	return err
}

// FormatDisc formats a disc (BD-RE, DVD-RAM, etc.)
func (s *BurnService) FormatDisc(devicePath string, mode string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	_, err := s.executor.RunWithProgress(ctx, func(p xorriso.Progress) {
		if app := application.Get(); app != nil {
			app.Event.Emit(models.EventBurnProgress, models.BurnProgress{
				Phase:   "formatting",
				Percent: p.Percent,
			})
		}
	}, "-dev", devicePath, "-format", mode)

	return err
}
