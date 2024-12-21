package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// Position represents coordinates on a keypad
type Position struct {
	x, y int
}

// Keypad represents the layout of a keypad
type Keypad struct {
	layout map[string]Position
	keys   map[Position]string
}

// NewDirectionalKeypad creates a directional keypad
func NewDirectionalKeypad() *Keypad {
	layout := map[string]Position{
		"^": {0, 0},
		"A": {1, 0},
		"<": {0, 1},
		"v": {1, 1},
		">": {2, 1},
	}
	keys := make(map[Position]string)
	for k, v := range layout {
		keys[v] = k
	}
	return &Keypad{layout: layout, keys: keys}
}

// NewNumericKeypad creates a numeric keypad
func NewNumericKeypad() *Keypad {
	layout := map[string]Position{
		"7": {0, 0}, "8": {1, 0}, "9": {2, 0},
		"4": {0, 1}, "5": {1, 1}, "6": {2, 1},
		"1": {0, 2}, "2": {1, 2}, "3": {2, 2},
		"0": {1, 3}, "A": {2, 3},
	}
	keys := make(map[Position]string)
	for k, v := range layout {
		keys[v] = k
	}
	return &Keypad{layout: layout, keys: keys}
}

// State represents the current state in BFS
type State struct {
	yourPos    Position // Your keypad's arm position
	robot1Pos  Position // Robot 1's keypad's arm position
	robot2Pos  Position // Robot 2's keypad's arm position
	codeIndex  int      // Current index in the code
	pressCount int      // Total button presses so far
}

// Queue represents a queue for BFS
type Queue []State

// Enqueue adds a state to the queue
func (q *Queue) Enqueue(s State) {
	*q = append(*q, s)
}

// Dequeue removes and returns the first state from the queue
func (q *Queue) Dequeue() State {
	s := (*q)[0]
	*q = (*q)[1:]
	return s
}

// isEmpty checks if the queue is empty
func (q *Queue) isEmpty() bool {
	return len(*q) == 0
}

func main() {
	startTotal := time.Now()

	if len(os.Args) < 2 {
		fmt.Println("Please provide the input file as an argument.")
		return
	}
	filePath := os.Args[1]
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Initialize keypads
	yourKeypad := NewDirectionalKeypad()
	robot1Keypad := NewDirectionalKeypad()
	robot2Keypad := NewDirectionalKeypad()
	numericKeypad := NewNumericKeypad()

	// Read codes from the input file
	scanner := bufio.NewScanner(file)
	var codes []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			codes = append(codes, line)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	elapsedTotal := time.Since(startTotal)
	fmt.Printf("Input parsed in %v\n", elapsedTotal)

	totalComplexity := 0

	for _, code := range codes {
		start := time.Now()
		sequenceLength := bfs(yourKeypad, robot1Keypad, robot2Keypad, numericKeypad, code)
		if sequenceLength == -1 {
			fmt.Printf("Code: %s cannot be typed on the keypad.\n", code)
			continue
		}
		numericValue := getNumericValue(code)
		complexity := sequenceLength * numericValue
		totalComplexity += complexity
		elapsed := time.Since(start)
		fmt.Printf("Code: %s, Sequence Length: %d, Numeric Value: %d, Complexity: %d, Time: %v\n",
			code, sequenceLength, numericValue, complexity, elapsed)
	}

	fmt.Printf("Total Complexity: %d\n", totalComplexity)
}

// bfs performs a breadth-first search to find the minimal number of button presses
func bfs(yourKeypad, robot1Keypad, robot2Keypad, numericKeypad *Keypad, code string) int {
	target := strings.ToUpper(code)
	codeRunes := []rune(target)
	codeLength := len(codeRunes)

	// Initialize BFS
	initialState := State{
		yourPos:    yourKeypad.layout["A"],
		robot1Pos:  robot1Keypad.layout["A"],
		robot2Pos:  robot2Keypad.layout["A"],
		codeIndex:  0,
		pressCount: 0,
	}

	queue := &Queue{}
	queue.Enqueue(initialState)

	visited := make(map[string]bool)
	visited[stateKey(initialState)] = true

	for !queue.isEmpty() {
		current := queue.Dequeue()

		// If all characters are pressed, return the press count
		if current.codeIndex == codeLength {
			return current.pressCount
		}

		// Try all possible button presses: directions and 'A'
		// Directions: '^', 'v', '<', '>'
		for direction := range yourKeypad.layout {
			if direction == "A" {
				continue
			}

			// Attempt to press a direction
			newYourPos, valid := moveDirection(current.yourPos, direction, yourKeypad)
			if !valid {
				continue
			}
			newRobot1Pos, valid1 := moveDirection(current.robot1Pos, direction, robot1Keypad)
			if !valid1 {
				continue
			}
			newRobot2Pos, valid2 := moveDirection(current.robot2Pos, direction, robot2Keypad)
			if !valid2 {
				continue
			}

			newState := State{
				yourPos:    newYourPos,
				robot1Pos:  newRobot1Pos,
				robot2Pos:  newRobot2Pos,
				codeIndex:  current.codeIndex,
				pressCount: current.pressCount + 1,
			}

			key := stateKey(newState)
			if !visited[key] {
				visited[key] = true
				queue.Enqueue(newState)
			}
		}

		// Attempt to press 'A'
		currentDesiredChar := string(codeRunes[current.codeIndex])

		// Robot2's current position should point to the desired character on numeric keypad
		currentRobot2Button := getKeyAt(current.robot2Pos, numericKeypad)
		if currentRobot2Button == currentDesiredChar {
			newState := State{
				yourPos:    current.yourPos,
				robot1Pos:  current.robot1Pos,
				robot2Pos:  current.robot2Pos,
				codeIndex:  current.codeIndex + 1,
				pressCount: current.pressCount + 1,
			}
			key := stateKey(newState)
			if !visited[key] {
				visited[key] = true
				queue.Enqueue(newState)
			}
		}
	}

	// If the code cannot be typed
	return -1
}

// moveDirection attempts to move the arm in the given direction on the specified keypad
func moveDirection(current Position, direction string, keypad *Keypad) (Position, bool) {
	newX, newY := current.x, current.y
	switch direction {
	case "^":
		newY--
	case "v":
		newY++
	case "<":
		newX--
	case ">":
		newX++
	default:
		return current, false
	}

	// Check if the new position exists on the keypad
	if _, exists := keypad.keys[Position{newX, newY}]; exists {
		return Position{newX, newY}, true
	}
	return current, false
}

// stateKey generates a unique key for a state
func stateKey(state State) string {
	return fmt.Sprintf("%d,%d,%d,%d,%d", state.yourPos.x, state.yourPos.y,
		state.robot1Pos.x, state.robot1Pos.y,
		state.robot2Pos.x, state.robot2Pos.y, state.codeIndex)
}

// getKeyAt returns the key at a given position on the keypad
func getKeyAt(pos Position, keypad *Keypad) string {
	if key, exists := keypad.keys[pos]; exists {
		return key
	}
	return ""
}

// getNumericValue extracts the numeric value from the code
func getNumericValue(code string) int {
	num := 0
	for _, char := range code {
		if char >= '0' && char <= '9' {
			num = num*10 + int(char-'0')
		}
	}
	return num
}
