package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"github.com/foo/bar", "https://github.com/foo/bar"},
		{"https://github.com/foo/bar", "https://github.com/foo/bar"},
		{"http://github.com/foo/bar", "http://github.com/foo/bar"},
		{"git@github.com:foo/bar.git", "git@github.com:foo/bar.git"},
	}
	for _, tc := range cases {
		got := NormalizeURL(tc.input)
		if got != tc.want {
			t.Errorf("NormalizeURL(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestRepoName(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"https://github.com/foo/bar", "bar"},
		{"https://github.com/foo/bar.git", "bar"},
		{"git@github.com:foo/bar.git", "bar"},
		{"git@github.com:bar.git", "bar"}, // single-level SSH path
		{"github.com/foo/bar", "bar"},
	}
	for _, tc := range cases {
		got := RepoName(tc.input)
		if got != tc.want {
			t.Errorf("RepoName(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

// setupRepo creates a temporary git repo with one commit for testing.
func setupRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	runCmd(t, dir, "git", "init")
	runCmd(t, dir, "git", "config", "user.email", "test@example.com")
	runCmd(t, dir, "git", "config", "user.name", "Test")
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}
	runCmd(t, dir, "git", "add", ".")
	runCmd(t, dir, "git", "commit", "-m", "init")
	return dir
}

func runCmd(t *testing.T, dir string, name string, args ...string) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%s %v: %v\n%s", name, args, err, out)
	}
}

func TestClone_success(t *testing.T) {
	src := setupRepo(t)
	dst := filepath.Join(t.TempDir(), "cloned")

	if err := Clone(src, dst); err != nil {
		t.Fatalf("Clone failed: %v", err)
	}
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		t.Fatal("expected destination directory to exist after clone")
	}
}

func TestClone_gitNotFound(t *testing.T) {
	// Temporarily clear PATH so git cannot be found
	t.Setenv("PATH", "")
	err := Clone("https://example.com/foo/bar", "/tmp/irrelevant")
	if err == nil {
		t.Fatal("expected error when git not on PATH")
	}
	if err.Error() != "git not found on PATH" {
		t.Errorf("unexpected error message: %v", err)
	}
}
