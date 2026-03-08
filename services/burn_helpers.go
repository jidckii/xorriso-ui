package services

import (
	"time"

	"xorriso-ui/pkg/models"
)

func (s *BurnService) updateState(jobID string, state models.BurnState) {
	s.mu.Lock()
	if s.currentJob != nil && s.currentJob.ID == jobID {
		s.currentJob.State = state
	}
	s.mu.Unlock()

	s.emitEvent(models.EventBurnStateChanged, string(state))
}

// updatePhase updates the current job's progress phase
func (s *BurnService) updatePhase(jobID string, phase string) {
	s.mu.Lock()
	if s.currentJob != nil && s.currentJob.ID == jobID {
		s.currentJob.Progress.Phase = phase
	}
	s.mu.Unlock()
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

	if state == models.BurnStateDone {
		s.emitEvent(models.EventBurnComplete, result)
	} else if state == models.BurnStateError {
		s.emitEvent(models.EventBurnError, errMsg)
	}
}

// emitLog sends a single log message via Wails event
func (s *BurnService) emitLog(msg string) {
	s.emitEvent(models.EventBurnLogLine, msg)
}

// emitLogLines отправляет информационные строки через событие лога
func (s *BurnService) emitLogLines(lines []string) {
	for _, line := range lines {
		s.emitEvent(models.EventBurnLogLine, line)
	}
}
