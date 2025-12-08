package internal

import (
	"bufio"
	"context"
	"fmt"
	"maps"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/ungerik/go3d/vec3"
	"github.com/urfave/cli/v3"
)

type PairDistance struct {
	point1   vec3.T
	point2   vec3.T
	distance float32
}

type ConnectedComponents struct {
	components    map[vec3.T]int
	nextComponent int
}

func newConnectedComponents() *ConnectedComponents {
	return &ConnectedComponents{make(map[vec3.T]int), 0}
}

func (c *ConnectedComponents) Add(point vec3.T) {
	c.components[point] = c.nextComponent
	c.nextComponent++
}

func (c *ConnectedComponents) Connect(point1 vec3.T, point2 vec3.T) {
	pt1Component := c.components[point1]
	pt2Component := c.components[point2]
	if pt1Component == pt2Component {
		return
	}

	fmt.Printf("merging components %d and %d\n", pt1Component, pt2Component)
	for pt := range c.components {
		if c.components[pt] == pt2Component {
			c.components[pt] = pt1Component
		}
	}
}

func (c *ConnectedComponents) SortByDistance() []PairDistance {
	var distances []PairDistance
	points := slices.Collect(maps.Keys(c.components))
	for i, point1 := range points {
		for j, point2 := range points {
			if j == i {
				break
			}
			distances = append(distances,
				PairDistance{point1, point2, vec3.Distance(&point1, &point2)})
		}
	}
	sort.Slice(distances, func(i, j int) bool {
		return distances[i].distance < distances[j].distance
	})
	return distances
}

func (c *ConnectedComponents) ByComponentCount() map[int]int {
	componentCount := make(map[int]int)
	for _, component := range c.components {
		componentCount[component]++
	}
	return componentCount
}

func Day8(_ context.Context, cmd *cli.Command) error {
	path := cmd.StringArg("path")
	if path == "" {
		return fmt.Errorf("path is required")
	}
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	components := newConnectedComponents()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, ",")
		if len(tokens) != 3 {
			return fmt.Errorf("bad point string %s", line)
		}
		x, err := strconv.Atoi(tokens[0])
		if err != nil {
			return fmt.Errorf("bad x %s", tokens[0])
		}
		y, err := strconv.Atoi(tokens[1])
		if err != nil {
			return fmt.Errorf("bad x %s", tokens[1])
		}
		z, err := strconv.Atoi(tokens[2])
		if err != nil {
			return fmt.Errorf("bad x %s", tokens[2])
		}
		components.Add(vec3.T{float32(x), float32(y), float32(z)})
		fmt.Printf("%s\n", line)
	}

	sortedByDistance := components.SortByDistance()
	// fmt.Println(sortedByDistance)
	for i, pairDistance := range sortedByDistance {
		if i == 1000 {
			break
		}
		components.Connect(pairDistance.point1, pairDistance.point2)
	}

	componentCount := components.ByComponentCount()
	sortedValues := slices.Collect(maps.Values(componentCount))
	sort.Ints(sortedValues)
	fmt.Println(sortedValues)
	fmt.Println(sortedValues[len(sortedValues)-1] * sortedValues[len(sortedValues)-2] * sortedValues[len(sortedValues)-3])
	return nil
}
