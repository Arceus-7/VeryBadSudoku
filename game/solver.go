package game

import (
	"math/rand"
)

// Solve fills the board in-place using backtracking.
// Returns true if the board is solvable.
func Solve(b *Board) bool {
	return solveBacktrack(b)
}

// SolveRandom fills the board using backtracking with randomized candidate
// ordering. Used by the generator to produce varied complete grids.
func SolveRandom(b *Board, rng *rand.Rand) bool {
	return solveRandomBacktrack(b, rng)
}

// CountSolutions counts the number of solutions for the board, stopping
// once it reaches the given limit. Used to verify unique solvability.
func CountSolutions(b *Board, limit int) int {
	count := 0
	countBacktrack(b, &count, limit)
	return count
}

// solveBacktrack is the deterministic backtracking solver.
func solveBacktrack(b *Board) bool {
	row, col, found := findEmpty(b)
	if !found {
		return true // all cells filled — solved
	}

	for num := 1; num <= 9; num++ {
		if IsValidPlacement(b, row, col, num) {
			b.Cells[row][col].Value = num
			if solveBacktrack(b) {
				return true
			}
			b.Cells[row][col].Value = 0
		}
	}

	return false
}

// solveRandomBacktrack is the randomized backtracking solver.
func solveRandomBacktrack(b *Board, rng *rand.Rand) bool {
	row, col, found := findEmpty(b)
	if !found {
		return true
	}

	// Randomized candidate order
	candidates := rng.Perm(9)
	for _, idx := range candidates {
		num := idx + 1
		if IsValidPlacement(b, row, col, num) {
			b.Cells[row][col].Value = num
			if solveRandomBacktrack(b, rng) {
				return true
			}
			b.Cells[row][col].Value = 0
		}
	}

	return false
}

// countBacktrack counts solutions via backtracking, stopping at limit.
func countBacktrack(b *Board, count *int, limit int) {
	if *count >= limit {
		return
	}

	row, col, found := findEmpty(b)
	if !found {
		*count++
		return
	}

	for num := 1; num <= 9; num++ {
		if *count >= limit {
			return
		}
		if IsValidPlacement(b, row, col, num) {
			b.Cells[row][col].Value = num
			countBacktrack(b, count, limit)
			b.Cells[row][col].Value = 0
		}
	}
}

// findEmpty returns the position of the first empty cell (row-major order).
func findEmpty(b *Board) (int, int, bool) {
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if b.Cells[r][c].Value == 0 {
				return r, c, true
			}
		}
	}
	return 0, 0, false
}
