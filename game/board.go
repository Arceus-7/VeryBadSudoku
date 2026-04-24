package game

import "fmt"

// Cell represents a single cell in the Sudoku grid.
type Cell struct {
	Value   int      // 0 = empty, 1-9 = filled
	Fixed   bool     // true if part of the original puzzle (non-editable)
	Notes   [10]bool // pencil marks; indices 1-9 are used
	Invalid bool     // true if flagged by the validator as conflicting
}

// HasNotes returns true if the cell has any pencil marks set.
func (c *Cell) HasNotes() bool {
	for i := 1; i <= 9; i++ {
		if c.Notes[i] {
			return true
		}
	}
	return false
}

// Board represents a 9×9 Sudoku grid.
type Board struct {
	Cells [9][9]Cell
}

// NewBoard returns an empty 9×9 board with all cells zeroed.
func NewBoard() *Board {
	return &Board{}
}

// Get returns the cell at the given row and column.
func (b *Board) Get(row, col int) *Cell {
	return &b.Cells[row][col]
}

// Set places a value in the cell at (row, col).
// Returns an error if the cell is fixed.
func (b *Board) Set(row, col, val int) error {
	if b.Cells[row][col].Fixed {
		return fmt.Errorf("cell (%d,%d) is fixed", row, col)
	}
	b.Cells[row][col].Value = val
	// Clear notes when a value is placed
	if val != 0 {
		b.Cells[row][col].Notes = [10]bool{}
	}
	return nil
}

// Clear removes the value from a non-fixed cell.
func (b *Board) Clear(row, col int) error {
	if b.Cells[row][col].Fixed {
		return fmt.Errorf("cell (%d,%d) is fixed", row, col)
	}
	b.Cells[row][col].Value = 0
	return nil
}

// ToggleNote toggles a pencil mark for the given number in the cell.
func (b *Board) ToggleNote(row, col, num int) {
	if b.Cells[row][col].Fixed || num < 1 || num > 9 {
		return
	}
	b.Cells[row][col].Notes[num] = !b.Cells[row][col].Notes[num]
}

// ClearNotes removes all pencil marks from a cell.
func (b *Board) ClearNotes(row, col int) {
	b.Cells[row][col].Notes = [10]bool{}
}

// IsFull returns true when every cell has a non-zero value.
func (b *Board) IsFull() bool {
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if b.Cells[r][c].Value == 0 {
				return false
			}
		}
	}
	return true
}

// Clone returns a deep copy of the board.
func (b *Board) Clone() *Board {
	nb := &Board{}
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			nb.Cells[r][c] = b.Cells[r][c]
		}
	}
	return nb
}
