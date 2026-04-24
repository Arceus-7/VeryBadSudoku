package game

// IsValidPlacement checks whether placing val at (row, col) would violate
// Sudoku constraints (row, column, and 3×3 box uniqueness).
// It ignores the cell's current value (useful for testing a hypothetical placement).
func IsValidPlacement(b *Board, row, col, val int) bool {
	if val == 0 {
		return true
	}

	// Check row
	for c := 0; c < 9; c++ {
		if c != col && b.Cells[row][c].Value == val {
			return false
		}
	}

	// Check column
	for r := 0; r < 9; r++ {
		if r != row && b.Cells[r][col].Value == val {
			return false
		}
	}

	// Check 3×3 box
	boxRow, boxCol := (row/3)*3, (col/3)*3
	for r := boxRow; r < boxRow+3; r++ {
		for c := boxCol; c < boxCol+3; c++ {
			if (r != row || c != col) && b.Cells[r][c].Value == val {
				return false
			}
		}
	}

	return true
}

// Validate checks the entire board for constraint violations and sets
// Cell.Invalid on every conflicting cell. Returns true if the board
// has no conflicts among filled cells.
func Validate(b *Board) bool {
	// Reset all invalid flags first
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			b.Cells[r][c].Invalid = false
		}
	}

	valid := true

	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			v := b.Cells[r][c].Value
			if v == 0 {
				continue
			}
			if !IsValidPlacement(b, r, c, v) {
				b.Cells[r][c].Invalid = true
				valid = false
			}
		}
	}

	return valid
}
