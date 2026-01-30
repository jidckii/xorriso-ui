package services

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"xorriso-ui/pkg/models"
	"xorriso-ui/pkg/xorriso"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type DeviceService struct {
	executor *xorriso.Executor
	mu       sync.RWMutex
	devices  []models.Device
	stopCh   chan struct{}
}

func NewDeviceService(executor *xorriso.Executor) *DeviceService {
	return &DeviceService{
		executor: executor,
		stopCh:   make(chan struct{}),
	}
}

// ServiceName returns the name of this service
func (s *DeviceService) ServiceName() string {
	return "DeviceService"
}

// ServiceStartup is called by Wails when the app starts
func (s *DeviceService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	go s.pollDevices()
	return nil
}

// ServiceShutdown is called by Wails when the app stops
func (s *DeviceService) ServiceShutdown() error {
	close(s.stopCh)
	return nil
}

// ListDevices detects optical drives via /proc/sys/dev/cdrom/info and /sys/block/
func (s *DeviceService) ListDevices() ([]models.Device, error) {
	devices, err := discoverDrivesFromProc()
	if err != nil {
		return nil, fmt.Errorf("failed to discover drives: %w", err)
	}

	// Enrich each device with supported profiles from xorriso
	for i := range devices {
		profiles, err := s.GetDriveProfiles(devices[i].Path)
		if err == nil {
			devices[i].Profiles = profiles
		}
	}

	s.mu.Lock()
	s.devices = devices
	s.mu.Unlock()

	return devices, nil
}

// discoverDrivesFromProc reads /proc/sys/dev/cdrom/info to get drive names,
// then reads /sys/block/<name>/device/ for vendor and model info.
func discoverDrivesFromProc() ([]models.Device, error) {
	f, err := os.Open("/proc/sys/dev/cdrom/info")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var driveNames []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "drive name:") {
			parts := strings.Fields(line)
			// "drive name:\tsr0 sr1 ..."
			if len(parts) > 2 {
				driveNames = append(driveNames, parts[2:]...)
			}
		}
	}

	var devices []models.Device
	for _, name := range driveNames {
		dev := models.Device{
			Path:     fmt.Sprintf("/dev/%s", name),
			LinkPath: resolveSymlink(fmt.Sprintf("/dev/%s", name)),
		}

		sysPath := filepath.Join("/sys/block", name, "device")
		dev.Vendor = readSysFile(filepath.Join(sysPath, "vendor"))
		dev.Model = readSysFile(filepath.Join(sysPath, "model"))
		dev.Revision = readSysFile(filepath.Join(sysPath, "rev"))

		devices = append(devices, dev)
	}

	return devices, nil
}

func readSysFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func resolveSymlink(path string) string {
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return path
	}
	return resolved
}

// GetDriveProfiles returns all media profiles supported by the drive (via xorriso)
func (s *DeviceService) GetDriveProfiles(devicePath string) ([]models.MediaProfile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := s.executor.Run(ctx,
		"-outdev", devicePath,
		"-list_profiles", "all",
	)
	if err != nil {
		return nil, err
	}

	return xorriso.ParseProfiles(result.ResultLines), nil
}

// GetMediaInfo returns information about the media in the specified drive
func (s *DeviceService) GetMediaInfo(devicePath string) (*models.MediaInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := s.executor.Run(ctx,
		"-dev", devicePath,
		"-toc",
		"-tell_media_space",
	)
	if err != nil {
		return nil, err
	}

	info := &models.MediaInfo{
		DevicePath: devicePath,
	}

	// Parse media space
	freeBlocks, _ := xorriso.ParseMediaSpace(result.ResultLines)
	info.FreeSpace = freeBlocks * 2048 // blocks are 2048 bytes

	// Parse info lines for media type and status
	for _, line := range result.InfoLines {
		if strings.Contains(line, "Media current:") {
			info.MediaType = extractAfterColon(line)
		}
		if strings.Contains(line, "Media status :") {
			info.MediaStatus = extractAfterColon(line)
		}
		if strings.Contains(line, "Media erasable") {
			info.Erasable = strings.Contains(line, "is erasable")
		}
	}

	return info, nil
}

// GetSpeeds returns available write speeds for the device
func (s *DeviceService) GetSpeeds(devicePath string) ([]models.SpeedDescriptor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := s.executor.Run(ctx,
		"-dev", devicePath,
		"-list_speeds",
	)
	if err != nil {
		return nil, err
	}

	return xorriso.ParseSpeeds(result.ResultLines), nil
}

// EjectDisc ejects the disc from the drive
func (s *DeviceService) EjectDisc(devicePath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := s.executor.Run(ctx,
		"-dev", devicePath,
		"-eject", "all",
	)
	return err
}

func (s *DeviceService) pollDevices() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			devices, err := s.ListDevices()
			if err == nil {
				if app := application.Get(); app != nil {
					app.Event.Emit(models.EventDeviceListUpdated, devices)
				}
			}
		}
	}
}

func extractAfterColon(line string) string {
	idx := strings.Index(line, ":")
	if idx < 0 || idx+1 >= len(line) {
		return ""
	}
	return strings.TrimSpace(line[idx+1:])
}
