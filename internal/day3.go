package internal

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/urfave/cli/v3"
)

func Day3(_ context.Context, cmd *cli.Command) error {
	path := cmd.StringArg("path")
	if path == "" {
		return fmt.Errorf("path is required")
	}
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	joltageSum := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("%s\n", line)
		if len(line) < 2 {
			return fmt.Errorf("bad line: %s", line)
		}

		maxPair := 0
		for i := 0; i < len(line)-1; i++ {
			firstChar := line[i]
			for j := i + 1; j < len(line); j++ {
				pairStr := string(firstChar) + string(line[j])
				pair, err := strconv.Atoi(pairStr)
				if err != nil {
					return fmt.Errorf("bad pair string: %s", pairStr)
				}
				if pair > maxPair {
					maxPair = pair
				}
			}
		}
		fmt.Printf("max pair %d\n", maxPair)
		joltageSum += maxPair
	}

	fmt.Printf("total output joltage %d\n", joltageSum)
	return nil
}
