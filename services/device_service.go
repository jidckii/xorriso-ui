package services

import (
	"context"
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
	// Start device polling in background
	go s.pollDevices()
	return nil
}

// ServiceShutdown is called by Wails when the app stops
func (s *DeviceService) ServiceShutdown() error {
	close(s.stopCh)
	return nil
}

// ListDevices returns all detected optical drives
func (s *DeviceService) ListDevices() ([]models.Device, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := s.executor.Run(ctx, "-device_links")
	if err != nil {
		return nil, err
	}

	devices := xorriso.ParseDevices(result.ResultLines)

	s.mu.Lock()
	s.devices = devices
	s.mu.Unlock()

	return devices, nil
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
		if contains(line, "Media current:") {
			info.MediaType = extractAfterColon(line)
		}
		if contains(line, "Media status :") {
			info.MediaStatus = extractAfterColon(line)
		}
		if contains(line, "Media erasable") {
			info.Erasable = contains(line, "is erasable")
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

func contains(s, substr string) bool {
	return len(s) >= len(substr) && containsStr(s, substr)
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func extractAfterColon(line string) string {
	for i, c := range line {
		if c == ':' && i+1 < len(line) {
			result := line[i+1:]
			// Trim spaces
			start := 0
			for start < len(result) && result[start] == ' ' {
				start++
			}
			end := len(result)
			for end > start && (result[end-1] == ' ' || result[end-1] == '\n') {
				end--
			}
			return result[start:end]
		}
	}
	return ""
}
