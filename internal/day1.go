package internal

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/urfave/cli/v3"
)

func getMagnitude(c rune) (int, error) {
	switch c {
	case 'L':
		return -1, nil
	case 'R':
		return 1, nil
	default:
		return 0, fmt.Errorf("bad direction: %c", c)
	}
}

func Day1(_ context.Context, cmd *cli.Command) error {
	path := cmd.StringArg("path")
	if path == "" {
		return fmt.Errorf("path is required")
	}
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	dial := 50

	crossCount := 0
	zeroCount := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("%d + %s\n", dial, line)
		if len(line) < 2 {
			return fmt.Errorf("bad line: %s", line)
		}
		magnitude, err := getMagnitude(rune(line[0]))
		if err != nil {
			return err
		}
		value, err := strconv.Atoi(line[1:])
		if err != nil {
			return fmt.Errorf("bad amount: %s", line)
		}
		if value < 1 {
			return fmt.Errorf("bad amount: %d", value)
		}
		cycles := value / 100
		if cycles > 0 {
			fmt.Printf("adding cross count %d\n", cycles)
		}
		crossCount += cycles
		if value >= 100 {
			value %= 100
			fmt.Printf("modded: %d\n", value)
		}
		if value != 0 {
			preDial := dial
			dial += magnitude * value
			if preDial != 0 && (dial > 99 || dial <= 0) {
				fmt.Println("Adding 1 to cross count")
				crossCount++
			}
			dial = (dial%100 + 100) % 100
		}
		fmt.Printf("dial post: %d\n", dial)
		if dial == 0 {
			zeroCount++
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to scan file: %w", err)
	}

	fmt.Println(zeroCount)
	fmt.Println(crossCount)
	return nil
}
