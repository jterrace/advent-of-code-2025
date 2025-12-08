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

	var lines [][]int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("%s\n", line)
		lines = append(lines, make([]int, len(line)))
		if len(line) != len(lines[0]) {
			return fmt.Errorf("found line of size %d doesn't match existing %d", len(line), len(lines[0]))
		}
		for col, c := range line {
			switch c {
			case 'S':
				lines[len(lines)-1][col] = 1
			case '^':
				lines[len(lines)-1][col] = -1
			default:
				lines[len(lines)-1][col] = 0
			}
		}
	}

	for row, line := range lines {
		if row == len(lines)-1 {
			break
		}
		for col, v := range line {
			if v <= 0 {
				continue
			}
			if lines[row+1][col] != -1 {
				lines[row+1][col] += v
				continue
			}
			lines[row+1][col-1] += v
			lines[row+1][col+1] += v
		}
	}

	fmt.Println(lines[len(lines)-1])
	sum := 0
	for _, v := range lines[len(lines)-1] {
		sum += v
	}
	fmt.Printf("total endings %d\n", sum)
	return nil
}
