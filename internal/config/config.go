package config

import (
	"errors"
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
// TUCKIT_BASE_DIR takes precedence over the config file value.
// Returns an error if neither source provides a value.
func Load() (*Config, error) {
	baseDir := ""

	// Try config file first. A missing file is fine (fall back to the env
	// var); a malformed file is a real error worth reporting.
	cfgPath, err := configFilePath()
	if err == nil {
		var cf configFile
		if _, derr := toml.DecodeFile(cfgPath, &cf); derr != nil {
			if !errors.Is(derr, os.ErrNotExist) {
				return nil, fmt.Errorf("reading %s: %w", cfgPath, derr)
			}
		} else if cf.BaseDir != "" {
			baseDir = cf.BaseDir
		}
	}

	// Env var takes precedence
	if env := os.Getenv("TUCKIT_BASE_DIR"); env != "" {
		baseDir = env
	}

	if baseDir == "" {
		return nil, fmt.Errorf("no base directory configured: set TUCKIT_BASE_DIR or add base_dir to ~/.config/tuckit/config.toml")
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
	return filepath.Join(home, ".config", "tuckit", "config.toml"), nil
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
