package tui

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/huh"
)

// Confirm shows a yes/no prompt. Returns true if the user confirms.
// Returns an error if no TTY is available or the prompt otherwise fails.
// If the user aborts (Ctrl+C), returns huh.ErrUserAborted — callers should treat this as exit 0.
func Confirm(question string) (bool, error) {
	var result bool
	err := huh.NewConfirm().
		Title(question).
		Value(&result).
		Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return false, huh.ErrUserAborted
		}
		return false, fmt.Errorf("interactive prompts require a terminal")
	}
	return result, nil
}

// AskURL shows a text input prompt asking for a git repository URL.
// Returns an error if no TTY is available.
// If the user aborts (Ctrl+C), returns huh.ErrUserAborted — callers should treat this as exit 0.
func AskURL() (string, error) {
	var repoURL string
	err := huh.NewInput().
		Title("Git repository URL").
		Placeholder("github.com/user/repo").
		Value(&repoURL).
		Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return "", huh.ErrUserAborted
		}
		return "", fmt.Errorf("interactive prompts require a terminal")
	}
	return repoURL, nil
}
