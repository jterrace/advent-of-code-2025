package internal

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/urfave/cli/v3"
)

type Graph struct {
	neighbors map[string][]string
	parents   map[string][]string
	nodes     []string
	indexMap  map[string]int
}

func newGraph() *Graph {
	return &Graph{make(map[string][]string), make(map[string][]string), nil, make(map[string]int)}
}

func (g *Graph) Add(node string, neighbors []string) {
	g.neighbors[node] = neighbors
	g.indexMap[node] = len(g.nodes)
	g.nodes = append(g.nodes, node)
	for _, neighbor := range neighbors {
		if !slices.Contains(g.parents[neighbor], node) {
			g.parents[neighbor] = append(g.parents[neighbor], node)
		}
	}
}

func (g *Graph) NodesTo(node string) map[string]bool {
	seen := make(map[string]bool)
	queue := []string{node}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		seen[cur] = true
		for _, parent := range g.parents[cur] {
			if !seen[parent] {
				queue = append(queue, parent)
				seen[parent] = true
			}
		}
	}
	return seen
}

func (g *Graph) AllPathsCount(start string, end string, dacReachable map[string]bool, fftReachable map[string]bool) int {
	fftIndex := g.indexMap["fft"]
	dacIndex := g.indexMap["dac"]
	numFound := 0
	queue := [][]int{{g.indexMap[start]}}
	for len(queue) > 0 {
		curPath := queue[0]
		queue = queue[1:]
		cur := g.nodes[curPath[len(curPath)-1]]
		if len(queue)%100000 == 0 && len(queue) > 0 {
			fmt.Println("queue len", len(queue))
			lastQueue := queue[len(queue)-1]
			path := make([]string, len(lastQueue))
			for i, idx := range lastQueue {
				path[i] = g.nodes[idx]
			}
			fmt.Println(path)
			// time.Sleep(time.Second)
		}
		for _, neighbor := range g.neighbors[cur] {
			if neighbor == "dac" && !slices.Contains(curPath, fftIndex) {
				continue
			}
			if neighbor == end {
				numFound++
				fmt.Println("found path", curPath)
				// fmt.Println(seen)
				// fmt.Println(queue)
				continue
			}
			if !dacReachable[neighbor] && !slices.Contains(curPath, dacIndex) {
				continue
			}
			if !fftReachable[neighbor] && !slices.Contains(curPath, fftIndex) {
				continue
			}
			neighborIndex := g.indexMap[neighbor]
			// fmt.Println("cur", cur, "neighbor", neighbor, "seen", seen)
			if slices.Contains(curPath, neighborIndex) {
				continue
			}
			nextPath := slices.Clone(curPath)
			nextPath = append(nextPath, neighborIndex)
			// fmt.Println(nextPath)
			// time.Sleep(time.Second)
			queue = append(queue, nextPath)
		}
	}
	return numFound
}

func (g *Graph) AllPaths(start string, end string, dacReachable map[string]bool) int {
	numFound := 0
	queue := [][]int{{g.indexMap[start]}}
	for len(queue) > 0 {
		curPath := queue[0]
		queue = queue[1:]
		cur := g.nodes[curPath[len(curPath)-1]]
		if len(queue)%100000 == 0 && len(queue) > 0 {
			// fmt.Println("queue len", len(queue), "fund paths", numFound)
			lastQueue := queue[len(queue)-1]
			path := make([]string, len(lastQueue))
			for i, idx := range lastQueue {
				path[i] = g.nodes[idx]
			}
			// fmt.Println(path)
			// time.Sleep(time.Second)
		}
		for _, neighbor := range g.neighbors[cur] {
			if dacReachable != nil && !dacReachable[neighbor] {
				continue
			}
			if neighbor == end {
				numFound++
				// fmt.Println("found path", curPath)
				// fmt.Println(queue)
				continue
			}
			neighborIndex := g.indexMap[neighbor]
			// fmt.Println("cur", cur, "neighbor", neighbor, "seen", seen)
			if slices.Contains(curPath, neighborIndex) {
				continue
			}
			nextPath := slices.Clone(curPath)
			nextPath = append(nextPath, neighborIndex)
			// fmt.Println(nextPath)
			// time.Sleep(time.Second)
			queue = append(queue, nextPath)
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
		neighbors := strings.Split(strings.TrimSpace(tokens[1]), " ")
		fmt.Println("node", node, "neighborslen", len(neighbors), "neighbors", neighbors)
		graph.Add(node, neighbors)
	}

	fmt.Println(graph)
	fmt.Println("total nodes", len(graph.nodes))
	dacReachable := graph.NodesTo("dac")
	fmt.Println("nodes that can reach dac", len(dacReachable))
	fmt.Println(dacReachable)
	fftReachable := graph.NodesTo("fft")
	fmt.Println("nodes that can reach fft", len(fftReachable))
	fmt.Println(fftReachable)
	//fmt.Println("num paths", graph.AllPathsCount("svr", "out", dacReachable, fftReachable))
	fmt.Println("paths from fft to dac", graph.AllPaths("fft", "dac", dacReachable))
	fmt.Println("paths from svr to fft", graph.AllPaths("svr", "fft", fftReachable))
	fmt.Println("paths from dac to out", graph.AllPaths("dac", "out", nil))
	return nil
}
