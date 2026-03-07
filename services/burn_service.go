package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"xorriso-ui/pkg/models"
	"xorriso-ui/pkg/xorriso"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type BurnService struct {
	executor   *xorriso.Executor
	mu         sync.Mutex
	currentJob *models.BurnJob
	cancelFn   context.CancelFunc
}

func NewBurnService(executor *xorriso.Executor) *BurnService {
	return &BurnService{executor: executor}
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
