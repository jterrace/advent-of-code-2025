package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	// This controls the maxprocs environment variable in container runtimes.
	// see https://martin.baillie.id/wrote/gotchas-in-the-go-network-packages-defaults/#bonus-gomaxprocs-containers-and-the-cfs
	"go.uber.org/automaxprocs/maxprocs"

	"github.com/jterrace/advent-of-code-2025/internal"
	"github.com/jterrace/advent-of-code-2025/internal/log"
	"github.com/urfave/cli/v3"
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
	cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:  "day1",
				Usage: "run day 1",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name: "path",
					},
				},
				Action: internal.Day1,
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		return fmt.Errorf("failed to run command: %w", err)
	}

	return nil
}
