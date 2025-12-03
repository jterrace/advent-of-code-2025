package internal

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"slices"
	"strconv"

	"github.com/urfave/cli/v3"
)

var WIDTH = 12

func toIntArray(s string) ([]int, error) {
	ints := make([]int, len(s))
	for i, c := range s {
		v, err := strconv.Atoi(string(c))
		if err != nil {
			return nil, fmt.Errorf("bad character: %s", c)
		}
		ints[i] = v
	}
	return ints, nil
}

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
		if len(line) < 2 || len(line) < WIDTH {
			return fmt.Errorf("bad line: %s", line)
		}

		joltDigits, err := toIntArray(line)
		if err != nil {
			return err
		}

		maxJoltage := bytes.NewBufferString("")
		nextStart := 0
		for remaining := WIDTH; remaining > 0; remaining-- {
			// fmt.Printf("next start %d remaining %d\n", nextStart, remaining)
			checkBatch := joltDigits[nextStart : len(joltDigits)-remaining+1]
			maxValue := slices.Max(checkBatch)
			maxIndex := slices.Index(checkBatch, maxValue) + nextStart
			if maxIndex == -1 {
				return fmt.Errorf("somehow did not get max index from max value %d on line %s", maxValue, line)
			}
			maxJoltage.WriteByte(line[maxIndex])
			nextStart = maxIndex + 1
		}
		maxJoltageStr := maxJoltage.String()
		fmt.Printf("max voltage %s\n", maxJoltageStr)
		maxJoltageNum, err := strconv.Atoi(maxJoltageStr)
		if err != nil {
			return fmt.Errorf("somehow failed to parse back to int %s", maxJoltageStr)
		}
		joltageSum += maxJoltageNum
	}

	fmt.Printf("total output joltage %d\n", joltageSum)
	return nil
}
