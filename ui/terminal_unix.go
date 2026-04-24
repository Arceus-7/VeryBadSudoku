//go:build !windows

package ui

import (
	"os"

	"golang.org/x/term"
)

// RestoreFunc restores the terminal to its previous state.
type RestoreFunc func() error

// EnableRawMode puts the Unix terminal into raw mode.
// Returns a function that restores the original mode.
func EnableRawMode() (RestoreFunc, error) {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return nil, err
	}
	return func() error {
		return term.Restore(fd, oldState)
	}, nil
}

// EnableANSI is a no-op on Unix — ANSI is always available.
func EnableANSI() error {
	return nil
}

// GetTerminalSize returns the current terminal width and height.
func GetTerminalSize() (int, int, error) {
	return term.GetSize(int(os.Stdout.Fd()))
}
