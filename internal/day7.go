package internal

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func recurseLocation(row int, col int, grid [][]rune, numEnds *int) {
	if row == len(grid)-1 {
		*numEnds++
		return
	}

	for nextRow := row + 1; nextRow < len(grid); nextRow++ {
		if nextRow == len(grid)-1 {
			*numEnds++
			return
		}
		if grid[nextRow][col] == '^' {
			recurseLocation(nextRow, col-1, grid, numEnds)
			recurseLocation(nextRow, col+1, grid, numEnds)
			return
		}
	}
}

func Day7(_ context.Context, cmd *cli.Command) error {
	path := cmd.StringArg("path")
	if path == "" {
		return fmt.Errorf("path is required")
	}
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var lines [][]rune
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("%s\n", line)
		lines = append(lines, []rune(line))
		if len(line) != len(lines[0]) {
			return fmt.Errorf("found line of size %d doesn't match existing %d", len(line), len(lines[0]))
		}
	}

	sLocation := -1
	for col, c := range lines[0] {
		if c == 'S' {
			sLocation = col
			break
		}
	}
	if sLocation == -1 {
		return fmt.Errorf("could not find S in first line %v", lines[0])
	}

	var numEnds int
	recurseLocation(0, sLocation, lines, &numEnds)
	fmt.Printf("total ends %d\n", numEnds)
	return nil
}
