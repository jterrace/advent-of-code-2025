package internal

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

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

	splitCount := 0
	for row, line := range lines {
		if row == len(lines)-1 {
			break
		}
		for col, c := range line {
			switch c {
			case 'S':
				lines[row+1][col] = '|'
			case '|':
				if lines[row+1][col] != '^' {
					lines[row+1][col] = '|'
					continue
				}
				splitCount += 1
				lines[row+1][col-1] = '|'
				lines[row+1][col+1] = '|'
			}
		}
	}

	fmt.Printf("total splits %d\n", splitCount)
	return nil
}
