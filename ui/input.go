package ui

import (
	"bufio"
	"os"
)

// KeyEvent represents a parsed key input.
type KeyEvent int

const (
	KeyNone    KeyEvent = iota
	KeyUp               // arrow up
	KeyDown             // arrow down
	KeyLeft             // arrow left
	KeyRight            // arrow right
	KeyNum1             // digit 1
	KeyNum2             // digit 2
	KeyNum3             // digit 3
	KeyNum4             // digit 4
	KeyNum5             // digit 5
	KeyNum6             // digit 6
	KeyNum7             // digit 7
	KeyNum8             // digit 8
	KeyNum9             // digit 9
	KeyDelete           // backspace or delete — clear cell
	KeyNote             // 'n' or 'N' — toggle note mode
	KeyHint             // 'h' or 'H' — reveal a cell
	KeyNewGame          // 'r' or 'R' — new game
	KeyQuit             // 'q' or 'Q' or Ctrl-C
	KeyEnter            // enter / return
	KeyEscape           // escape (standalone, not part of a sequence)
)

var stdinReader = bufio.NewReaderSize(os.Stdin, 16)

// ReadKey reads and parses a single key event from stdin.
// Handles arrow keys (escape sequences), printable characters, and control keys.
func ReadKey() KeyEvent {
	b, err := stdinReader.ReadByte()
	if err != nil {
		return KeyNone
	}

	switch b {
	// Escape or escape sequence
	case 0x1b:
		return handleEscape()

	// Enter / Return
	case '\r', '\n':
		return KeyEnter

	// Backspace (0x7F on Unix, 0x08 on Windows)
	case 0x7f, 0x08:
		return KeyDelete

	// Ctrl-C
	case 0x03:
		return KeyQuit

	// Number keys
	case '1':
		return KeyNum1
	case '2':
		return KeyNum2
	case '3':
		return KeyNum3
	case '4':
		return KeyNum4
	case '5':
		return KeyNum5
	case '6':
		return KeyNum6
	case '7':
		return KeyNum7
	case '8':
		return KeyNum8
	case '9':
		return KeyNum9

	// Commands
	case 'n', 'N':
		return KeyNote
	case 'h', 'H':
		return KeyHint
	case 'r', 'R':
		return KeyNewGame
	case 'q', 'Q':
		return KeyQuit

	// WASD movement (alternative to arrows)
	case 'w', 'W':
		return KeyUp
	case 's', 'S':
		return KeyDown
	case 'a', 'A':
		return KeyLeft
	case 'd', 'D':
		return KeyRight

	// Delete key
	case 0x00:
		return KeyDelete

	default:
		return KeyNone
	}
}

// handleEscape processes escape sequences (arrow keys, etc.)
// or returns KeyEscape if it's a standalone escape press.
func handleEscape() KeyEvent {
	// Check if there are more bytes (part of an escape sequence)
	if stdinReader.Buffered() == 0 {
		// Standalone escape key
		return KeyEscape
	}

	b2, err := stdinReader.ReadByte()
	if err != nil {
		return KeyEscape
	}

	if b2 != '[' {
		return KeyEscape
	}

	// Read the final byte of the CSI sequence
	b3, err := stdinReader.ReadByte()
	if err != nil {
		return KeyEscape
	}

	switch b3 {
	case 'A':
		return KeyUp
	case 'B':
		return KeyDown
	case 'C':
		return KeyRight
	case 'D':
		return KeyLeft
	case '3':
		// Delete key sends \033[3~ — consume the tilde
		if stdinReader.Buffered() > 0 {
			_, _ = stdinReader.ReadByte()
		}
		return KeyDelete
	default:
		return KeyNone
	}
}
