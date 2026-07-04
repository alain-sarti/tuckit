package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/alain-sarti/tuckit/internal/config"
	"github.com/alain-sarti/tuckit/internal/git"
	"github.com/alain-sarti/tuckit/internal/tui"
	"github.com/charmbracelet/huh"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "tuckit [url]",
	Short:        "Clone a git repository into your base directory",
	Version:      versionString(),
	Args:         cobra.MaximumNArgs(1),
	RunE:         run,
	SilenceUsage: true,
}

// versionString reports the build version, derived from the module version
// embedded by `go install`. Falls back to "dev" for local/unversioned builds.
func versionString() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		if v := info.Main.Version; v != "" && v != "(devel)" {
			return v
		}
	}
	return "dev"
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// resolve validates rawURL and computes the clone URL and the target path
// under baseDir. It returns an error for input that has no usable repository
// name rather than silently targeting baseDir itself.
func resolve(rawURL, baseDir string) (url, target string, err error) {
	if err := git.ValidateURL(rawURL); err != nil {
		return "", "", err
	}
	url = git.NormalizeURL(strings.TrimSpace(rawURL))
	target = filepath.Join(baseDir, git.RepoName(url))
	return url, target, nil
}

// targetAvailable reports whether a clone can be placed at target. git clone
// refuses any existing file or non-empty directory, so those are unavailable;
// a missing path or an empty directory is fine.
func targetAvailable(target string) (bool, error) {
	info, err := os.Stat(target)
	if errors.Is(err, os.ErrNotExist) {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	if !info.IsDir() {
		return false, nil
	}
	entries, err := os.ReadDir(target)
	if err != nil {
		return false, err
	}
	return len(entries) == 0, nil
}

func run(cmd *cobra.Command, args []string) error {
	// Get URL: from args or interactive prompt
	var rawURL string
	if len(args) == 1 {
		rawURL = args[0]
	} else {
		var err error
		rawURL, err = tui.AskURL()
		if err != nil {
			if errors.Is(err, huh.ErrUserAborted) {
				return nil
			}
			return err
		}
	}

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// Resolve target path
	url, targetPath, err := resolve(rawURL, cfg.BaseDir)
	if err != nil {
		return err
	}

	// Conflict check: git clone won't write into an occupied target.
	avail, err := targetAvailable(targetPath)
	if err != nil {
		return err
	}
	if !avail {
		return fmt.Errorf("target already exists: %s", targetPath)
	}

	// Clone
	if err := git.Clone(url, targetPath); err != nil {
		return err
	}

	fmt.Println(targetPath)

	// Optional zoxide registration — only when we can actually prompt.
	// Without a terminal (piped/CI), skip it so a successful clone still exits 0.
	if isatty.IsTerminal(os.Stdin.Fd()) {
		ok, err := tui.Confirm("Add to zoxide?")
		if err != nil {
			if errors.Is(err, huh.ErrUserAborted) {
				return nil
			}
			return err
		}
		if ok {
			if err := git.ZoxideAdd(targetPath); err != nil {
				fmt.Fprintf(os.Stderr, "warning: could not add to zoxide: %v\n", err)
			}
		}
	}

	return nil
}
