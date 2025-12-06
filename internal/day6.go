package internal

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/urfave/cli/v3"
)

func strArrayToIntArray(strs []string) ([]int, error) {
	ints := make([]int, len(strs))
	for i, s := range strs {
		v, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("bad string: %s", s)
		}
		ints[i] = v
	}
	return ints, nil
}

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

	finalSum := 0
	var accumulation [][]int = nil
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// fmt.Printf("%s\n", line)
		items := slices.DeleteFunc(strings.Split(line, " "), func(e string) bool {
			return e == ""
		})

		if items[0] == "*" || items[0] == "+" {
			fmt.Println(items)
			fmt.Println("found last line")
			if len(items) != len(accumulation) {
				return fmt.Errorf("found end of size %d which doesn't match existing %d", len(items), len(accumulation))
			}
			for i, operand := range items {
				if operand != "*" && operand != "+" {
					return fmt.Errorf("bad operand %s", operand)
				}
				output := 0
				if operand == "*" {
					output = 1
				}
				for _, v := range accumulation[i] {
					if operand == "+" {
						output += v
					} else {
						output *= v
					}
				}
				fmt.Printf("result %d\n", output)
				finalSum += output
			}
			break
		}

		ints, err := strArrayToIntArray(items)
		if err != nil {
			return err
		}
		fmt.Println(ints)
		if accumulation == nil {
			accumulation = make([][]int, len(ints))
		}
		if len(ints) != len(accumulation) {
			return fmt.Errorf("found line of size %d which doesn't match existing %d", len(ints), len(accumulation))
		}
		for i, v := range ints {
			accumulation[i] = append(accumulation[i], v)
		}
	}

	fmt.Printf("final sum %d\n", finalSum)
	return nil
}
