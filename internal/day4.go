package internal

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

type Grid struct {
	grid   [][]rune
	width  int
	height int
}

func NewGrid(width int, height int) *Grid {
	newgrid := make([][]rune, height)
	for i := range height {
		newgrid[i] = make([]rune, width)
	}
	return &Grid{newgrid, width, height}
}

func (g Grid) Set(row int, col int, v rune) {
	g.grid[row][col] = v
}

func (g Grid) Get(row int, col int) (rune, error) {
	if row < 0 || row > g.height-1 {
		return 0, fmt.Errorf("out of bounds row")
	}
	if col < 0 || col > g.width-1 {
		return 0, fmt.Errorf("out of bounds cow")
	}
	return g.grid[row][col], nil
}

type Adjacent struct {
	col int
	row int
}

var _ADJACENTS = [...]Adjacent{
	{-1, -1},
	{-1, 0},
	{-1, 1},
	{0, 1},
	{1, 1},
	{1, 0},
	{1, -1},
	{0, -1},
}

const FILLED = '@'
const MARKED = 'x'
const CLEARED = '.'

func (g Grid) GetAndMarkAccessible() int {
	numFound := 0
	for row := 0; row <= g.height; row++ {
		for col := 0; col <= g.width; col++ {
			val, _ := g.Get(row, col)
			if val != FILLED {
				continue
			}

			numPresent := 0
			for _, adjacent := range _ADJACENTS {
				val, err := g.Get(row+adjacent.row, col+adjacent.col)
				if err != nil {
					continue
				}
				if val == FILLED || val == MARKED {
					numPresent++
				}
			}
			if numPresent < 4 {
				numFound++
				g.Set(row, col, MARKED)
			}
		}
	}
	fmt.Printf("Found %d total\n", numFound)
	return numFound
}

func (g Grid) ClearMarked() {
	for row := 0; row <= g.height; row++ {
		for col := 0; col <= g.width; col++ {
			val, _ := g.Get(row, col)
			if val == MARKED {
				g.Set(row, col, CLEARED)
			}
		}
	}
}

func Day4(_ context.Context, cmd *cli.Command) error {
	path := cmd.StringArg("path")
	if path == "" {
		return fmt.Errorf("path is required")
	}
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var lines []string
	width := -1
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if width == -1 {
			width = len(line)
		}
		if len(line) != width {
			return fmt.Errorf("line %s does not match width %d", line, width)
		}
		fmt.Printf("%s\n", line)
		lines = append(lines, line)
	}

	grid := NewGrid(width, len(lines))
	for row, line := range lines {
		for col, c := range line {
			grid.grid[row][col] = c
		}
	}

	totalRemoved := 0
	for numFound := grid.GetAndMarkAccessible(); numFound > 0; numFound = grid.GetAndMarkAccessible() {
		grid.ClearMarked()
		totalRemoved += numFound
	}

	fmt.Printf("Removed %d total\n", totalRemoved)
	return nil
}
