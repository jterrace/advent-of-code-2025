package internal

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/urfave/cli/v3"
)

func Day6(_ context.Context, cmd *cli.Command) error {
	path := cmd.StringArg("path")
	if path == "" {
		return fmt.Errorf("path is required")
	}
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var accumulation []strings.Builder = nil
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("%s\n", line)
		if accumulation == nil {
			accumulation = make([]strings.Builder, len(line))
		}

		if len(line) != len(accumulation) {
			return fmt.Errorf("found line of size %d doesn't match existing %d", len(line), len(accumulation))
		}

		for i, c := range line {
			accumulation[i].WriteRune(c)
		}
	}

	totalSum := 0
	var curAccum []int
	for col := len(accumulation) - 1; col >= 0; col-- {
		colStr := strings.TrimSpace(accumulation[col].String())
		if len(colStr) == 0 {
			continue
		}
		fmt.Println(colStr)
		operator := colStr[len(colStr)-1]
		if operator == '+' || operator == '*' {
			colStr = strings.TrimSpace(colStr[:len(colStr)-1])
		}

		colNum, err := strconv.Atoi(colStr)
		if err != nil {
			return fmt.Errorf("invalid num %s", colStr)
		}
		curAccum = append(curAccum, colNum)

		if operator == '+' || operator == '*' {
			sectionResult := 0
			if operator == '*' {
				sectionResult = 1
			}
			for _, v := range curAccum {
				if operator == '+' {
					sectionResult += v
				} else {
					sectionResult *= v
				}
			}
			fmt.Printf("section result = %d\n", sectionResult)
			totalSum += sectionResult

			curAccum = nil
		}
	}

	fmt.Printf("total sum %d\n", totalSum)
	return nil
}
