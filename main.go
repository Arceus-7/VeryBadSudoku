package main

import (
	"fmt"
	"os"
	"time"

	"github.com/arceus-7/sudoku/game"
	"github.com/arceus-7/sudoku/state"
	"github.com/arceus-7/sudoku/ui"
)

func main() {
	// 1. Enable ANSI VT processing (Windows: sets console mode; Unix: no-op)
	if err := ui.EnableANSI(); err != nil {
		fmt.Fprintln(os.Stderr, "ANSI init failed:", err)
		os.Exit(1)
	}

	// 2. Switch to raw input mode
	restore, err := ui.EnableRawMode()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Raw mode failed:", err)
		os.Exit(1)
	}
	defer restore()

	// 3. Enter alternate screen, hide cursor
	ui.EnterAlternateScreen()
	ui.HideCursor()
	defer func() {
		ui.ShowCursor()
		ui.ExitAlternateScreen()
	}()

	// 4. Initialize theme with chroma16
	theme, err := ui.NewTheme("#FF6B35")
	if err != nil {
		ui.ShowCursor()
		ui.ExitAlternateScreen()
		restore()
		fmt.Fprintln(os.Stderr, "Theme init failed:", err)
		os.Exit(1)
	}

	renderer := ui.NewRenderer(theme)

	// Main application loop
	for {
		// Difficulty selection
		diff, quit := difficultySelect(renderer)
		if quit {
			return
		}

		// Show generating message
		renderer.RenderGenerating()

		// Generate puzzle and play
		gs := state.NewGame(diff)

		if playGame(renderer, gs) {
			return // user chose to quit entirely
		}
		// Otherwise, loop back to difficulty selection
	}
}

// difficultySelect shows the difficulty selection screen.
// Returns the selected difficulty and whether the user wants to quit.
func difficultySelect(renderer *ui.Renderer) (game.Difficulty, bool) {
	selected := 0
	difficulties := []game.Difficulty{game.Easy, game.Medium, game.Hard}

	renderer.RenderDifficultySelect(selected)

	for {
		key := ui.ReadKey()
		switch key {
		case ui.KeyUp:
			if selected > 0 {
				selected--
			}
		case ui.KeyDown:
			if selected < 2 {
				selected++
			}
		case ui.KeyEnter:
			return difficulties[selected], false
		case ui.KeyQuit, ui.KeyEscape:
			return game.Easy, true
		default:
			continue
		}

		renderer.RenderDifficultySelect(selected)
	}
}

// playGame runs the main game loop.
// Returns true if the user wants to quit the application entirely.
func playGame(renderer *ui.Renderer, gs *state.GameState) bool {
	// Render the initial board
	renderer.RenderGame(gs)

	// Set up a ticker for the timer display
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// We need to handle both key input and timer ticks.
	// Use a goroutine for key input and a channel.
	keyChan := make(chan ui.KeyEvent, 1)
	go func() {
		for {
			key := ui.ReadKey()
			keyChan <- key
			if key == ui.KeyQuit || key == ui.KeyEscape {
				return
			}
		}
	}()

	for {
		select {
		case key := <-keyChan:
			if gs.Won {
				// On win screen, only handle R (new game) or Q (quit)
				switch key {
				case ui.KeyNewGame:
					return false // go back to difficulty select
				case ui.KeyQuit, ui.KeyEscape:
					return true // quit app
				}
				continue
			}

			switch key {
			case ui.KeyUp:
				gs.MoveCursor(-1, 0)
			case ui.KeyDown:
				gs.MoveCursor(1, 0)
			case ui.KeyLeft:
				gs.MoveCursor(0, -1)
			case ui.KeyRight:
				gs.MoveCursor(0, 1)

			case ui.KeyNum1:
				gs.PlaceNumber(1)
			case ui.KeyNum2:
				gs.PlaceNumber(2)
			case ui.KeyNum3:
				gs.PlaceNumber(3)
			case ui.KeyNum4:
				gs.PlaceNumber(4)
			case ui.KeyNum5:
				gs.PlaceNumber(5)
			case ui.KeyNum6:
				gs.PlaceNumber(6)
			case ui.KeyNum7:
				gs.PlaceNumber(7)
			case ui.KeyNum8:
				gs.PlaceNumber(8)
			case ui.KeyNum9:
				gs.PlaceNumber(9)

			case ui.KeyDelete:
				gs.ClearCell()
			case ui.KeyNote:
				gs.ToggleNoteMode()
			case ui.KeyHint:
				gs.UseHint()
			case ui.KeyNewGame:
				return false // go back to difficulty select
			case ui.KeyQuit, ui.KeyEscape:
				return true // quit app
			}

			if gs.Won {
				renderer.RenderWinScreen(gs)
			} else {
				renderer.RenderGame(gs)
			}

		case <-ticker.C:
			// Re-render to update the timer display
			if !gs.Won {
				renderer.RenderGame(gs)
			}
		}
	}
}
