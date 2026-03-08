package services

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"xorriso-ui/pkg/models"
	"xorriso-ui/pkg/xorriso"
)

// noopEmit заглушка для emitEvent
func noopEmit(name string, data ...any) {}

func TestCheckDiskSpace(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.iso")
	os.WriteFile(tmpFile, []byte("test"), 0644)

	svc := NewBurnService(&mockRunner{})
	svc.emitEvent = noopEmit

	// Запрашиваем 1 байт свободного места -- должно хватить
	ok, available, err := svc.CheckDiskSpace(tmpFile, 1)
	if err != nil {
		t.Fatalf("CheckDiskSpace: %v", err)
	}
	if !ok {
		t.Error("expected enough disk space for 1 byte")
	}
	if available <= 0 {
		t.Errorf("available = %d, expected > 0", available)
	}
}

func TestStartBurn_AlreadyInProgress(t *testing.T) {
	svc := NewBurnService(&mockRunner{})
	svc.emitEvent = noopEmit

	// Устанавливаем текущий job в состоянии writing
	svc.currentJob = &models.BurnJob{
		ID:        "existing-job",
		State:     models.BurnStateWriting,
		StartedAt: time.Now(),
	}

	project := &models.Project{
		Name:     "Test",
		VolumeID: "TEST",
		Entries:  []models.FileEntry{},
	}
	opts := models.BurnOptions{}

	_, err := svc.StartBurn(project, "/dev/sr0", opts)
	if err == nil {
		t.Fatal("expected error when burn already in progress")
	}
	if err.Error() != "burn already in progress" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCancelBurn_NoJob(t *testing.T) {
	svc := NewBurnService(&mockRunner{})
	svc.emitEvent = noopEmit

	err := svc.CancelBurn("some-id")
	if err == nil {
		t.Fatal("expected error when no job exists")
	}
}

func TestCancelBurn_WrongID(t *testing.T) {
	svc := NewBurnService(&mockRunner{})
	svc.emitEvent = noopEmit

	svc.currentJob = &models.BurnJob{
		ID:    "correct-id",
		State: models.BurnStateWriting,
	}

	err := svc.CancelBurn("wrong-id")
	if err == nil {
		t.Fatal("expected error for wrong job ID")
	}
	if err.Error() != "no matching burn job found" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestGetJobStatus_NoJob(t *testing.T) {
	svc := NewBurnService(&mockRunner{})
	svc.emitEvent = noopEmit

	_, err := svc.GetJobStatus("some-id")
	if err == nil {
		t.Fatal("expected error when no job exists")
	}
	if err.Error() != "job not found" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestGetJobStatus_Found(t *testing.T) {
	svc := NewBurnService(&mockRunner{})
	svc.emitEvent = noopEmit

	job := &models.BurnJob{
		ID:        "test-job-123",
		State:     models.BurnStatePending,
		StartedAt: time.Now(),
	}
	svc.currentJob = job

	got, err := svc.GetJobStatus("test-job-123")
	if err != nil {
		t.Fatalf("GetJobStatus: %v", err)
	}
	if got.ID != "test-job-123" {
		t.Errorf("ID = %q, want test-job-123", got.ID)
	}
	if got.State != models.BurnStatePending {
		t.Errorf("State = %q, want pending", got.State)
	}
}

func TestBuildISOCommand(t *testing.T) {
	runner := &mockRunner{}
	svc := NewBurnService(runner)
	svc.emitEvent = noopEmit

	project := &models.Project{
		VolumeID: "MY_VOLUME",
		Entries: []models.FileEntry{
			{SourcePath: "/home/user/file.txt", DestPath: "/file.txt"},
			{SourcePath: "/home/user/dir", DestPath: "/dir", IsDir: true},
		},
		ISOOptions: models.ISOOptions{
			ISOLevel:   3,
			RockRidge:  true,
			Joliet:     true,
			HFSPlus:    true,
			Zisofs:     true,
			MD5:        true,
			BackupMode: true,
		},
	}

	cmd := xorriso.NewCommand()
	svc.buildISOCommand(cmd, project)
	args := cmd.Build()

	// Проверяем наличие ожидаемых аргументов
	argStr := joinArgs(args)

	expectations := []string{
		"-volid", "MY_VOLUME",
		"-iso_level", "3",
		"-rockridge", "on",
		"-joliet", "on",
		"-hfsplus", "on",
		"-zisofs", "by_magic",
		"-md5", "on",
		"-for_backup",
		"-map", "/home/user/file.txt", "/file.txt",
		"-map", "/home/user/dir", "/dir",
	}

	for _, exp := range expectations {
		if !containsArg(args, exp) {
			t.Errorf("expected argument %q in command, got: %s", exp, argStr)
		}
	}
}

func TestBuildISOCommand_Minimal(t *testing.T) {
	svc := NewBurnService(&mockRunner{})
	svc.emitEvent = noopEmit

	project := &models.Project{
		Entries: []models.FileEntry{},
		ISOOptions: models.ISOOptions{
			RockRidge: false,
			Joliet:    false,
		},
	}

	cmd := xorriso.NewCommand()
	svc.buildISOCommand(cmd, project)
	args := cmd.Build()

	// Не должно быть -volid (VolumeID пустой)
	if containsArg(args, "-volid") {
		t.Error("should not have -volid when VolumeID is empty")
	}
	// Должно быть -rockridge off
	if !containsSequence(args, "-rockridge", "off") {
		t.Error("expected -rockridge off")
	}
	// Joliet=false — не должно быть -joliet вообще (xorriso default — off)
	if containsArg(args, "-joliet") {
		t.Error("should not have -joliet when Joliet is false (xorriso default is off)")
	}
}

func TestBuildISOCommand_UDF(t *testing.T) {
	svc := NewBurnService(&mockRunner{})
	svc.emitEvent = noopEmit

	project := &models.Project{
		Entries: []models.FileEntry{
			{SourcePath: "/tmp/data", DestPath: "/data"},
		},
		ISOOptions: models.ISOOptions{
			UDF: true,
		},
	}

	cmd := xorriso.NewCommand()
	svc.buildISOCommand(cmd, project)
	args := cmd.Build()

	if !containsSequence(args, "-udf", "on") {
		t.Errorf("expected -udf on in args, got: %s", joinArgs(args))
	}
}

func TestBuildISOCommand_UDFOff(t *testing.T) {
	svc := NewBurnService(&mockRunner{})
	svc.emitEvent = noopEmit

	project := &models.Project{
		Entries: []models.FileEntry{
			{SourcePath: "/tmp/data", DestPath: "/data"},
		},
		ISOOptions: models.ISOOptions{
			UDF: false,
		},
	}

	cmd := xorriso.NewCommand()
	svc.buildISOCommand(cmd, project)
	args := cmd.Build()

	// UDF=false — не должно быть -udf вообще (xorriso default — off)
	if containsArg(args, "-udf") {
		t.Errorf("should not have -udf when UDF is false (xorriso default is off), got: %s", joinArgs(args))
	}
}

func TestGetBurnCommand_Full(t *testing.T) {
	svc := NewBurnService(&mockRunner{})
	svc.emitEvent = noopEmit

	project := &models.Project{
		VolumeID: "DISC_1",
		Entries: []models.FileEntry{
			{SourcePath: "/home/user/file.txt", DestPath: "/file.txt"},
			{SourcePath: "/home/user/docs", DestPath: "/docs", IsDir: true},
		},
		ISOOptions: models.ISOOptions{
			ISOLevel:  3,
			RockRidge: true,
			Joliet:    true,
		},
	}

	opts := models.BurnOptions{
		Speed:           "4x",
		DummyMode:       true,
		CloseDisc:       true,
		StreamRecording: true,
		Eject:           true,
		BurnMode:        "SAO",
		Padding:         300,
	}

	result, err := svc.GetBurnCommand(project, "/dev/sr0", opts)
	if err != nil {
		t.Fatalf("GetBurnCommand: %v", err)
	}

	// Должна начинаться с "xorriso"
	if result[:8] != "xorriso " {
		t.Errorf("expected command to start with 'xorriso ', got: %s", result[:20])
	}

	// Проверяем наличие ключевых фрагментов
	expectations := []string{
		"-dev /dev/sr0",
		"-volid DISC_1",
		"-iso_level 3",
		"-rockridge on",
		"-joliet on",
		"-speed 4x",
		"-write_type SAO",
		"-padding 300k",
		"-dummy on",
		"-close on",
		"-stream_recording on",
		"-map /home/user/file.txt /file.txt",
		"-map /home/user/docs /docs",
		"-commit",
		"-eject all",
	}

	for _, exp := range expectations {
		if !containsSubstring(result, exp) {
			t.Errorf("expected %q in command, got: %s", exp, result)
		}
	}
}

func TestGetBurnCommand_Minimal(t *testing.T) {
	svc := NewBurnService(&mockRunner{})
	svc.emitEvent = noopEmit

	project := &models.Project{
		Entries: []models.FileEntry{
			{SourcePath: "/tmp/data", DestPath: "/data"},
		},
		ISOOptions: models.ISOOptions{},
	}

	opts := models.BurnOptions{}

	result, err := svc.GetBurnCommand(project, "/dev/sr1", opts)
	if err != nil {
		t.Fatalf("GetBurnCommand: %v", err)
	}

	if !containsSubstring(result, "-dev /dev/sr1") {
		t.Error("expected -dev /dev/sr1")
	}
	if !containsSubstring(result, "-commit") {
		t.Error("expected -commit")
	}
	// Без eject
	if containsSubstring(result, "-eject") {
		t.Error("should not have -eject when opts.Eject is false")
	}
	// Speed auto — не должно быть -speed
	if containsSubstring(result, "-speed") {
		t.Error("should not have -speed when speed is empty")
	}
}

func TestGetBurnCommand_Multisession(t *testing.T) {
	svc := NewBurnService(&mockRunner{})
	svc.emitEvent = noopEmit

	project := &models.Project{
		Entries: []models.FileEntry{
			{SourcePath: "/tmp/data", DestPath: "/data"},
		},
	}

	// CloseDisc + Multisession — взаимоисключающие опции, должна быть ошибка
	opts := models.BurnOptions{
		Multisession: true,
		CloseDisc:    true,
	}

	_, err := svc.GetBurnCommand(project, "/dev/sr0", opts)
	if err == nil {
		t.Fatal("expected error for mutually exclusive closeDisc and multisession")
	}

	// Multisession без CloseDisc — должно работать
	opts2 := models.BurnOptions{
		Multisession: true,
		CloseDisc:    false,
	}

	result, err := svc.GetBurnCommand(project, "/dev/sr0", opts2)
	if err != nil {
		t.Fatalf("GetBurnCommand: %v", err)
	}

	// Multisession: close off (диск остаётся открытым)
	if !containsSubstring(result, "-close off") {
		t.Errorf("expected -close off for multisession, got: %s", result)
	}
}

func TestGetBurnCommand_NilProject(t *testing.T) {
	svc := NewBurnService(&mockRunner{})
	svc.emitEvent = noopEmit

	_, err := svc.GetBurnCommand(nil, "/dev/sr0", models.BurnOptions{})
	if err == nil {
		t.Fatal("expected error for nil project")
	}
}

func TestGetBurnCommand_EmptyDevice(t *testing.T) {
	svc := NewBurnService(&mockRunner{})
	svc.emitEvent = noopEmit

	project := &models.Project{
		Entries: []models.FileEntry{{SourcePath: "/tmp/f", DestPath: "/f"}},
	}

	_, err := svc.GetBurnCommand(project, "", models.BurnOptions{})
	if err == nil {
		t.Fatal("expected error for empty device path")
	}
}

func TestGetBurnCommand_NoEntries(t *testing.T) {
	svc := NewBurnService(&mockRunner{})
	svc.emitEvent = noopEmit

	project := &models.Project{
		Entries: []models.FileEntry{},
	}

	_, err := svc.GetBurnCommand(project, "/dev/sr0", models.BurnOptions{})
	if err == nil {
		t.Fatal("expected error for project with no entries")
	}
}

// --- Вспомогательные функции ---

func joinArgs(args []string) string {
	var result strings.Builder
	for i, a := range args {
		if i > 0 {
			result.WriteString(" ")
		}
		result.WriteString(a)
	}
	return result.String()
}

func containsArg(args []string, target string) bool {
	return slices.Contains(args, target)
}

func containsSubstring(s, substr string) bool {
	return strings.Contains(s, substr)
}

func containsSequence(args []string, a, b string) bool {
	for i := 0; i < len(args)-1; i++ {
		if args[i] == a && args[i+1] == b {
			return true
		}
	}
	return false
}
