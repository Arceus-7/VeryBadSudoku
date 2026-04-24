package ui

import (
	"fmt"

	"github.com/arceus-7/chroma16"
)

// Theme holds all pre-built ANSI color escape sequences for the game.
// This is the ONLY file in the project that constructs ANSI color codes.
// All colors are derived from a chroma16 palette.
type Theme struct {
	palette chroma16.Palette

	// Foreground colors
	BoardBorder   string // grid lines and box-drawing characters
	CellFixed     string // given numbers (original puzzle clues)
	CellUser      string // user-placed numbers
	CellInvalid   string // conflicting numbers
	NoteColor     string // pencil marks (dim)
	StatusBar     string // status bar text
	TitleColor    string // titles and headers
	DimText       string // dimmed helper text
	AccentColor   string // accent for highlights

	// Background colors
	CellSelectedBg  string // cursor highlight background
	CellHighlightBg string // same-number highlight background
	CellBoxAltBg    string // alternating 3×3 box background (subtle)
	WinBg           string // win screen background

	// Foreground+Background combos
	CellSelectedFg string // text on selected background

	// Reset
	Reset string
}

// fgColor builds a 24-bit foreground ANSI escape from RGB values.
func fgColor(r, g, b uint8) string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
}

// bgColor builds a 24-bit background ANSI escape from RGB values.
func bgColor(r, g, b uint8) string {
	return fmt.Sprintf("\033[48;2;%d;%d;%dm", r, g, b)
}

// NewTheme creates a Theme by generating a chroma16 palette from the seed.
func NewTheme(seed string) (*Theme, error) {
	palette, err := chroma16.New().
		Seed(seed).
		Mood(chroma16.Dark).
		Contrast(chroma16.High).
		Build()
	if err != nil {
		return nil, fmt.Errorf("chroma16 palette: %w", err)
	}

	rgb := palette.RGB()

	// 16-color slot mapping:
	//  0=Black  1=Red     2=Green   3=Yellow  4=Blue    5=Magenta  6=Cyan    7=White
	//  8=BrBlk  9=BrRed  10=BrGrn  11=BrYel  12=BrBlu  13=BrMag  14=BrCyn  15=BrWht

	t := &Theme{
		palette: palette,

		// Foreground colors from palette slots
		BoardBorder: fgColor(rgb[14][0], rgb[14][1], rgb[14][2]), // Bright Cyan
		CellFixed:   fgColor(rgb[15][0], rgb[15][1], rgb[15][2]), // Bright White
		CellUser:    fgColor(rgb[10][0], rgb[10][1], rgb[10][2]), // Bright Green
		CellInvalid: fgColor(rgb[9][0], rgb[9][1], rgb[9][2]),    // Bright Red
		NoteColor:   fgColor(rgb[8][0], rgb[8][1], rgb[8][2]),    // Bright Black (gray)
		StatusBar:   fgColor(rgb[13][0], rgb[13][1], rgb[13][2]), // Bright Magenta
		TitleColor:  fgColor(rgb[14][0], rgb[14][1], rgb[14][2]), // Bright Cyan
		DimText:     fgColor(rgb[8][0], rgb[8][1], rgb[8][2]),    // Bright Black
		AccentColor: fgColor(rgb[11][0], rgb[11][1], rgb[11][2]), // Bright Yellow

		// Background colors
		CellSelectedBg:  bgColor(rgb[4][0], rgb[4][1], rgb[4][2]),  // Blue bg
		CellHighlightBg: bgColor(rgb[0][0]+20, rgb[0][1]+20, rgb[0][2]+30), // Subtle highlight
		CellBoxAltBg:    bgColor(rgb[0][0]+8, rgb[0][1]+8, rgb[0][2]+8),    // Very subtle alt bg
		WinBg:           bgColor(rgb[2][0], rgb[2][1], rgb[2][2]),           // Green bg

		// Selected cell foreground (bright white on blue bg)
		CellSelectedFg: fgColor(rgb[15][0], rgb[15][1], rgb[15][2]),

		Reset: "\033[0m",
	}

	return t, nil
}

// Colorize wraps text with the given foreground color and reset.
func (t *Theme) Colorize(text, color string) string {
	return color + text + t.Reset
}

// ColorBg wraps text with foreground + background colors and reset.
func (t *Theme) ColorBg(text, fg, bg string) string {
	return bg + fg + text + t.Reset
}

// Palette returns the underlying chroma16 palette.
func (t *Theme) Palette() chroma16.Palette {
	return t.palette
}
