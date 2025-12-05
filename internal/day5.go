package internal

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/andrewwphillips/rangeset"
	"github.com/urfave/cli/v3"
)

func parseRange(line string) (uint64, uint64, error) {
	strNums := strings.Split(line, "-")
	if len(strNums) != 2 {
		return 0, 0, fmt.Errorf("incorrect range numbers %s", line)
	}
	startRange, err := strconv.ParseUint(strNums[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("error parsing start range %w", err)
	}

	endRange, err := strconv.ParseUint(strNums[1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("error parsing end range %w", err)
	}

	return startRange, endRange, nil
}

func Day5(_ context.Context, cmd *cli.Command) error {
	path := cmd.StringArg("path")
	if path == "" {
		return fmt.Errorf("path is required")
	}
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	numFresh := uint64(0)
	freshSet := rangeset.Make[uint64]()
	testingMode := false
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			testingMode = true
			continue
		}

		if testingMode {
			testNum, err := strconv.ParseUint(line, 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse %s", line)
			}
			hasit := freshSet.Contains(testNum)
			fmt.Printf("test %d = %t\n", testNum, hasit)
			if hasit {
				numFresh++
			}
		} else {
			startRange, endRange, err := parseRange(line)
			if err != nil {
				return err
			}
			fmt.Printf("start %d end %d\n", startRange, endRange)
			freshSet.AddRange(startRange, endRange+1)
		}
	}

	fmt.Printf("total fresh: %d\n", numFresh)

	fmt.Printf("%s\n", freshSet.String())

	numFresh = 0
	for v := range freshSet.SpansSeq() {
		fmt.Printf("start %d end %d\n", v.Bot, v.Top)
		numFresh += (v.Top - v.Bot)
	}
	fmt.Printf("total fresh: %d\n", numFresh)

	return nil
}
