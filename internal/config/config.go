package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// Config holds resolved configuration values.
type Config struct {
	BaseDir string
}

type configFile struct {
	BaseDir string `toml:"base_dir"`
}

// Load resolves the base directory from config file and environment.
// GCO_BASE_DIR takes precedence over the config file value.
// Returns an error if neither source provides a value.
func Load() (*Config, error) {
	baseDir := ""

	// Try config file first
	cfgPath, err := configFilePath()
	if err == nil {
		var cf configFile
		if _, err := toml.DecodeFile(cfgPath, &cf); err == nil && cf.BaseDir != "" {
			baseDir = cf.BaseDir
		}
	}

	// Env var takes precedence
	if env := os.Getenv("GCO_BASE_DIR"); env != "" {
		baseDir = env
	}

	if baseDir == "" {
		return nil, fmt.Errorf("No base directory configured. Set GCO_BASE_DIR or add base_dir to ~/.config/gco/config.toml")
	}

	expanded, err := expandTilde(baseDir)
	if err != nil {
		return nil, fmt.Errorf("expanding base_dir: %w", err)
	}

	return &Config{BaseDir: expanded}, nil
}

func configFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "gco", "config.toml"), nil
}

func expandTilde(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, path[1:]), nil
}
