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
	if s.currentJob != nil && s.currentJob.State == models.BurnStateWriting {
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

func (s *BurnService) runBurn(ctx context.Context, project *models.Project, devicePath string, opts models.BurnOptions, jobID string) {
	s.updateState(jobID, models.BurnStateWriting)

	// Build xorriso command
	cmd := xorriso.NewCommand()
	cmd.Device(devicePath)

	// ISO options
	if project.VolumeID != "" {
		cmd.VolumeID(project.VolumeID)
	}
	cmd.RockRidge(project.ISOOptions.RockRidge)
	cmd.Joliet(project.ISOOptions.Joliet)
	if project.ISOOptions.MD5 {
		cmd.MD5("on")
	}
	if project.ISOOptions.BackupMode {
		cmd.ForBackup()
	}

	// Add files
	for _, entry := range project.Entries {
		cmd.Map(entry.SourcePath, entry.DestPath)
	}

	// Burn options
	if opts.Speed != "" && opts.Speed != "auto" {
		cmd.WriteSpeed(opts.Speed)
	}
	cmd.Dummy(opts.DummyMode)
	cmd.Close(opts.CloseDisc)
	cmd.StreamRecording(opts.StreamRecording)

	cmd.Commit()

	if opts.Eject {
		cmd.Eject("all")
	}

	// Execute with progress
	result, err := s.executor.RunWithProgress(ctx, func(p xorriso.Progress) {
		progress := models.BurnProgress{
			Phase:   p.Phase,
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
	}, cmd.Build()...)

	if err != nil {
		s.finishJob(jobID, models.BurnStateError, nil, err.Error())
		return
	}

	if result.ExitCode != 0 {
		errMsg := "xorriso exited with code " + fmt.Sprintf("%d", result.ExitCode)
		if len(result.InfoLines) > 0 {
			errMsg = result.InfoLines[len(result.InfoLines)-1]
		}
		s.finishJob(jobID, models.BurnStateError, nil, errMsg)
		return
	}

	// Optional verification
	if opts.Verify {
		s.updateState(jobID, models.BurnStateVerifying)
		// TODO: implement verification using -check_media
	}

	s.finishJob(jobID, models.BurnStateDone, &models.BurnResult{
		Success: true,
	}, "")
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
