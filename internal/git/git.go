package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

// ValidateURL reports whether input is a usable git URL — non-empty and
// yielding a real repository name. It rejects blank input and URLs from which
// no meaningful name can be derived (which would otherwise resolve to the base
// directory itself).
func ValidateURL(input string) error {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return errors.New("repository URL is required")
	}
	name := RepoName(NormalizeURL(trimmed))
	switch name {
	case "", ".", "..", "/":
		return fmt.Errorf("could not derive a repository name from %q", input)
	}
	return nil
}

// NormalizeURL ensures the URL has a protocol.
// Short form (e.g. github.com/foo/bar) gets https:// prepended.
// SCP-like SSH (git@…) and URLs that already carry a scheme
// (https://, http://, ssh://, git://, …) are returned unchanged.
func NormalizeURL(input string) string {
	if strings.HasPrefix(input, "git@") || strings.Contains(input, "://") {
		return input
	}
	return "https://" + input
}

// Clone runs git clone, placing the repo at targetPath.
// git's stderr is passed through directly so the user sees progress.
// Returns a clear error if git is not on PATH.
func Clone(url, targetPath string) error {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		return fmt.Errorf("git not found on PATH")
	}
	// "--" stops git from interpreting a url or path as an option.
	cmd := exec.Command(gitPath, "clone", "--", url, targetPath)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}
	return nil
}

// ZoxideAdd registers path with zoxide.
// Returns an error if zoxide is not on PATH or exits non-zero,
// but callers should treat this as a warning, not a fatal error.
func ZoxideAdd(path string) error {
	zoxidePath, err := exec.LookPath("zoxide")
	if err != nil {
		return fmt.Errorf("zoxide not found on PATH")
	}
	cmd := exec.Command(zoxidePath, "add", path)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("zoxide add: %s", out)
	}
	return nil
}

// RepoName derives the repository name from a URL.
// Returns the last path segment with any .git suffix stripped.
func RepoName(url string) string {
	// SSH URLs: git@github.com:foo/bar.git — extract path after ":"
	if strings.HasPrefix(url, "git@") {
		parts := strings.SplitN(url, ":", 2)
		if len(parts) == 2 {
			url = parts[1]
		}
	}
	// Strip protocol prefix
	if idx := strings.Index(url, "://"); idx != -1 {
		url = url[idx+3:]
	}
	base := path.Base(url)
	return strings.TrimSuffix(base, ".git")
}
