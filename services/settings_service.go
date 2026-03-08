package services

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"xorriso-ui/pkg/xorriso"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type AppSettings struct {
	XorrisoPath string `json:"xorrisoPath"`
	Language    string `json:"language"`
	Theme       string `json:"theme"`
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
		Language:    "en",
		Theme:       "dark",
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

// ToolInfo contains name, path and version of an external tool
type ToolInfo struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Version string `json:"version"`
	Found   bool   `json:"found"`
}

// GetXorrisoVersion returns xorriso version info
func (s *SettingsService) GetXorrisoVersion() (string, error) {
	return s.executor.Version(context.Background())
}

// GetToolsInfo returns information about required external tools (xorriso, mkisofs)
func (s *SettingsService) GetToolsInfo() ([]ToolInfo, error) {
	tools := []ToolInfo{
		{Name: "xorriso"},
		{Name: "mkisofs"},
	}

	for i := range tools {
		p, err := exec.LookPath(tools[i].Name)
		if err != nil {
			continue
		}
		tools[i].Path = p
		tools[i].Found = true

		cmd := exec.CommandContext(context.Background(), p, "--version")
		output, err := cmd.CombinedOutput()
		if err == nil {
			// Take only the first line of version output
			lines := strings.SplitN(strings.TrimSpace(string(output)), "\n", 2)
			if len(lines) > 0 {
				tools[i].Version = lines[0]
			}
		}
	}

	return tools, nil
}
