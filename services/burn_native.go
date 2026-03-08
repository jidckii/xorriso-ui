package services

import (
	"context"
	"fmt"
	"time"

	"xorriso-ui/pkg/models"
	"xorriso-ui/pkg/xorriso"
)

func (s *BurnService) runBurn(ctx context.Context, project *models.Project, devicePath string, opts models.BurnOptions, jobID string) {
	startTime := time.Now()

	if len(project.Entries) == 0 {
		s.finishJob(jobID, models.BurnStateError, nil, "project has no entries")
		return
	}

	s.updateState(jobID, models.BurnStateWriting)

	// Формируем команду xorriso
	cmd := xorriso.NewCommand()
	cmd.Device(devicePath)
	cmd.AbortOn("FAILURE")

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

		s.emitEvent(models.EventBurnProgress, progress)
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
		verifyCmd.AbortOn("FAILURE")
		if project.ISOOptions.MD5 {
			verifyCmd.MD5("on")
			verifyCmd.CheckMD5Recursive("FAILURE", "/")
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

			s.emitEvent(models.EventBurnProgress, progress)
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
		ejectCtx, ejectCancel := context.WithTimeout(context.Background(), ejectTimeout)
		defer ejectCancel()
		if _, err := s.executor.Run(ejectCtx, ejectCmd.Build()...); err != nil {
			s.emitLog(fmt.Sprintf("eject failed: %s", err))
		}
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

func (s *BurnService) runCreateISO(ctx context.Context, project *models.Project, outputPath string, jobID string) {
	startTime := time.Now()

	if len(project.Entries) == 0 {
		s.finishJob(jobID, models.BurnStateError, nil, "project has no entries")
		return
	}

	s.updateState(jobID, models.BurnStateCreatingISO)

	cmd := xorriso.NewCommand()
	cmd.StdioOutDevice(outputPath)
	cmd.AbortOn("FAILURE")

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

		s.emitEvent(models.EventBurnProgress, progress)
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
