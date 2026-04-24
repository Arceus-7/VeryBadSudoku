# VeryBadSudoku

A terminal-based Sudoku game written in pure Go. No external game engines, no TUI frameworks -- just raw terminal manipulation with Unicode box-drawing and 24-bit color.

Runs on Linux, macOS, and Windows from a single static binary.

## Screenshot

```
╔═══════════════════════╗
║  V E R Y  B A D      ║
║    S U D O K U       ║
╚═══════════════════════╝

╔═══╤═══╤═══╦═══╤═══╤═══╦═══╤═══╤═══╗
║ 5 │ 3 │ · ║ · │ 7 │ · ║ · │ · │ · ║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║ 6 │ · │ · ║ 1 │ 9 │ 5 ║ · │ · │ · ║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║ · │ 9 │ 8 ║ · │ · │ · ║ · │ 6 │ · ║
╠═══╪═══╪═══╬═══╪═══╪═══╬═══╪═══╪═══╣
║ · │ · │ · ║ · │ 6 │ · ║ · │ · │ · ║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║ · │ · │ · ║ · │ · │ · ║ · │ · │ · ║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║ · │ · │ · ║ · │ 8 │ · ║ · │ · │ · ║
╠═══╪═══╪═══╬═══╪═══╪═══╬═══╪═══╪═══╣
║ · │ 6 │ · ║ · │ · │ · ║ 2 │ 8 │ · ║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║ · │ · │ · ║ 4 │ 1 │ 9 ║ · │ · │ 5 ║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║ · │ · │ · ║ · │ 8 │ · ║ · │ 7 │ 9 ║
╚═══╧═══╧═══╩═══╧═══╧═══╩═══╧═══╧═══╝
```

## Features

- Three difficulty levels: Easy, Medium, Hard
- Puzzle generation with guaranteed unique solutions
- Real-time conflict detection and highlighting
- Pencil mark / note mode
- Same-number highlighting across the grid
- Hint system that reveals the correct value
- Live game timer with error and hint tracking
- Centered, colored UI driven entirely by chroma16 palettes
- Alternate screen buffer for clean terminal restore on exit

## Requirements

- Go 1.21 or later
- Windows 10 1903+, macOS, or Linux
- A terminal with UTF-8 and ANSI color support

## Install

```bash
go install github.com/arceus-7/sudoku@latest
```

## Build from source

```bash
git clone https://github.com/arceus-7/sudoku.git
cd sudoku
go build -o sudoku .
./sudoku
```

Cross-compile for Windows from Linux/macOS:

```bash
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o sudoku.exe .
```

## Controls

| Key | Action |
|-----|--------|
| Arrow keys / WASD | Move cursor |
| 1-9 | Place number (or toggle pencil mark in note mode) |
| Backspace / Delete | Clear cell |
| N | Toggle pencil mark mode |
| H | Reveal correct number (hint) |
| R | Start a new game |
| Q | Quit |

## Project Structure

```
sudoku/
├── main.go              Entry point and game loop
├── main_windows.go      Windows UTF-8 codepage init
├── game/
│   ├── board.go         Cell and Board types
│   ├── difficulty.go    Difficulty levels
│   ├── generator.go     Puzzle generation with uniqueness check
│   ├── solver.go        Backtracking solver
│   └── validator.go     Row/column/box constraint validation
├── ui/
│   ├── input.go         Cross-platform key input
│   ├── renderer.go      Full-screen rendering with box drawing
│   ├── terminal_unix.go    Unix raw mode (x/term)
│   ├── terminal_windows.go Windows raw mode (x/sys/windows)
│   └── theme.go         Color theme (sole ANSI code source, via chroma16)
└── state/
    └── gamestate.go     Game state, cursor, hints, win detection
```

## Dependencies

| Package | Purpose |
|---------|---------|
| [chroma16](https://github.com/arceus-7/chroma16) | 16-color palette generation from a seed |
| [golang.org/x/term](https://pkg.go.dev/golang.org/x/term) | Unix terminal raw mode |
| [golang.org/x/sys](https://pkg.go.dev/golang.org/x/sys) | Windows console API access |

No TUI frameworks. No game engines. No CGO.

## Design Constraints

- All ANSI color codes are constructed exclusively in `ui/theme.go`
- All platform-specific code is isolated to `terminal_unix.go`, `terminal_windows.go`, and `main_windows.go`
- The game produces a single static binary with no runtime assets
- Puzzles are generated at startup with a backtracking solver and verified for unique solvability

## License

MIT. See [LICENSE](LICENSE).
