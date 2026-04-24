package ui

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/arceus-7/sudoku/game"
	"github.com/arceus-7/sudoku/state"
)

// displayWidth returns the visual display width of a string.
// Uses rune count which is correct for ASCII, box-drawing, and symbols used here.
func displayWidth(s string) int {
	return utf8.RuneCountInString(s)
}

// Box-drawing characters for the Sudoku grid.
const (
	// Heavy borders (3×3 block boundaries)
	topLeft     = "╔"
	topRight    = "╗"
	bottomLeft  = "╚"
	bottomRight = "╝"
	heavyH      = "═"
	heavyV      = "║"
	heavyTDown  = "╦"
	heavyTUp    = "╩"
	heavyTRight = "╠"
	heavyTLeft  = "╣"
	heavyCross  = "╬"

	// Light borders (within 3×3 blocks)
	lightH     = "─"
	lightV     = "│"
	lightCross = "┼"

	// Mixed intersections
	heavyDownLightH  = "╤"
	heavyUpLightH    = "╧"
	heavyRightLightV = "╟"
	heavyLeftLightV  = "╢"
	lightCrossHeavyH = "╪"
	lightCrossHeavyV = "╫"
)

// Renderer handles all screen drawing.
type Renderer struct {
	theme *Theme
	buf   strings.Builder
}

// NewRenderer creates a renderer with the given theme.
func NewRenderer(theme *Theme) *Renderer {
	return &Renderer{
		theme: theme,
	}
}

// EnterAlternateScreen switches to the alternate screen buffer.
func EnterAlternateScreen() {
	fmt.Print("\033[?1049h\033[2J\033[H")
}

// ExitAlternateScreen returns to the main screen buffer.
func ExitAlternateScreen() {
	fmt.Print("\033[?1049l")
}

// HideCursor hides the terminal cursor.
func HideCursor() {
	fmt.Print("\033[?25l")
}

// ShowCursor shows the terminal cursor.
func ShowCursor() {
	fmt.Print("\033[?25h")
}

// flush writes the buffer to stdout in a single call to minimize flicker.
func (r *Renderer) flush() {
	os.Stdout.WriteString(r.buf.String())
	r.buf.Reset()
}

// moveTo appends cursor-movement escape to position (row, col) — 1-indexed.
func (r *Renderer) moveTo(row, col int) {
	fmt.Fprintf(&r.buf, "\033[%d;%dH", row, col)
}

// clearScreen appends the clear-screen escape sequence.
func (r *Renderer) clearScreen() {
	r.buf.WriteString("\033[2J\033[H")
}

// ─────────────────── Game Screen ───────────────────

// RenderGame draws the complete game screen: title, info bar, board, controls.
func (r *Renderer) RenderGame(gs *state.GameState) {
	r.buf.Reset()
	r.clearScreen()

	t := r.theme
	width, height, _ := GetTerminalSize()

	if width < 40 || height < 25 {
		r.buf.WriteString("Terminal too small!\r\nMinimum: 40×25\r\n")
		r.flush()
		return
	}

	// Board dimensions (Unicode box-drawing): 37 chars wide × 19 rows tall
	boardWidth := 37
	boardHeight := 19

	// Center board
	startCol := (width - boardWidth) / 2
	if startCol < 1 {
		startCol = 1
	}
	startRow := (height-boardHeight-8)/2 + 1
	if startRow < 6 {
		startRow = 6
	}

	// ── Title ──
	titleLines := []struct {
		text  string
		color string
	}{
		{"╔═══════════════════════╗", t.TitleColor},
		{"║  V E R Y  B A D      ║", t.AccentColor},
		{"║    S U D O K U       ║", t.AccentColor},
		{"╚═══════════════════════╝", t.TitleColor},
	}
	titleW := 25
	titleCol := (width - titleW) / 2
	if titleCol < 1 {
		titleCol = 1
	}
	for i, tl := range titleLines {
		r.moveTo(1+i, titleCol)
		r.buf.WriteString(t.Colorize(tl.text, tl.color))
	}

	// ── Note-mode indicator ──
	if gs.NoteMode {
		label := " PENCIL MODE "
		r.moveTo(startRow-2, (width-displayWidth(label))/2)
		r.buf.WriteString(t.ColorBg(label, t.CellSelectedFg, t.CellSelectedBg))
	}

	// ── Info bar ──
	elapsed := gs.ElapsedTime()
	min := int(elapsed.Minutes())
	sec := int(elapsed.Seconds()) % 60
	info := fmt.Sprintf(" %s  |  %02d:%02d  |  Errors: %d  |  Hints: %d ",
		gs.Difficulty.String(), min, sec, gs.Errors, gs.HintsUsed)
	r.moveTo(startRow-1, (width-len(info))/2)
	r.buf.WriteString(t.Colorize(info, t.StatusBar))

	// ── Board ──
	curRow := startRow

	// Top border
	r.moveTo(curRow, startCol)
	r.buf.WriteString(t.Colorize(buildTopBorder(), t.BoardBorder))
	curRow++

	for gridRow := 0; gridRow < 9; gridRow++ {
		r.moveTo(curRow, startCol)
		r.writeDataRow(gs, gridRow)
		curRow++

		if gridRow < 8 {
			r.moveTo(curRow, startCol)
			if (gridRow+1)%3 == 0 {
				r.buf.WriteString(t.Colorize(buildHeavySeparator(), t.BoardBorder))
			} else {
				r.buf.WriteString(t.Colorize(buildLightSeparator(), t.BoardBorder))
			}
			curRow++
		}
	}

	// Bottom border
	r.moveTo(curRow, startCol)
	r.buf.WriteString(t.Colorize(buildBottomBorder(), t.BoardBorder))
	curRow++

	// ── Controls legend ──
	r.drawControls(curRow+1, width)

	r.flush()
}

// writeDataRow appends one row of cell values to the buffer.
func (r *Renderer) writeDataRow(gs *state.GameState, gridRow int) {
	t := r.theme
	board := gs.Board
	cursorVal := board.Get(gs.CursorRow, gs.CursorCol).Value

	for gridCol := 0; gridCol < 9; gridCol++ {
		// Vertical separator
		if gridCol%3 == 0 {
			r.buf.WriteString(t.Colorize(heavyV, t.BoardBorder))
		} else {
			r.buf.WriteString(t.Colorize(lightV, t.BoardBorder))
		}

		cell := board.Get(gridRow, gridCol)
		isSelected := gridRow == gs.CursorRow && gridCol == gs.CursorCol
		isSameNum := cursorVal != 0 && cell.Value == cursorVal && !isSelected

		r.writeCell(cell, isSelected, isSameNum)
	}

	// Closing vertical
	r.buf.WriteString(t.Colorize(heavyV, t.BoardBorder))
}

// writeCell appends the 3-character cell content (" X ") to the buffer.
func (r *Renderer) writeCell(cell *game.Cell, isSelected, isSameNumber bool) {
	t := r.theme

	var content string
	var fg string

	if cell.Value != 0 {
		content = fmt.Sprintf(" %d ", cell.Value)
		switch {
		case cell.Invalid:
			fg = t.CellInvalid
		case cell.Fixed:
			fg = t.CellFixed
		default:
			fg = t.CellUser
		}
	} else if cell.HasNotes() {
		// Compact note display: show up to 3 notes in the cell
		var notes []byte
		for i := 1; i <= 9; i++ {
			if cell.Notes[i] {
				notes = append(notes, byte('0'+i))
			}
		}
		switch {
		case len(notes) >= 3:
			content = string(notes[:3])
		case len(notes) == 2:
			content = " " + string(notes[:2])
		case len(notes) == 1:
			content = " " + string(notes[:1]) + " "
		default:
			content = " · "
		}
		fg = t.NoteColor
	} else {
		content = " · "
		fg = t.DimText
	}

	switch {
	case isSelected:
		r.buf.WriteString(t.ColorBg(content, t.CellSelectedFg, t.CellSelectedBg))
	case isSameNumber:
		r.buf.WriteString(t.ColorBg(content, fg, t.CellHighlightBg))
	default:
		r.buf.WriteString(t.Colorize(content, fg))
	}
}

// drawControls renders the key legend.
func (r *Renderer) drawControls(row, width int) {
	t := r.theme
	line := "Arrows Move  |  1-9 Place  |  Del Clear  |  N Notes  |  H Hint  |  R New  |  Q Quit"
	col := (width - displayWidth(line)) / 2
	if col < 1 {
		col = 1
	}
	r.moveTo(row, col)
	r.buf.WriteString(t.Colorize(line, t.DimText))
}

// ─────────────────── Grid Line Builders ───────────────────

func buildTopBorder() string {
	var sb strings.Builder
	sb.WriteString(topLeft)
	for block := 0; block < 3; block++ {
		for cell := 0; cell < 3; cell++ {
			sb.WriteString(heavyH + heavyH + heavyH)
			if cell < 2 {
				sb.WriteString(heavyDownLightH)
			}
		}
		if block < 2 {
			sb.WriteString(heavyTDown)
		}
	}
	sb.WriteString(topRight)
	return sb.String()
}

func buildBottomBorder() string {
	var sb strings.Builder
	sb.WriteString(bottomLeft)
	for block := 0; block < 3; block++ {
		for cell := 0; cell < 3; cell++ {
			sb.WriteString(heavyH + heavyH + heavyH)
			if cell < 2 {
				sb.WriteString(heavyUpLightH)
			}
		}
		if block < 2 {
			sb.WriteString(heavyTUp)
		}
	}
	sb.WriteString(bottomRight)
	return sb.String()
}

func buildHeavySeparator() string {
	var sb strings.Builder
	sb.WriteString(heavyTRight)
	for block := 0; block < 3; block++ {
		for cell := 0; cell < 3; cell++ {
			sb.WriteString(heavyH + heavyH + heavyH)
			if cell < 2 {
				sb.WriteString(lightCrossHeavyH)
			}
		}
		if block < 2 {
			sb.WriteString(heavyCross)
		}
	}
	sb.WriteString(heavyTLeft)
	return sb.String()
}

func buildLightSeparator() string {
	var sb strings.Builder
	sb.WriteString(heavyRightLightV)
	for block := 0; block < 3; block++ {
		for cell := 0; cell < 3; cell++ {
			sb.WriteString(lightH + lightH + lightH)
			if cell < 2 {
				sb.WriteString(lightCross)
			}
		}
		if block < 2 {
			sb.WriteString(lightCrossHeavyV)
		}
	}
	sb.WriteString(heavyLeftLightV)
	return sb.String()
}

// ─────────────────── Difficulty Select Screen ───────────────────

// RenderDifficultySelect draws the difficulty selection menu.
func (r *Renderer) RenderDifficultySelect(selected int) {
	r.buf.Reset()
	r.clearScreen()
	t := r.theme

	width, height, _ := GetTerminalSize()

	titleLines := []struct {
		text  string
		color string
	}{
		{"╔═════════════════════════════╗", t.TitleColor},
		{"║                             ║", t.TitleColor},
		{"║    V E R Y   B A D         ║", t.AccentColor},
		{"║      S U D O K U           ║", t.AccentColor},
		{"║                             ║", t.TitleColor},
		{"╚═════════════════════════════╝", t.TitleColor},
	}

	topRow := height/2 - 8
	if topRow < 2 {
		topRow = 2
	}

	for i, tl := range titleLines {
		col := (width - 31) / 2
		if col < 1 {
			col = 1
		}
		r.moveTo(topRow+i, col)
		r.buf.WriteString(t.Colorize(tl.text, tl.color))
	}

	// Subtitle
	sub := "Select Difficulty"
	r.moveTo(topRow+7, (width-displayWidth(sub))/2)
	r.buf.WriteString(t.Colorize(sub, t.StatusBar))

	// Options
	labels := []string{"  Easy  ", " Medium ", "  Hard  "}
	descs := []string{
		"Perfect for beginners - 35 cells removed",
		"A balanced challenge - 45 cells removed",
		"For Sudoku masters - 52 cells removed",
	}

	for i, label := range labels {
		optRow := topRow + 9 + (i * 3)
		optCol := (width - displayWidth(label) - 4) / 2

		if i == selected {
			r.moveTo(optRow, optCol)
			r.buf.WriteString(t.Colorize(" > ", t.AccentColor))
			r.buf.WriteString(t.ColorBg(label, t.CellSelectedFg, t.CellSelectedBg))
		} else {
			r.moveTo(optRow, optCol)
			r.buf.WriteString("   ")
			r.buf.WriteString(t.Colorize(label, t.CellFixed))
		}

		r.moveTo(optRow+1, (width-displayWidth(descs[i]))/2)
		r.buf.WriteString(t.Colorize(descs[i], t.DimText))
	}

	// Nav hint
	hint := "Up/Down Select   Enter Confirm   Q Quit"
	r.moveTo(topRow+19, (width-displayWidth(hint))/2)
	r.buf.WriteString(t.Colorize(hint, t.DimText))

	r.flush()
}

// ─────────────────── Win Screen ───────────────────

// RenderWinScreen draws the victory screen.
func (r *Renderer) RenderWinScreen(gs *state.GameState) {
	r.buf.Reset()
	r.clearScreen()
	t := r.theme

	width, height, _ := GetTerminalSize()
	centerRow := height / 2

	elapsed := gs.ElapsedTime()
	min := int(elapsed.Minutes())
	sec := int(elapsed.Seconds()) % 60

	lines := []struct {
		text  string
		color string
	}{
		{"╔═══════════════════════════╗", t.AccentColor},
		{"║                           ║", t.AccentColor},
		{"║   *  P U Z Z L E  *      ║", t.AccentColor},
		{"║     S O L V E D !        ║", t.AccentColor},
		{"║                           ║", t.AccentColor},
		{"╚═══════════════════════════╝", t.AccentColor},
		{"", ""},
		{fmt.Sprintf("Difficulty: %s", gs.Difficulty.String()), t.CellFixed},
		{fmt.Sprintf("Time: %02d:%02d", min, sec), t.CellUser},
		{fmt.Sprintf("Errors: %d   Hints: %d", gs.Errors, gs.HintsUsed), t.StatusBar},
		{"", ""},
		{"R: New Game   Q: Quit", t.DimText},
	}

	rowStart := centerRow - len(lines)/2
	for i, l := range lines {
		if l.text == "" {
			continue
		}
		col := (width - displayWidth(l.text)) / 2
		if col < 1 {
			col = 1
		}
		r.moveTo(rowStart+i, col)
		r.buf.WriteString(t.Colorize(l.text, l.color))
	}

	r.flush()
}

// ─────────────────── Loading Screen ───────────────────

// RenderGenerating shows a "generating puzzle" message.
func (r *Renderer) RenderGenerating() {
	r.buf.Reset()
	r.clearScreen()
	t := r.theme

	width, height, _ := GetTerminalSize()
	msg := "Generating puzzle..."
	r.moveTo(height/2, (width-displayWidth(msg))/2)
	r.buf.WriteString(t.Colorize(msg, t.TitleColor))
	r.flush()
}
