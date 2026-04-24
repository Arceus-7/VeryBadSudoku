package state

import (
	"time"

	"github.com/arceus-7/sudoku/game"
)

// GameState holds all mutable game state.
type GameState struct {
	Board      *game.Board
	Solution   *game.Board
	Difficulty game.Difficulty
	CursorRow  int
	CursorCol  int
	NoteMode   bool
	StartTime  time.Time
	Errors     int
	Won        bool
	HintsUsed  int
}

// NewGame generates a new puzzle and returns a fresh GameState.
func NewGame(diff game.Difficulty) *GameState {
	puzzle, solution := game.Generate(diff)
	return &GameState{
		Board:      puzzle,
		Solution:   solution,
		Difficulty: diff,
		CursorRow:  4, // start in center
		CursorCol:  4,
		StartTime:  time.Now(),
	}
}

// MoveCursor moves the cursor by (dr, dc), clamping to the grid bounds.
func (gs *GameState) MoveCursor(dr, dc int) {
	gs.CursorRow += dr
	gs.CursorCol += dc

	if gs.CursorRow < 0 {
		gs.CursorRow = 0
	}
	if gs.CursorRow > 8 {
		gs.CursorRow = 8
	}
	if gs.CursorCol < 0 {
		gs.CursorCol = 0
	}
	if gs.CursorCol > 8 {
		gs.CursorCol = 8
	}
}

// PlaceNumber places a number at the cursor position.
// In note mode, it toggles a pencil mark instead.
func (gs *GameState) PlaceNumber(num int) {
	if gs.Won {
		return
	}

	cell := gs.Board.Get(gs.CursorRow, gs.CursorCol)
	if cell.Fixed {
		return
	}

	if gs.NoteMode {
		gs.Board.ToggleNote(gs.CursorRow, gs.CursorCol, num)
		return
	}

	// Check against solution
	solCell := gs.Solution.Get(gs.CursorRow, gs.CursorCol)
	if num != solCell.Value {
		gs.Errors++
	}

	_ = gs.Board.Set(gs.CursorRow, gs.CursorCol, num)
	game.Validate(gs.Board)

	if gs.Board.IsFull() && game.Validate(gs.Board) {
		gs.Won = true
	}
}

// ClearCell removes the value from the cell at the cursor position.
func (gs *GameState) ClearCell() {
	if gs.Won {
		return
	}

	cell := gs.Board.Get(gs.CursorRow, gs.CursorCol)
	if cell.Fixed {
		return
	}

	_ = gs.Board.Clear(gs.CursorRow, gs.CursorCol)
	gs.Board.ClearNotes(gs.CursorRow, gs.CursorCol)
	game.Validate(gs.Board)
}

// ToggleNoteMode toggles between normal and pencil mark mode.
func (gs *GameState) ToggleNoteMode() {
	gs.NoteMode = !gs.NoteMode
}

// UseHint reveals one empty cell using the solution.
// Returns true if a hint was given, false if no empty cells remain.
func (gs *GameState) UseHint() bool {
	if gs.Won {
		return false
	}

	// Try cursor position first
	cell := gs.Board.Get(gs.CursorRow, gs.CursorCol)
	if cell.Value == 0 && !cell.Fixed {
		solVal := gs.Solution.Get(gs.CursorRow, gs.CursorCol).Value
		_ = gs.Board.Set(gs.CursorRow, gs.CursorCol, solVal)
		gs.Board.Cells[gs.CursorRow][gs.CursorCol].Fixed = true // mark as given
		gs.HintsUsed++
		game.Validate(gs.Board)
		if gs.Board.IsFull() && game.Validate(gs.Board) {
			gs.Won = true
		}
		return true
	}

	// Otherwise find any empty cell
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if gs.Board.Cells[r][c].Value == 0 {
				solVal := gs.Solution.Get(r, c).Value
				_ = gs.Board.Set(r, c, solVal)
				gs.Board.Cells[r][c].Fixed = true
				gs.HintsUsed++
				game.Validate(gs.Board)
				if gs.Board.IsFull() && game.Validate(gs.Board) {
					gs.Won = true
				}
				return true
			}
		}
	}

	return false
}

// ElapsedTime returns the duration since the game started.
func (gs *GameState) ElapsedTime() time.Duration {
	return time.Since(gs.StartTime)
}
