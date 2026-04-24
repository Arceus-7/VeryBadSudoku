//go:build windows

package ui

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows"
)

const (
	enableVirtualTerminalProcessing = 0x0004
	enableVirtualTerminalInput      = 0x0200
	disableNewlineAutoReturn        = 0x0008
)

// RestoreFunc restores the terminal to its previous state.
type RestoreFunc func() error

// EnableRawMode puts the Windows console into raw input mode.
// Returns a function that restores the original mode.
func EnableRawMode() (RestoreFunc, error) {
	stdinHandle := windows.Handle(os.Stdin.Fd())
	stdoutHandle := windows.Handle(os.Stdout.Fd())

	var oldStdinMode, oldStdoutMode uint32

	if err := windows.GetConsoleMode(stdinHandle, &oldStdinMode); err != nil {
		return nil, fmt.Errorf("GetConsoleMode stdin: %w", err)
	}
	if err := windows.GetConsoleMode(stdoutHandle, &oldStdoutMode); err != nil {
		return nil, fmt.Errorf("GetConsoleMode stdout: %w", err)
	}

	// Raw input: disable line input, echo, processed input;
	// enable virtual terminal input for escape sequence support.
	rawStdinMode := (oldStdinMode &^ (windows.ENABLE_LINE_INPUT |
		windows.ENABLE_ECHO_INPUT |
		windows.ENABLE_PROCESSED_INPUT)) |
		enableVirtualTerminalInput

	if err := windows.SetConsoleMode(stdinHandle, rawStdinMode); err != nil {
		return nil, fmt.Errorf("SetConsoleMode stdin raw: %w", err)
	}

	restore := func() error {
		_ = windows.SetConsoleMode(stdinHandle, oldStdinMode)
		_ = windows.SetConsoleMode(stdoutHandle, oldStdoutMode)
		return nil
	}

	return restore, nil
}

// EnableANSI enables ANSI/VT escape sequence processing on Windows stdout.
func EnableANSI() error {
	stdoutHandle := windows.Handle(os.Stdout.Fd())

	var mode uint32
	if err := windows.GetConsoleMode(stdoutHandle, &mode); err != nil {
		return fmt.Errorf("GetConsoleMode stdout: %w", err)
	}

	mode |= enableVirtualTerminalProcessing | disableNewlineAutoReturn

	if err := windows.SetConsoleMode(stdoutHandle, mode); err != nil {
		return fmt.Errorf("SetConsoleMode stdout ANSI: %w", err)
	}
	return nil
}

// GetTerminalSize returns the current terminal width and height.
func GetTerminalSize() (int, int, error) {
	stdoutHandle := windows.Handle(os.Stdout.Fd())
	var info windows.ConsoleScreenBufferInfo
	if err := windows.GetConsoleScreenBufferInfo(stdoutHandle, &info); err != nil {
		return 80, 24, nil // safe fallback
	}
	width := int(info.Window.Right - info.Window.Left + 1)
	height := int(info.Window.Bottom - info.Window.Top + 1)
	return width, height, nil
}
