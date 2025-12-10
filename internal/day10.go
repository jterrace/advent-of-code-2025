package internal

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/mowshon/iterium"
	"github.com/urfave/cli/v3"
)

type Button struct {
	toggles []int
}

func newButton(s string) (*Button, error) {
	if s[0] != '(' || s[len(s)-1] != ')' {
		return nil, fmt.Errorf("bad button section %s", s)
	}
	button := &Button{}
	for _, indexStr := range strings.Split(s[1:len(s)-1], ",") {
		index, err := strconv.Atoi(indexStr)
		if err != nil {
			return nil, fmt.Errorf("bad num %s", indexStr)
		}
		button.toggles = append(button.toggles, index)
	}

	return button, nil
}

type Lights struct {
	vals []bool
}

func (l *Lights) Press(button Button) {
	for _, idx := range button.toggles {
		l.vals[idx] = !l.vals[idx]
	}
}

func (l *Lights) PressAll(buttons []Button) {
	for _, button := range buttons {
		l.Press(button)
	}
}

func newLightsCount(count int) *Lights {
	vals := make([]bool, count)
	return &Lights{vals}
}

func newLights(lightStr string) (*Lights, error) {
	lights := &Lights{}
	if lightStr[0] != '[' || lightStr[len(lightStr)-1] != ']' {
		return nil, fmt.Errorf("bad lights section %s", lightStr)
	}
	lightStr = lightStr[1 : len(lightStr)-1]
	for _, c := range lightStr {
		val := false
		if c == '#' {
			val = true
		}
		lights.vals = append(lights.vals, val)
	}
	return lights, nil
}

type Joltages struct {
	vals []int
}

func (l *Joltages) Press(button Button) {
	for _, idx := range button.toggles {
		l.vals[idx] += 1
	}
}

func (l *Joltages) Subtract(button Button) bool {
	result := true
	for _, idx := range button.toggles {
		l.vals[idx] -= 1
		if l.vals[idx] < 0 {
			result = false
		}
	}
	return result
}

func (l *Joltages) PressAll(buttons []Button) {
	for _, button := range buttons {
		l.Press(button)
	}
}

func newJoltagesCount(count int) *Joltages {
	vals := make([]int, count)
	return &Joltages{vals}
}

func newJoltages(joltageStr string) (*Joltages, error) {
	joltages := &Joltages{}
	if joltageStr[0] != '{' || joltageStr[len(joltageStr)-1] != '}' {
		return nil, fmt.Errorf("bad joltages section %s", joltageStr)
	}
	joltageStr = joltageStr[1 : len(joltageStr)-1]
	// fmt.Println(joltageStr)
	joltagesStrs := strings.Split(joltageStr, ",")
	for _, joltageStr := range joltagesStrs {
		joltage, err := strconv.Atoi(joltageStr)
		if err != nil {
			return nil, fmt.Errorf("bad joltage num %s", joltageStr)
		}
		joltages.vals = append(joltages.vals, joltage)
	}
	return joltages, nil
}

type Machine struct {
	targetLights   *Lights
	buttons        []Button
	targetJoltages *Joltages
}

func newMachine(line string) (*Machine, error) {
	machine := &Machine{}
	tokens := strings.Split(line, " ")

	targetLights, err := newLights(tokens[0])
	if err != nil {
		return nil, err
	}
	machine.targetLights = targetLights

	joltages, err := newJoltages(tokens[len(tokens)-1])
	if err != nil {
		return nil, err
	}
	machine.targetJoltages = joltages

	for i := 1; i < len(tokens)-1; i++ {
		button, err := newButton(tokens[i])
		if err != nil {
			return nil, err
		}
		machine.buttons = append(machine.buttons, *button)
	}

	// fmt.Println(machine)
	return machine, nil
}

type TraversalCase struct {
	curJoltages Joltages
	numPresses  int
}

func Day10(_ context.Context, cmd *cli.Command) error {
	path := cmd.StringArg("path")
	if path == "" {
		return fmt.Errorf("path is required")
	}
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var machines []Machine
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("%s\n", line)
		machine, err := newMachine(line)
		if err != nil {
			return err
		}
		machines = append(machines, *machine)
	}

	var foundCounts []int
	for _, machine := range machines {
		fmt.Println(machine)
		foundPress := false
		// fmt.Println("Start machine", machine.targetLights)
		for numPress := 1; !foundPress; numPress++ {
			fmt.Printf("checking combinations of count %d\n", numPress)
			result := iterium.CombinationsWithReplacement(machine.buttons, numPress)
			for {
				buttonPresses, err := result.Next()
				if err != nil {
					break
				}
				lights := newLightsCount(len(machine.targetLights.vals))
				lights.PressAll(buttonPresses)
				// fmt.Println("checking", buttonPresses)
				// fmt.Println("result", lights.vals)
				if slices.Equal(lights.vals, machine.targetLights.vals) {
					// fmt.Printf("found combo %v\n", buttonPresses)
					foundPress = true
					foundCounts = append(foundCounts, numPress)
					break
				}
			}
		}
	}

	fmt.Println("found counts", foundCounts)
	sum := 0
	for _, v := range foundCounts {
		sum += v
	}
	fmt.Printf("total sum %d\n", sum)

	foundCounts = nil
	for _, machine := range machines {
		fmt.Println(machine)
		targetZeros := make([]int, len(machine.targetJoltages.vals))
		foundPress := false
		// fmt.Println("Start machine", machine.targetLights)
		for !foundPress {
			checkCount := 0
			startJoltage := Joltages{slices.Clone(machine.targetJoltages.vals)}
			start := TraversalCase{startJoltage, 0}
			queue := []TraversalCase{start}
			seen := make(map[string]bool)
			for cur := queue[0]; !foundPress && len(queue) > 0; {
				cur = queue[0]
				queue = queue[1:]
				if cur.numPresses > checkCount {
					checkCount = cur.numPresses
					fmt.Println("at press count", checkCount)
				}
				// fmt.Println("current", cur)
				for _, button := range machine.buttons {
					nextJoltages := Joltages{slices.Clone(cur.curJoltages.vals)}
					// fmt.Println(nextJoltages.vals)
					result := nextJoltages.Subtract(button)
					stupid := fmt.Sprint(nextJoltages.vals)
					if !result || seen[stupid] {
						continue
					}
					if slices.Equal(targetZeros, nextJoltages.vals) {
						fmt.Println("found zero!")
						foundPress = true
						foundCounts = append(foundCounts, cur.numPresses+1)
						break
					}
					seen[stupid] = true
					queue = append(queue, TraversalCase{nextJoltages, cur.numPresses + 1})
				}
			}
		}
	}

	fmt.Println("found counts", foundCounts)
	sum = 0
	for _, v := range foundCounts {
		sum += v
	}
	fmt.Printf("total sum %d\n", sum)

	return nil
}
