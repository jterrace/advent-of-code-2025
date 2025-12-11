package internal

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"

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
		l.vals[idx]++
	}
}

func (l *Joltages) PressCheck(button *Button, target *Joltages) bool {
	result := true
	for _, idx := range button.toggles {
		l.vals[idx]++
		if l.vals[idx] > target.vals[idx] {
			result = false
		}
	}
	return result
}

func (l *Joltages) Subtract(button Button) bool {
	result := true
	for _, idx := range button.toggles {
		l.vals[idx]--
		if l.vals[idx] < 0 {
			result = false
		}
	}
	return result
}

func (l *Joltages) Add(other *Joltages, target *Joltages) bool {
	for i, v := range other.vals {
		l.vals[i] += v
		if l.vals[i] > target.vals[i] {
			return false
		}
	}
	return true
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

type IndexIncrementer struct {
	indexes  []int
	bounds   [][]*Joltages
	done     bool
	sumCount int
}

func (incr *IndexIncrementer) IncrementFrom(start int) {
	for i := start; i >= 0; i-- {
		incr.indexes[i]++
		// if i == 4 {
		//	fmt.Println(incr.indexes)
		// }
		incr.sumCount++
		if incr.indexes[i] < len(incr.bounds[i]) {
			break
		}
		if i == 0 {
			incr.done = true
		}
		incr.sumCount -= incr.indexes[i]
		incr.indexes[i] = 0
	}
}

func (incr *IndexIncrementer) Increment(broke bool) {
	if broke {
		// fmt.Println("broke before", incr.indexes)
		incr.SkipFromBreak()
		// fmt.Println("broke after", incr.indexes)
		// time.Sleep(time.Second)
	} else {
		incr.IncrementFrom(len(incr.indexes) - 1)
	}
}

func (incr *IndexIncrementer) SkipFromBreak() {
	for i := len(incr.indexes) - 1; i >= 0; i-- {
		if incr.indexes[i] != 0 {
			incr.sumCount -= incr.indexes[i]
			incr.indexes[i] = 0
			incr.IncrementFrom(i - 1)
			break
		}
	}
}

func newIndexIncrementer(bounds [][]*Joltages) *IndexIncrementer {
	return &IndexIncrementer{make([]int, len(bounds)), bounds, false, 0}
}

func doJoltages(ch chan int, wg *sync.WaitGroup, machine *Machine) {
	defer wg.Done()
	fmt.Println(machine)
	perStep := make([][]*Joltages, len(machine.buttons))
	minCount := -1
	for buttonIdx, button := range machine.buttons {
		curJoltages := newJoltagesCount(len(machine.targetJoltages.vals))
		for stepCount := 0; ; stepCount++ {
			perStep[buttonIdx] = append(perStep[buttonIdx], &Joltages{slices.Clone(curJoltages.vals)})
			if !curJoltages.PressCheck(&button, machine.targetJoltages) {
				break
			}
		}
		fmt.Println(buttonIdx, perStep[buttonIdx])
	}

	numVoltages := len(machine.targetJoltages.vals)
	broke := false
	for incr := newIndexIncrementer(perStep); !incr.done; incr.Increment(broke) {
		broke = false
		// fmt.Println(incr.indexes)
		if minCount != -1 && incr.sumCount >= minCount {
			continue
		}

		allEqual := true
		for j := range numVoltages {
			joltageSum := 0
			for i := 0; i < len(incr.indexes); i++ {
				joltageSum += perStep[i][incr.indexes[i]].vals[j]
				if joltageSum > machine.targetJoltages.vals[j] {
					broke = true
					break
				}
			}
			if broke {
				break
			}
			if joltageSum != machine.targetJoltages.vals[j] {
				allEqual = false
				break
			}
		}
		if broke {
			continue
		}
		if allEqual {
			fmt.Println("found one at size", incr.sumCount)
			fmt.Println(incr.indexes)
			if minCount == -1 || incr.sumCount < minCount {
				minCount = incr.sumCount
			}
		}
	}
	fmt.Println("min count", minCount)
	ch <- minCount
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
	var wg sync.WaitGroup
	ch := make(chan int)

	numExpected := len(machines)
	for _, machine := range machines {
		wg.Add(1)
		go doJoltages(ch, &wg, &machine)
	}

	var receiveWg sync.WaitGroup
	receiveWg.Add(1)
	go func() {
		defer receiveWg.Done()
		received := 0
		for v := range ch {
			received++
			fmt.Println("got result!", v, received, "out of", numExpected)
			foundCounts = append(foundCounts, v)
		}
	}()

	wg.Wait()
	close(ch)
	receiveWg.Wait()

	fmt.Println("found counts", foundCounts)
	sum = 0
	for _, v := range foundCounts {
		sum += v
	}
	fmt.Printf("total sum %d\n", sum)

	return nil
}
