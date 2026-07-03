package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoad_envVarOnly(t *testing.T) {
	t.Setenv("HOME", t.TempDir()) // isolate from real config file
	t.Setenv("TUCKIT_BASE_DIR", "/tmp/test-dev")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.BaseDir != "/tmp/test-dev" {
		t.Errorf("got %q, want %q", cfg.BaseDir, "/tmp/test-dev")
	}
}

func TestLoad_noConfig_returnsError(t *testing.T) {
	t.Setenv("TUCKIT_BASE_DIR", "")
	t.Setenv("HOME", t.TempDir()) // no config file in temp home
	_, err := Load()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "no base directory configured") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestLoad_tildeExpansion(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	t.Setenv("TUCKIT_BASE_DIR", "~/Dev")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := filepath.Join(tmpHome, "Dev")
	if cfg.BaseDir != expected {
		t.Errorf("got %q, want %q", cfg.BaseDir, expected)
	}
}

func TestLoad_envVarOverridesConfigFile(t *testing.T) {
	// Write a config file in a temp HOME
	tmpHome := t.TempDir()
	cfgDir := filepath.Join(tmpHome, ".config", "tuckit")
	os.MkdirAll(cfgDir, 0755)
	os.WriteFile(filepath.Join(cfgDir, "config.toml"), []byte(`base_dir = "/from-file"`), 0644)

	t.Setenv("HOME", tmpHome)
	t.Setenv("TUCKIT_BASE_DIR", "/from-env")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.BaseDir != "/from-env" {
		t.Errorf("got %q, want env value %q", cfg.BaseDir, "/from-env")
	}
}

func TestLoad_configFileOnly(t *testing.T) {
	tmpHome := t.TempDir()
	cfgDir := filepath.Join(tmpHome, ".config", "tuckit")
	os.MkdirAll(cfgDir, 0755)
	os.WriteFile(filepath.Join(cfgDir, "config.toml"), []byte(`base_dir = "/from-file"`), 0644)

	t.Setenv("HOME", tmpHome)
	t.Setenv("TUCKIT_BASE_DIR", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.BaseDir != "/from-file" {
		t.Errorf("got %q, want %q", cfg.BaseDir, "/from-file")
	}
}
