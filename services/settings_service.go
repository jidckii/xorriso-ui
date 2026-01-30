package services

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"

	"xorriso-ui/pkg/models"
	"xorriso-ui/pkg/xorriso"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type AppSettings struct {
	XorrisoPath        string             `json:"xorrisoPath"`
	DefaultBurn        models.BurnOptions `json:"defaultBurn"`
	DefaultISO         models.ISOOptions  `json:"defaultIso"`
	BDXLSafeMode       bool               `json:"bdxlSafeMode"`
	AutoEject          bool               `json:"autoEject"`
	DevicePollInterval int                `json:"devicePollInterval"`
	Language           string             `json:"language"`
	Theme              string             `json:"theme"`
}

type SettingsService struct {
	executor *xorriso.Executor
}

func NewSettingsService(executor *xorriso.Executor) *SettingsService {
	return &SettingsService{executor: executor}
}

func (s *SettingsService) ServiceName() string {
	return "SettingsService"
}

func (s *SettingsService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	return nil
}

func (s *SettingsService) ServiceShutdown() error {
	return nil
}

func (s *SettingsService) settingsPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = os.TempDir()
	}
	return filepath.Join(configDir, "xorriso-ui", "settings.json")
}

// GetSettings loads settings from disk
func (s *SettingsService) GetSettings() (*AppSettings, error) {
	xorrisoPath := "xorriso"
	if p, err := exec.LookPath("xorriso"); err == nil {
		xorrisoPath = p
	}

	settings := &AppSettings{
		XorrisoPath: xorrisoPath,
		DefaultBurn: models.BurnOptions{
			Speed:    "auto",
			Verify:   true,
			Eject:    true,
			BurnMode: "auto",
			Padding:  300,
		},
		DefaultISO: models.ISOOptions{
			UDF: true,
			MD5: true,
		},
		BDXLSafeMode:       true,
		AutoEject:          true,
		DevicePollInterval: 5,
		Language:           "en",
		Theme:              "dark",
	}

	data, err := os.ReadFile(s.settingsPath())
	if err != nil {
		return settings, nil // return defaults
	}

	if err := json.Unmarshal(data, settings); err != nil {
		return settings, nil // return defaults on error
	}

	return settings, nil
}

// SaveSettings saves settings to disk
func (s *SettingsService) SaveSettings(settings *AppSettings) error {
	dir := filepath.Dir(s.settingsPath())
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.settingsPath(), data, 0644)
}

// GetXorrisoVersion returns xorriso version info
func (s *SettingsService) GetXorrisoVersion() (string, error) {
	return s.executor.Version(context.Background())
}
