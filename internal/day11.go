package internal

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
)

type Graph struct {
	neighbors map[string][]string
}

func newGraph() *Graph {
	return &Graph{make(map[string][]string)}
}

func (g *Graph) Add(node string, neighbors []string) {
	g.neighbors[node] = neighbors
}

func (g *Graph) AllPathsCount(start string, end string) int {
	numFound := 0
	seen := make(map[string]bool)
	queue := []string{start}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		seen[cur] = true
		for _, neighbor := range g.neighbors[cur] {
			if neighbor == end {
				numFound++
				continue
			}
			if seen[neighbor] {
				continue
			}
			queue = append(queue, neighbor)
		}
	}
	return numFound
}

func Day11(_ context.Context, cmd *cli.Command) error {
	path := cmd.StringArg("path")
	if path == "" {
		return fmt.Errorf("path is required")
	}
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	graph := newGraph()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("%s\n", line)
		tokens := strings.Split(line, ":")
		if len(tokens) != 2 {
			return fmt.Errorf("bad line %s", line)
		}
		node := tokens[0]
		neighbors := strings.Split(tokens[1], " ")
		graph.Add(node, neighbors)
	}

	fmt.Println(graph)
	fmt.Println("num paths", graph.AllPathsCount("you", "out"))
	return nil
}
