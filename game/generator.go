package game

import (
	"math/rand"
	"time"
)

// Generate creates a new Sudoku puzzle at the given difficulty.
// Returns the puzzle board (with blanks) and the complete solution.
func Generate(diff Difficulty) (*Board, *Board) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Step 1: Generate a complete, valid grid
	solution := NewBoard()
	SolveRandom(solution, rng)

	// Step 2: Clone it as the puzzle, then remove cells
	puzzle := solution.Clone()

	toRemove := diff.CellsToRemove()
	removed := 0

	// Build a list of all cell positions for random removal
	type pos struct{ r, c int }
	positions := make([]pos, 0, 81)
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			positions = append(positions, pos{r, c})
		}
	}

	// Shuffle positions
	rng.Shuffle(len(positions), func(i, j int) {
		positions[i], positions[j] = positions[j], positions[i]
	})

	for _, p := range positions {
		if removed >= toRemove {
			break
		}

		r, c := p.r, p.c
		if puzzle.Cells[r][c].Value == 0 {
			continue
		}

		// Save current value
		saved := puzzle.Cells[r][c].Value
		puzzle.Cells[r][c].Value = 0

		// Check that the puzzle still has a unique solution
		check := puzzle.Clone()
		if CountSolutions(check, 2) != 1 {
			// Restore — removing this cell creates ambiguity
			puzzle.Cells[r][c].Value = saved
			continue
		}

		removed++
	}

	// Mark remaining filled cells as fixed (given clues)
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if puzzle.Cells[r][c].Value != 0 {
				puzzle.Cells[r][c].Fixed = true
			}
		}
	}

	// Also mark solution cells as fixed for consistency
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			solution.Cells[r][c].Fixed = true
		}
	}

	return puzzle, solution
}
