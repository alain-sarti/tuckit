package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alain-sarti/tuckit/internal/config"
	"github.com/alain-sarti/tuckit/internal/git"
	"github.com/alain-sarti/tuckit/internal/tui"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "tuckit [url]",
	Short:        "Clone a git repository into your base directory",
	Args:         cobra.MaximumNArgs(1),
	RunE:         run,
	SilenceUsage: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
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
	url := git.NormalizeURL(rawURL)
	repoName := git.RepoName(url)
	targetPath := filepath.Join(cfg.BaseDir, repoName)

	// Conflict check
	if _, err := os.Stat(targetPath); err == nil {
		ok, err := tui.Confirm(fmt.Sprintf("%s already exists. Continue anyway?", targetPath))
		if err != nil {
			if errors.Is(err, huh.ErrUserAborted) {
				return nil
			}
			return err
		}
		if !ok {
			return nil
		}
	}

	// Clone
	if err := git.Clone(url, targetPath); err != nil {
		return err
	}

	fmt.Println(targetPath)

	// Optional zoxide registration
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

	return nil
}
