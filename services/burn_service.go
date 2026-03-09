package services

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"xorriso-ui/pkg/models"
	"xorriso-ui/pkg/xorriso"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v3/pkg/application"
)

const (
	ejectTimeout       = 30 * time.Second
	blankFormatTimeout = 30 * time.Minute
)

const (
	appApplicationID = "XORRISO-UI (C) Evgeniy Medvedev"
	appSystemID      = "LINUX"
)

type BurnService struct {
	executor   xorriso.Runner
	mu         sync.Mutex
	currentJob *models.BurnJob
	cancelFn   context.CancelFunc
	emitEvent  func(name string, data ...any)
}

func NewBurnService(executor xorriso.Runner) *BurnService {
	return &BurnService{
		executor:  executor,
		emitEvent: defaultEmitEvent,
	}
}

func defaultEmitEvent(name string, data ...any) {
	if app := application.Get(); app != nil {
		app.Event.Emit(name, data...)
	}
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

	if err := validateBurnOptions(opts); err != nil {
		return "", err
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.cancelFn = cancel

	go s.runBurn(ctx, project, devicePath, opts, jobID)

	return jobID, nil
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

// BlankDisc blanks a rewritable disc
func (s *BurnService) BlankDisc(devicePath string, mode string) error {
	ctx, cancel := context.WithTimeout(context.Background(), blankFormatTimeout)
	defer cancel()

	_, err := s.executor.RunWithProgress(ctx, func(p xorriso.Progress) {
		s.emitEvent(models.EventBurnProgress, models.BurnProgress{
			Phase:   "blanking",
			Percent: p.Percent,
		})
	}, "-dev", devicePath, "-blank", mode)

	return err
}

// FormatDisc formats a disc (BD-RE, DVD-RAM, etc.)
func (s *BurnService) FormatDisc(devicePath string, mode string) error {
	ctx, cancel := context.WithTimeout(context.Background(), blankFormatTimeout)
	defer cancel()

	_, err := s.executor.RunWithProgress(ctx, func(p xorriso.Progress) {
		s.emitEvent(models.EventBurnProgress, models.BurnProgress{
			Phase:   "formatting",
			Percent: p.Percent,
		})
	}, "-dev", devicePath, "-format", mode)

	return err
}

// validateBurnOptions проверяет корректность опций записи
func validateBurnOptions(opts models.BurnOptions) error {
	if opts.BurnMode != "" && opts.BurnMode != "auto" {
		validModes := map[string]bool{"DAO": true, "TAO": true, "SAO": true}
		if !validModes[opts.BurnMode] {
			return fmt.Errorf("invalid burn mode: %s (must be DAO, TAO, or SAO)", opts.BurnMode)
		}
	}
	if opts.Padding < 0 {
		return fmt.Errorf("padding cannot be negative: %d", opts.Padding)
	}
	if opts.CloseDisc && opts.Multisession {
		return fmt.Errorf("closeDisc and multisession are mutually exclusive")
	}
	return nil
}

// buildISOCommand формирует общую часть команды xorriso для ISO-опций и файлов проекта
func (s *BurnService) buildISOCommand(cmd *xorriso.CommandBuilder, project *models.Project) {
	if project.VolumeID != "" {
		cmd.VolumeID(project.VolumeID)
	}

	cmd.ApplicationID(appApplicationID)
	cmd.SystemID(appSystemID)
	if project.ISOOptions.PublisherID != "" {
		cmd.Publisher(project.ISOOptions.PublisherID)
	}

	// ISO level
	if project.ISOOptions.ISOLevel > 0 {
		cmd.ISOLevel(project.ISOOptions.ISOLevel)
	}

	cmd.RockRidge(project.ISOOptions.RockRidge)
	if project.ISOOptions.Joliet {
		cmd.Joliet(true)
	}
	if project.ISOOptions.UDF {
		cmd.UDF(true)
	}

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

// GetBurnCommand формирует полную строку команды xorriso для записи диска.
// Результат можно скопировать в буфер обмена и выполнить в терминале.
func (s *BurnService) GetBurnCommand(project *models.Project, devicePath string, opts models.BurnOptions) (string, error) {
	if project == nil {
		return "", fmt.Errorf("project is nil")
	}
	if devicePath == "" {
		return "", fmt.Errorf("device path is empty")
	}
	if len(project.Entries) == 0 {
		return "", fmt.Errorf("project has no entries")
	}

	if err := validateBurnOptions(opts); err != nil {
		return "", err
	}

	cmd := xorriso.NewCommand()
	cmd.Device(devicePath)

	// ISO-опции и файлы проекта
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
		cmd.Close(false)
	} else {
		cmd.Close(opts.CloseDisc)
	}
	cmd.StreamRecording(opts.StreamRecording)

	cmd.Commit()

	if opts.Eject {
		cmd.Eject("all")
	}

	args := cmd.Build()
	return "xorriso " + strings.Join(args, " "), nil
}
