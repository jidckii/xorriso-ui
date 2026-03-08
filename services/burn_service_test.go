package services

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"xorriso-ui/pkg/mkisofs"
	"xorriso-ui/pkg/models"
	"xorriso-ui/pkg/xorriso"
)

// mockISOBuilder реализует mkisofs.ISOBuilder для тестов
type mockISOBuilder struct {
	BuildISOFn func(ctx context.Context, opts mkisofs.BuildOpts, progressFn mkisofs.ProgressFn) error
}

func (m *mockISOBuilder) BuildISO(ctx context.Context, opts mkisofs.BuildOpts, progressFn mkisofs.ProgressFn) error {
	if m.BuildISOFn != nil {
		return m.BuildISOFn(ctx, opts, progressFn)
	}
	return nil
}

// noopEmit заглушка для emitEvent
func noopEmit(name string, data ...interface{}) {}

func TestCheckDiskSpace(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.iso")
	os.WriteFile(tmpFile, []byte("test"), 0644)

	svc := NewBurnService(&mockRunner{}, nil)
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
	svc := NewBurnService(&mockRunner{}, nil)
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
	svc := NewBurnService(&mockRunner{}, nil)
	svc.emitEvent = noopEmit

	err := svc.CancelBurn("some-id")
	if err == nil {
		t.Fatal("expected error when no job exists")
	}
}

func TestCancelBurn_WrongID(t *testing.T) {
	svc := NewBurnService(&mockRunner{}, nil)
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
	svc := NewBurnService(&mockRunner{}, nil)
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
	svc := NewBurnService(&mockRunner{}, nil)
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

func TestMkisofsAvailable_True(t *testing.T) {
	svc := NewBurnService(&mockRunner{}, &mockISOBuilder{})
	if !svc.MkisofsAvailable() {
		t.Error("MkisofsAvailable should return true when mkisofsExecutor is set")
	}
}

func TestMkisofsAvailable_False(t *testing.T) {
	svc := NewBurnService(&mockRunner{}, nil)
	if svc.MkisofsAvailable() {
		t.Error("MkisofsAvailable should return false when mkisofsExecutor is nil")
	}
}

func TestBuildISOCommand(t *testing.T) {
	runner := &mockRunner{}
	svc := NewBurnService(runner, nil)
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
	svc := NewBurnService(&mockRunner{}, nil)
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
	// Должно быть -joliet off
	if !containsSequence(args, "-joliet", "off") {
		t.Error("expected -joliet off")
	}
}

// --- Вспомогательные функции ---

func joinArgs(args []string) string {
	result := ""
	for i, a := range args {
		if i > 0 {
			result += " "
		}
		result += a
	}
	return result
}

func containsArg(args []string, target string) bool {
	for _, a := range args {
		if a == target {
			return true
		}
	}
	return false
}

func containsSequence(args []string, a, b string) bool {
	for i := 0; i < len(args)-1; i++ {
		if args[i] == a && args[i+1] == b {
			return true
		}
	}
	return false
}
