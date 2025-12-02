package internal

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/urfave/cli/v3"
)

func allSame(s []string) bool {
	for i := 1; i < len(s); i++ {
		if s[i] != s[0] {
			return false
		}
	}
	return true
}

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
			// fmt.Printf("Checking %s\n", valStr)
			for repeatWidth := 1; repeatWidth <= len(valStr)/2; repeatWidth++ {
				if len(valStr)%repeatWidth != 0 {
					continue
				}
				var pieces []string
				for i := 0; i < len(valStr); i += repeatWidth {
					pieces = append(pieces, valStr[i:i+repeatWidth])
				}
				if allSame(pieces) {
					fmt.Printf("Found invalid label: %d\n", val)
					invalidSum += val
					break
				}
			}
		}
	}

	fmt.Printf("Sum of invalid labels: %d\n", invalidSum)
	return nil
}
