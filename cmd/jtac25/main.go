package main

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	// This controls the maxprocs environment variable in container runtimes.
	// see https://martin.baillie.id/wrote/gotchas-in-the-go-network-packages-defaults/#bonus-gomaxprocs-containers-and-the-cfs
	"go.uber.org/automaxprocs/maxprocs"

	"github.com/jterrace/advent-of-code-2025/internal/log"
)

func main() {
	// Logger configuration
	logger := log.New(
		log.WithLevel(os.Getenv("LOG_LEVEL")),
		log.WithSource(),
	)

	if err := run(logger); err != nil {
		logger.ErrorContext(context.Background(), "an error occurred", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	ctx := context.Background()

	_, err := maxprocs.Set(maxprocs.Logger(func(s string, i ...interface{}) {
		logger.DebugContext(ctx, fmt.Sprintf(s, i...))
	}))
	if err != nil {
		return fmt.Errorf("setting max procs: %w", err)
	}

	logger.InfoContext(ctx, "Hello world!", slog.String("location", "world"))

	if len(os.Args) != 2 {
		return fmt.Errorf("provide input file as first argument")
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	dial := 50

	zero_count := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		if len(line) < 2 {
			return fmt.Errorf("bad line: %s", line)
		}
		direction := line[0]
		magnitude := 0
		switch direction {
		case 'L':
			magnitude = -1
		case 'R':
			magnitude = 1
		default:
			return fmt.Errorf("bad direction: %s", line)
		}
		value, err := strconv.Atoi(line[1:])
		if err != nil {
			return fmt.Errorf("bad amount: %s", line)
		}
		if value < 1 {
			return fmt.Errorf("bad amount: %d", value)
		}
		dial += magnitude * value
		dial = (dial%100 + 100) % 100
		fmt.Println(dial)
		if dial == 0 {
			zero_count++
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to scan file: %w", err)
	}

	fmt.Println(zero_count)
	return nil
}
