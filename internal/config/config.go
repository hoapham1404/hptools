package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Config holds the application configuration
type Config struct {
	App    AppConfig    `json:"app"`
	Window WindowConfig `json:"window"`
	Log    LogConfig    `json:"log"`
}

// AppConfig holds general application settings
type AppConfig struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// WindowConfig holds window-specific settings
type WindowConfig struct {
	Title            string           `json:"title"`
	Width            int              `json:"width"`
	Height           int              `json:"height"`
	MinWidth         int              `json:"minWidth"`
	MinHeight        int              `json:"minHeight"`
	MaxWidth         int              `json:"maxWidth"`
	MaxHeight        int              `json:"maxHeight"`
	BackgroundColour application.RGBA `json:"backgroundColour"`
	Mac              MacWindowConfig  `json:"mac"`
}

// MacWindowConfig holds macOS-specific window settings
type MacWindowConfig struct {
	InvisibleTitleBarHeight int                     `json:"invisibleTitleBarHeight"`
	Backdrop                application.MacBackdrop `json:"backdrop"`
	TitleBar                application.MacTitleBar `json:"titleBar"`
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"` // "text" or "json"
}

// Default returns a configuration with sensible defaults
func Default() *Config {
	return &Config{
		App: AppConfig{
			Name:        "hptools",
			Description: "HP Tools - Window Management Application",
		},
		Window: WindowConfig{
			Title:            "HP Tools",
			Width:            1200,
			Height:           1000,
			MinWidth:         800,
			MinHeight:        600,
			MaxWidth:         1920,
			MaxHeight:        1080,
			BackgroundColour: application.NewRGB(27, 38, 54),
			Mac: MacWindowConfig{
				InvisibleTitleBarHeight: 50,
				Backdrop:                application.MacBackdropTranslucent,
				TitleBar:                application.MacTitleBarHiddenInset,
			},
		},
		Log: LogConfig{
			Level:  "info",
			Format: "text",
		},
	}
}

// Load loads configuration from file, falling back to defaults
func Load(configPath string) (*Config, error) {
	cfg := Default()

	if configPath == "" {
		return cfg, nil
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Config file doesn't exist, create it with defaults
		if err := Save(cfg, configPath); err != nil {
			return cfg, fmt.Errorf("creating default config file: %w", err)
		}
		return cfg, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return cfg, fmt.Errorf("reading config file: %w", err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return cfg, fmt.Errorf("parsing config file: %w", err)
	}

	return cfg, nil
}

// Save saves configuration to file
func Save(cfg *Config, configPath string) error {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	return nil
}

// GetConfigPath returns the default configuration file path
func GetConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "./config.json"
	}
	return filepath.Join(home, ".config", "hptools", "config.json")
}
