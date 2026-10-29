package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolve(t *testing.T) {
	cases := []struct {
		name       string
		rawURL     string
		baseDir    string
		wantURL    string
		wantTarget string
	}{
		{"short form", "github.com/foo/bar", "/base", "https://github.com/foo/bar", "/base/bar"},
		{"https with .git", "https://github.com/foo/bar.git", "/base", "https://github.com/foo/bar.git", "/base/bar"},
		{"ssh", "git@github.com:foo/bar.git", "/base", "git@github.com:foo/bar.git", "/base/bar"},
		{"surrounding whitespace", "  github.com/foo/bar  ", "/base", "https://github.com/foo/bar", "/base/bar"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			url, target, err := resolve(tc.rawURL, tc.baseDir)
			if err != nil {
				t.Fatalf("resolve(%q) unexpected error: %v", tc.rawURL, err)
			}
			if url != tc.wantURL {
				t.Errorf("url = %q, want %q", url, tc.wantURL)
			}
			if target != tc.wantTarget {
				t.Errorf("target = %q, want %q", target, tc.wantTarget)
			}
		})
	}
}

func TestResolve_rejectsInvalid(t *testing.T) {
	for _, raw := range []string{"", "   ", "https://"} {
		if _, _, err := resolve(raw, "/base"); err == nil {
			t.Errorf("resolve(%q) = nil error, want error", raw)
		}
	}
}

func TestVersionFlag(t *testing.T) {
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"--version"})
	t.Cleanup(func() {
		rootCmd.SetArgs(nil)
		rootCmd.SetOut(nil)
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("--version returned error: %v", err)
	}
	if !strings.Contains(buf.String(), "version") {
		t.Errorf("--version output = %q, want it to contain %q", buf.String(), "version")
	}
}

func TestVersionString_fallback(t *testing.T) {
	if versionString() == "" {
		t.Error("versionString() must not be empty")
	}
}

func TestTargetOccupied(t *testing.T) {
	t.Run("nonexistent path is not occupied", func(t *testing.T) {
		target := filepath.Join(t.TempDir(), "does-not-exist")
		ok, err := targetOccupied(target)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ok {
			t.Error("nonexistent path should not be occupied")
		}
	})

	t.Run("empty dir is not occupied", func(t *testing.T) {
		target := filepath.Join(t.TempDir(), "empty")
		if err := os.Mkdir(target, 0755); err != nil {
			t.Fatal(err)
		}
		ok, err := targetOccupied(target)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ok {
			t.Error("empty dir should not be occupied")
		}
	})

	t.Run("non-empty dir is occupied", func(t *testing.T) {
		target := filepath.Join(t.TempDir(), "full")
		if err := os.Mkdir(target, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(target, "f"), []byte("x"), 0644); err != nil {
			t.Fatal(err)
		}
		ok, err := targetOccupied(target)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ok {
			t.Error("non-empty dir should be occupied")
		}
	})

	t.Run("existing file is occupied", func(t *testing.T) {
		target := filepath.Join(t.TempDir(), "file")
		if err := os.WriteFile(target, []byte("x"), 0644); err != nil {
			t.Fatal(err)
		}
		ok, err := targetOccupied(target)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ok {
			t.Error("existing file should be occupied")
		}
	})
}
