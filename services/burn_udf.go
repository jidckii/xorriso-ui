package services

import (
	"context"
	"fmt"
	"os"
	"time"

	"xorriso-ui/pkg/mkisofs"
	"xorriso-ui/pkg/models"
	"xorriso-ui/pkg/xorriso"

	"github.com/wailsapp/wails/v3/pkg/application"
)

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

		s.emitEvent(models.EventBurnProgress, progress)
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

		s.emitEvent(models.EventBurnProgress, progress)
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
		ejectCtx, ejectCancel := context.WithTimeout(context.Background(), ejectTimeout)
		defer ejectCancel()
		if _, err := s.executor.Run(ejectCtx, ejectCmd.Build()...); err != nil {
			s.emitLog(fmt.Sprintf("eject failed: %s", err))
		}
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

		s.emitEvent(models.EventBurnProgress, progress)
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
