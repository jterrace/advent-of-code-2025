package internal

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/urfave/cli/v3"
)

func Day2(ctx context.Context, cmd *cli.Command) error {
	path := cmd.StringArg("path")
	if path == "" {
		return fmt.Errorf("path is required")
	}
	fileData, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	fileStr := string(fileData)

	invalidSum := 0
	idRanges := strings.Split(fileStr, ",")
	for _, idRange := range idRanges {
		idRange = strings.TrimSpace(idRange)
		if idRange == "" {
			continue
		}

		parts := strings.Split(idRange, "-")
		if len(parts) != 2 {
			return fmt.Errorf("invalid range: %s", idRange)
		}
		startRange, err := strconv.Atoi(parts[0])
		if err != nil {
			return fmt.Errorf("bad start range: %s", parts[0])
		}
		endRange, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("bad end range: %s", parts[1])
		}
		// fmt.Printf("Processing range %d-%d\n", startRange, endRange)

		for val := startRange; val <= endRange; val++ {
			valStr := strconv.Itoa(val)
			if len(valStr)%2 == 1 {
				continue
			}
			// fmt.Printf("Checking %s\n", valStr)
			mid := len(valStr) / 2
			if valStr[:mid] == valStr[mid:] {
				fmt.Printf("found invalid label %s\n", valStr)
				invalidSum += val
			}
		}
	}

	fmt.Printf("Sum of invalid labels: %d\n", invalidSum)
	return nil
}
