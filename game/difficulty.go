package game

// Difficulty represents the puzzle difficulty level.
type Difficulty int

const (
	// Easy removes ~35 cells from the complete grid.
	Easy Difficulty = iota
	// Medium removes ~45 cells from the complete grid.
	Medium
	// Hard removes ~52 cells from the complete grid.
	Hard
)

// String returns the human-readable name of the difficulty.
func (d Difficulty) String() string {
	switch d {
	case Easy:
		return "Easy"
	case Medium:
		return "Medium"
	case Hard:
		return "Hard"
	default:
		return "Unknown"
	}
}

// CellsToRemove returns how many cells to blank out for this difficulty.
func (d Difficulty) CellsToRemove() int {
	switch d {
	case Easy:
		return 35
	case Medium:
		return 45
	case Hard:
		return 52
	default:
		return 35
	}
}
