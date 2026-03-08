package services

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"xorriso-ui/pkg/xorriso"
)

// mockRunner реализует xorriso.Runner для тестов
type mockRunner struct {
	RunFn             func(ctx context.Context, args ...string) (*xorriso.CmdResult, error)
	RunWithProgressFn func(ctx context.Context, progressFn func(xorriso.Progress), args ...string) (*xorriso.CmdResult, error)
	VersionFn         func(ctx context.Context) (string, error)
}

func (m *mockRunner) Run(ctx context.Context, args ...string) (*xorriso.CmdResult, error) {
	if m.RunFn != nil {
		return m.RunFn(ctx, args...)
	}
	return &xorriso.CmdResult{}, nil
}

func (m *mockRunner) RunWithProgress(ctx context.Context, progressFn func(xorriso.Progress), args ...string) (*xorriso.CmdResult, error) {
	if m.RunWithProgressFn != nil {
		return m.RunWithProgressFn(ctx, progressFn, args...)
	}
	return &xorriso.CmdResult{}, nil
}

func (m *mockRunner) Version(ctx context.Context) (string, error) {
	if m.VersionFn != nil {
		return m.VersionFn(ctx)
	}
	return "xorriso 1.5.6", nil
}

// --- Тесты ---

func TestDiscoverDrivesFromProc(t *testing.T) {
	// Создаём временную структуру /proc и /sys
	tmpDir := t.TempDir()
	procPath := filepath.Join(tmpDir, "cdrom_info")
	sysPath := filepath.Join(tmpDir, "sys_block")

	// Содержимое /proc/sys/dev/cdrom/info
	procContent := `CD-ROM information, Id: cdrom.c 3.20 2003/12/17

drive name:		sr0
drive speed:		48
drive # of slots:	1
Can close tray:		1
Can lock tray:		1
Can change speed:	1
Can read multisession:	1
`
	if err := os.WriteFile(procPath, []byte(procContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Создаём /sys/block/sr0/device/{vendor,model,rev}
	deviceDir := filepath.Join(sysPath, "sr0", "device")
	if err := os.MkdirAll(deviceDir, 0755); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(deviceDir, "vendor"), []byte("ASUS    \n"), 0644)
	os.WriteFile(filepath.Join(deviceDir, "model"), []byte("BW-16D1HT\n"), 0644)
	os.WriteFile(filepath.Join(deviceDir, "rev"), []byte("3.02\n"), 0644)

	devices, err := discoverDrivesFromProc(procPath, sysPath)
	if err != nil {
		t.Fatalf("discoverDrivesFromProc: %v", err)
	}

	if len(devices) != 1 {
		t.Fatalf("expected 1 device, got %d", len(devices))
	}

	dev := devices[0]
	if dev.Path != "/dev/sr0" {
		t.Errorf("Path = %q, want /dev/sr0", dev.Path)
	}
	if dev.Vendor != "ASUS" {
		t.Errorf("Vendor = %q, want ASUS", dev.Vendor)
	}
	if dev.Model != "BW-16D1HT" {
		t.Errorf("Model = %q, want BW-16D1HT", dev.Model)
	}
	if dev.Revision != "3.02" {
		t.Errorf("Revision = %q, want 3.02", dev.Revision)
	}
	if dev.DriveSpeed != 48 {
		t.Errorf("DriveSpeed = %d, want 48", dev.DriveSpeed)
	}
	if !dev.CanCloseTray {
		t.Error("CanCloseTray should be true")
	}
	if !dev.CanLockTray {
		t.Error("CanLockTray should be true")
	}
}

func TestDiscoverDrivesFromProc_NoFile(t *testing.T) {
	_, err := discoverDrivesFromProc("/nonexistent/path/cdrom_info", "/nonexistent/sys")
	if err == nil {
		t.Error("expected error for nonexistent proc path, got nil")
	}
}

func TestExtractAfterColon(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"normal", "Media current: DVD+R", "DVD+R"},
		{"with spaces", "Volume Id    : MY_DISC", "MY_DISC"},
		{"empty value", "Key:", ""},
		{"no colon", "no colon here", ""},
		{"multiple colons", "Time: 2024:01:01", "2024:01:01"},
		{"only colon", ":", ""},
		{"colon at end", "test:", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractAfterColon(tt.input)
			if got != tt.expected {
				t.Errorf("extractAfterColon(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestExtractMediaProduct(t *testing.T) {
	input := "Media product: 97m26s66f/79m59s71f , CMC Magnetics Corporation"
	got := extractMediaProduct(input)
	if got != "CMC Magnetics Corporation" {
		t.Errorf("extractMediaProduct() = %q, want %q", got, "CMC Magnetics Corporation")
	}
}

func TestExtractMediaProduct_NoComma(t *testing.T) {
	input := "Media product: Unknown Media"
	got := extractMediaProduct(input)
	if got != "Unknown Media" {
		t.Errorf("extractMediaProduct() = %q, want %q", got, "Unknown Media")
	}
}

func TestGetDriveProfiles(t *testing.T) {
	runner := &mockRunner{
		RunFn: func(ctx context.Context, args ...string) (*xorriso.CmdResult, error) {
			return &xorriso.CmdResult{
				ResultLines: []string{
					"Profile      : 0x0009 (CD-R)",
					"Profile      : 0x000A (CD-RW)",
					"Profile      : 0x0011 (DVD-R sequential recording)",
					"Profile      : 0x001B (DVD+R) (current)",
					"Profile      : 0x0041 (BD-R sequential recording)",
					"Profile      : 0x0043 (BD-RE)",
				},
			}, nil
		},
	}

	svc := NewDeviceService(runner)
	svc.emitEvent = func(name string, data ...any) {} // no-op для тестов

	profiles, err := svc.GetDriveProfiles("/dev/sr0")
	if err != nil {
		t.Fatalf("GetDriveProfiles: %v", err)
	}

	if len(profiles) == 0 {
		t.Fatal("expected non-empty profiles list")
	}

	// Проверяем, что есть хотя бы один профиль с current=true
	var hasCurrent bool
	for _, p := range profiles {
		if p.Current {
			hasCurrent = true
			break
		}
	}
	if !hasCurrent {
		t.Error("expected at least one profile with Current=true")
	}
}

func TestGetSpeeds(t *testing.T) {
	runner := &mockRunner{
		RunFn: func(ctx context.Context, args ...string) (*xorriso.CmdResult, error) {
			return &xorriso.CmdResult{
				ResultLines: []string{
					"Write speed  :   1385kB/s  (CD  1x)",
					"Write speed  :   2770kB/s  (CD  2x)",
					"Write speed  :   5540kB/s  (CD  4x)",
					"Write speed  :  11080kB/s  (CD  8x)",
				},
			}, nil
		},
	}

	svc := NewDeviceService(runner)
	svc.emitEvent = func(name string, data ...any) {}

	speeds, err := svc.GetSpeeds("/dev/sr0")
	if err != nil {
		t.Fatalf("GetSpeeds: %v", err)
	}

	if len(speeds) == 0 {
		t.Fatal("expected non-empty speeds list")
	}
}

func TestProfileCaching(t *testing.T) {
	// Создаём temp proc/sys
	tmpDir := t.TempDir()
	procPath := filepath.Join(tmpDir, "cdrom_info")
	sysPath := filepath.Join(tmpDir, "sys_block")

	procContent := `CD-ROM information, Id: cdrom.c 3.20 2003/12/17

drive name:		sr0
drive speed:		48
drive # of slots:	1
Can close tray:		1
Can lock tray:		1
`
	os.WriteFile(procPath, []byte(procContent), 0644)

	deviceDir := filepath.Join(sysPath, "sr0", "device")
	os.MkdirAll(deviceDir, 0755)
	os.WriteFile(filepath.Join(deviceDir, "vendor"), []byte("ASUS\n"), 0644)
	os.WriteFile(filepath.Join(deviceDir, "model"), []byte("BW-16D1HT\n"), 0644)
	os.WriteFile(filepath.Join(deviceDir, "rev"), []byte("3.02\n"), 0644)

	callCount := 0
	runner := &mockRunner{
		RunFn: func(ctx context.Context, args ...string) (*xorriso.CmdResult, error) {
			callCount++
			return &xorriso.CmdResult{
				ResultLines: []string{
					"Profile      : 0x0009 (CD-R)",
					"Profile      : 0x001B (DVD+R) (current)",
				},
			}, nil
		},
	}

	svc := NewDeviceService(runner)
	svc.emitEvent = func(name string, data ...any) {}
	svc.procCdromInfoPath = procPath
	svc.sysBlockPath = sysPath

	// Первый вызов — должен обратиться к xorriso
	_, err := svc.ListDevices()
	if err != nil {
		t.Fatalf("first ListDevices: %v", err)
	}
	if callCount != 1 {
		t.Fatalf("expected 1 xorriso call after first ListDevices, got %d", callCount)
	}

	// Второй вызов — должен использовать кэш
	_, err = svc.ListDevices()
	if err != nil {
		t.Fatalf("second ListDevices: %v", err)
	}
	if callCount != 1 {
		t.Fatalf("expected still 1 xorriso call after second ListDevices, got %d", callCount)
	}
}

func TestGetMediaInfo(t *testing.T) {
	runner := &mockRunner{
		RunFn: func(ctx context.Context, args ...string) (*xorriso.CmdResult, error) {
			return &xorriso.CmdResult{
				ResultLines: []string{
					"Media current: DVD+R",
					"Media status : is written , is appendable",
					"Media erasable: is not erasable",
					"Media product: 97m26s66f/79m59s71f , Verbatim",
					"Volume Id    : TEST_DISC",
					"Volume Set Id: ",
					"Publisher Id : ",
					"Preparer Id  : XORRISO",
					"App Id       : XORRISO",
					"System Id    : LINUX",
					"Creation Time: 2024-01-15 10:00:00",
					"Modif. Time  : 2024-01-15 10:00:00",
					"Media summary: 1 session, 100000 blocks, 195.3m",
				},
			}, nil
		},
	}

	svc := NewDeviceService(runner)
	svc.emitEvent = func(name string, data ...any) {}

	info, err := svc.GetMediaInfo("/dev/sr0")
	if err != nil {
		t.Fatalf("GetMediaInfo: %v", err)
	}

	if info.DevicePath != "/dev/sr0" {
		t.Errorf("DevicePath = %q, want /dev/sr0", info.DevicePath)
	}
	if info.MediaType != "DVD+R" {
		t.Errorf("MediaType = %q, want DVD+R", info.MediaType)
	}
	if info.MediaStatus != "is written , is appendable" {
		t.Errorf("MediaStatus = %q, unexpected", info.MediaStatus)
	}
	if info.Erasable {
		t.Error("Erasable should be false for DVD+R")
	}
	if info.MediaProduct != "Verbatim" {
		t.Errorf("MediaProduct = %q, want Verbatim", info.MediaProduct)
	}
	if info.VolumeID != "TEST_DISC" {
		t.Errorf("VolumeID = %q, want TEST_DISC", info.VolumeID)
	}
	if info.PreparerID != "XORRISO" {
		t.Errorf("PreparerID = %q, want XORRISO", info.PreparerID)
	}
	if info.SystemID != "LINUX" {
		t.Errorf("SystemID = %q, want LINUX", info.SystemID)
	}
}
