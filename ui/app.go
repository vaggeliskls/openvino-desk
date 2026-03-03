package main

import (
	"context"
	"encoding/json"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/vaggeliskls/openvino-desk/ui/internal/setup"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed assets
var setupAssets embed.FS

const (
	defaultUvURL   = "https://github.com/astral-sh/uv/releases/download/0.7.3/uv-x86_64-pc-windows-msvc.zip"
	defaultOvmsURL = "https://github.com/openvinotoolkit/model_server/releases/download/v2026.0/ovms_windows_python_on.zip"
)

// Config holds user-configurable settings.
type Config struct {
	InstallDir string `json:"install_dir"`
	UvURL      string `json:"uv_url"`
	OvmsURL    string `json:"ovms_url"`
	StartupSet bool   `json:"startup_set"` // true once the startup preference has been written
}

// StatusResult reports whether each component is ready.
type StatusResult struct {
	UvReady   bool `json:"uv_ready"`
	DepsReady bool `json:"deps_ready"`
	OvmsReady bool `json:"ovms_ready"`
}

// App is the Wails application struct.
type App struct {
	ctx    context.Context
	config Config
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.loadConfig()
	// On first run, register the app to start with Windows by default.
	if !a.config.StartupSet {
		a.SetStartup(true) //nolint: errcheck — best-effort on first run
		a.config.StartupSet = true
		a.SaveConfig(a.config) //nolint: errcheck
	}
}

func configPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".openvino-desk", "config.json")
}

func defaultConfig() Config {
	home, _ := os.UserHomeDir()
	return Config{
		InstallDir: filepath.Join(home, "openvino-desk"),
		UvURL:      defaultUvURL,
		OvmsURL:    defaultOvmsURL,
	}
}

func (a *App) loadConfig() {
	data, err := os.ReadFile(configPath())
	if err != nil {
		a.config = defaultConfig()
		return
	}
	if err := json.Unmarshal(data, &a.config); err != nil {
		a.config = defaultConfig()
		return
	}
	// Fill in URL defaults for older configs that predate these fields.
	if a.config.UvURL == "" {
		a.config.UvURL = defaultUvURL
	}
	if a.config.OvmsURL == "" {
		a.config.OvmsURL = defaultOvmsURL
	}
}

// GetConfig returns the current configuration.
func (a *App) GetConfig() Config {
	return a.config
}

// SaveConfig persists the configuration to disk.
func (a *App) SaveConfig(config Config) error {
	a.config = config
	dir := filepath.Dir(configPath())
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath(), data, 0644)
}

// CheckStatus reports whether uv, the export venv, and OVMS are present.
func (a *App) CheckStatus() StatusResult {
	uvBin := filepath.Join(a.config.InstallDir, "uv.exe")
	venvPython := filepath.Join(a.config.InstallDir, "export", "Scripts", "python.exe")
	ovmsDir := filepath.Join(a.config.InstallDir, "ovms")

	_, uvErr := os.Stat(uvBin)
	_, depsErr := os.Stat(venvPython)
	_, ovmsErr := os.Stat(ovmsDir)

	return StatusResult{
		UvReady:   uvErr == nil,
		DepsReady: depsErr == nil,
		OvmsReady: ovmsErr == nil,
	}
}

// extractAssets writes embedded assets (requirements, scripts) to installDir,
// skipping uv.exe which is downloaded from the configured URL instead.
func (a *App) extractAssets() error {
	return fs.WalkDir(setupAssets, "assets", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Base(path) == "uv.exe" {
			return nil // uv is downloaded from URL, not embedded
		}
		rel, _ := filepath.Rel("assets", path)
		dest := filepath.Join(a.config.InstallDir, rel)
		if d.IsDir() {
			return os.MkdirAll(dest, 0755)
		}
		data, err := setupAssets.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(dest, data, 0755)
	})
}

func (a *App) emit(line string) {
	runtime.EventsEmit(a.ctx, "log", line)
}

// PrepareExport extracts bundled assets, downloads uv, then sets up the Python environment.
func (a *App) PrepareExport() error {
	if a.config.InstallDir == "" {
		return fmt.Errorf("install directory is not configured")
	}
	if a.config.UvURL == "" {
		return fmt.Errorf("uv download URL is not configured")
	}
	a.emit("Extracting bundled assets...")
	if err := a.extractAssets(); err != nil {
		return fmt.Errorf("extract assets: %w", err)
	}
	return setup.PrepareExport(a.config.InstallDir, a.config.UvURL, a.emit)
}

// PrepareOVMS downloads and extracts the OVMS server.
func (a *App) PrepareOVMS() error {
	if a.config.InstallDir == "" {
		return fmt.Errorf("install directory is not configured")
	}
	if a.config.OvmsURL == "" {
		return fmt.Errorf("OVMS download URL is not configured")
	}
	return setup.PrepareOVMS(a.config.InstallDir, a.config.OvmsURL, a.emit)
}
