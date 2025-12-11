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

type Point struct {
	col int
	row int
}

type Board struct {
	points       []Point
	grid         [][]bool
	borderPoints []Point
}

func newBoard() *Board {
	return &Board{}
}

func (b *Board) addPoint(p Point) {
	b.points = append(b.points, p)
}

func (b *Board) makeGrid() {
	maxRow := 0
	maxCol := 0
	for _, pt := range b.points {
		if pt.col > maxCol {
			maxCol = pt.col
		}
		if pt.row > maxRow {
			maxRow = pt.row
		}
	}
	b.grid = make([][]bool, maxRow+1)
	for row := 0; row <= maxRow; row++ {
		b.grid[row] = make([]bool, maxCol+1)
	}
	for _, pt := range b.points {
		b.grid[pt.row][pt.col] = true
	}
}

func (b *Board) fillBetween(pt1 Point, pt2 Point) {
	for curPoint := pt2; curPoint.row != pt1.row || curPoint.col != pt1.col; {
		b.borderPoints = append(b.borderPoints, Point{curPoint.col, curPoint.row})
		b.grid[curPoint.row][curPoint.col] = true
		if curPoint.row == pt1.row {
			if curPoint.col < pt1.col {
				curPoint.col++
			} else {
				curPoint.col--
			}
		} else if curPoint.row < pt1.row {
			curPoint.row++
		} else {
			curPoint.row--
		}
	}
}

func (b *Board) fillBorder() {
	fmt.Println("filling border")
	lastPoint := b.points[0]
	for i := 0; i < len(b.points); i++ {
		b.fillBetween(lastPoint, b.points[i])
		lastPoint = b.points[i]
	}
	b.fillBetween(lastPoint, b.points[0])
}

func (b *Board) findLargestFilled() int {
	fmt.Println("find largest")
	maxArea := 0
	for i, pt1 := range b.points {
		for j, pt2 := range b.points {
			if i == j {
				break
			}
			minCol := min(pt1.col, pt2.col)
			maxCol := max(pt1.col, pt2.col)
			minRow := min(pt1.row, pt2.row)
			maxRow := max(pt1.row, pt2.row)
			area := (maxCol - minCol + 1) * (maxRow - minRow + 1)
			if area < maxArea {
				continue
			}
			insideShape := true
			for _, borderPt := range b.borderPoints {
				if borderPt.col > minCol && borderPt.col < maxCol && borderPt.row > minRow && borderPt.row < maxRow {
					insideShape = false
					break
				}
			}
			if insideShape {
				fmt.Printf("found rect %v %v area %d\n", pt1, pt2, area)
				maxArea = area
			}
		}
	}
	return maxArea
}

func (b *Board) Print() {
	if len(b.grid) > 100 {
		return
	}
	for _, row := range b.grid {
		for _, val := range row {
			if val {
				fmt.Print("X")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Print("\n")
	}
}

func Day9(_ context.Context, cmd *cli.Command) error {
	path := cmd.StringArg("path")
	if path == "" {
		return fmt.Errorf("path is required")
	}
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	board := newBoard()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// fmt.Printf("%s\n", line)
		tokens := strings.Split(line, ",")
		if len(tokens) != 2 {
			return fmt.Errorf("failed to parse 2 ints from %s", line)
		}
		col, err := strconv.Atoi(tokens[0])
		if err != nil {
			return fmt.Errorf("bad int %s", tokens[0])
		}
		row, err := strconv.Atoi(tokens[1])
		if err != nil {
			return fmt.Errorf("bad int %s", tokens[1])
		}
		board.addPoint(Point{col, row})
	}

	board.makeGrid()
	board.Print()
	board.fillBorder()
	fmt.Println()
	board.Print()
	fmt.Println()
	board.Print()
	fmt.Printf("largest rect %d\n", board.findLargestFilled())

	return nil
}
