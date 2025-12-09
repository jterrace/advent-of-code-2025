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
	points []Point
	grid   [][]bool
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

func (b *Board) isFilledBetween(pt1 Point, pt2 Point) bool {
	for curPoint := pt2; curPoint.row != pt1.row || curPoint.col != pt1.col; {
		if !b.grid[curPoint.row][curPoint.col] {
			return false
		}
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
	return true
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

type FillingState int

const (
	StateOutside FillingState = iota
	StateOpening
	StateInside
)

func (b *Board) Fill() {
	fmt.Println("filling")
	var state FillingState
	firstCol := -1
	for row := 0; row < len(b.grid); row++ {
		state = StateOutside
		for col := 0; col < len(b.grid[row]); col++ {
			v := b.grid[row][col]
			switch state {
			case StateOutside:
				if v {
					state = StateOpening
					firstCol = col
				}
			case StateOpening:
				if v {
					firstCol = col
				} else {
					state = StateInside
				}
			case StateInside:
				if v {
					state = StateOutside
					b.fillBetween(Point{firstCol, row}, Point{col, row})
				}
			}
		}
	}
}

func (b *Board) findLargestFilled() int {
	fmt.Println("find largest")
	maxArea := 0
	for i, pt1 := range b.points {
		for j, pt2 := range b.points {
			if i == j {
				break
			}
			pt3 := Point{pt1.col, pt2.row}
			pt4 := Point{pt2.col, pt1.row}
			// fmt.Println(pt1, pt2, pt3, pt4)
			if b.isFilledBetween(pt1, pt3) && b.isFilledBetween(pt3, pt2) && b.isFilledBetween(pt2, pt4) && b.isFilledBetween(pt4, pt1) {
				area := (pt2.col - pt1.col + 1) * (pt2.row - pt1.row + 1)
				area = max(area, -area)
				if area > maxArea {
					fmt.Printf("found filled rect with %v and %v area %d\n", pt1, pt2, area)
					maxArea = area
				}
			}
		}
	}
	return maxArea
}

func (b *Board) Print(name string) {
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
	board.Print("initial")
	board.fillBorder()
	fmt.Println()
	board.Print("borderFilled")
	board.Fill()
	fmt.Println()
	board.Print("filled")
	fmt.Printf("largest rect %d\n", board.findLargestFilled())

	return nil
}
