package services

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetSettings_Defaults(t *testing.T) {
	// Используем XDG_CONFIG_HOME для изоляции от реальных настроек
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	svc := NewSettingsService(&mockRunner{})

	settings, err := svc.GetSettings()
	if err != nil {
		t.Fatalf("GetSettings: %v", err)
	}

	if settings.Language != "en" {
		t.Errorf("Language = %q, want en", settings.Language)
	}
	if settings.Theme != "dark" {
		t.Errorf("Theme = %q, want dark", settings.Theme)
	}
	// XorrisoPath должен быть либо "xorriso", либо путь к бинарнику
	if settings.XorrisoPath == "" {
		t.Error("XorrisoPath should not be empty")
	}
}

func TestSaveAndGetSettings(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	svc := NewSettingsService(&mockRunner{})

	// Сохраняем настройки
	toSave := &AppSettings{
		XorrisoPath: "/usr/bin/xorriso",
		Language:    "ru",
		Theme:       "light",
	}
	if err := svc.SaveSettings(toSave); err != nil {
		t.Fatalf("SaveSettings: %v", err)
	}

	// Загружаем обратно
	loaded, err := svc.GetSettings()
	if err != nil {
		t.Fatalf("GetSettings: %v", err)
	}

	if loaded.XorrisoPath != "/usr/bin/xorriso" {
		t.Errorf("XorrisoPath = %q, want /usr/bin/xorriso", loaded.XorrisoPath)
	}
	if loaded.Language != "ru" {
		t.Errorf("Language = %q, want ru", loaded.Language)
	}
	if loaded.Theme != "light" {
		t.Errorf("Theme = %q, want light", loaded.Theme)
	}
}

func TestGetSettings_CorruptFile(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	// Создаём невалидный JSON по пути настроек
	settingsDir := filepath.Join(tmpDir, "xorriso-ui")
	os.MkdirAll(settingsDir, 0755)
	settingsFile := filepath.Join(settingsDir, "settings.json")
	os.WriteFile(settingsFile, []byte("{invalid json!!!"), 0644)

	svc := NewSettingsService(&mockRunner{})
	settings, err := svc.GetSettings()
	if err != nil {
		t.Fatalf("GetSettings should not return error for corrupt file: %v", err)
	}

	// Должен вернуть defaults
	if settings.Language != "en" {
		t.Errorf("Language = %q, want en (default)", settings.Language)
	}
	if settings.Theme != "dark" {
		t.Errorf("Theme = %q, want dark (default)", settings.Theme)
	}
}

func TestGetXorrisoVersion(t *testing.T) {
	expectedVersion := "xorriso 1.5.6 : RockRidge filesystem manipulator"

	runner := &mockRunner{
		VersionFn: func(ctx context.Context) (string, error) {
			return expectedVersion, nil
		},
	}

	svc := NewSettingsService(runner)

	version, err := svc.GetXorrisoVersion()
	if err != nil {
		t.Fatalf("GetXorrisoVersion: %v", err)
	}
	if version != expectedVersion {
		t.Errorf("Version = %q, want %q", version, expectedVersion)
	}
}

func TestSettingsPath(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	svc := NewSettingsService(&mockRunner{})
	path := svc.settingsPath()

	if !strings.Contains(path, "xorriso-ui") {
		t.Errorf("settings path should contain 'xorriso-ui', got %q", path)
	}
	if !strings.HasSuffix(path, "settings.json") {
		t.Errorf("settings path should end with 'settings.json', got %q", path)
	}
}

func TestSaveSettings_CreatesDir(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	svc := NewSettingsService(&mockRunner{})

	settings := &AppSettings{
		XorrisoPath: "/usr/bin/xorriso",
		Language:    "en",
		Theme:       "dark",
	}
	if err := svc.SaveSettings(settings); err != nil {
		t.Fatalf("SaveSettings: %v", err)
	}

	// Проверяем, что файл был создан
	data, err := os.ReadFile(svc.settingsPath())
	if err != nil {
		t.Fatalf("settings file not created: %v", err)
	}

	var loaded AppSettings
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("invalid JSON in settings file: %v", err)
	}
	if loaded.Language != "en" {
		t.Errorf("Language in file = %q, want en", loaded.Language)
	}
}
