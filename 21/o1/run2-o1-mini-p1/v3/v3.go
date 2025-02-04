package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// Position represents coordinates on the keypad
type Position struct {
	x, y int
}

// State represents the current state in BFS
type State struct {
	pos     Position
	index   int
	presses int
}

// Keypad layout mappings
var numericKeypad = map[string]Position{
	"7": {0, 0}, "8": {1, 0}, "9": {2, 0},
	"4": {0, 1}, "5": {1, 1}, "6": {2, 1},
	"1": {0, 2}, "2": {1, 2}, "3": {2, 2},
	"0": {1, 3}, "A": {2, 3},
}

var directions = []Position{
	{0, -1}, // Up (^)
	{0, 1},  // Down (v)
	{-1, 0}, // Left (<)
	{1, 0},  // Right (>)
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
		sequenceLength := bfsSequenceLength(code, numericKeypad)
		numericValue := getNumericValue(code)
		if sequenceLength == -1 {
			fmt.Printf("Code: %s cannot be typed on the keypad.\n", code)
			continue
		}
		complexity := sequenceLength * numericValue
		totalComplexity += complexity
		elapsed := time.Since(start)
		fmt.Printf("Code: %s, Sequence Length: %d, Numeric Value: %d, Complexity: %d, Time: %v\n",
			code, sequenceLength, numericValue, complexity, elapsed)
	}

	fmt.Printf("Total Complexity: %d\n", totalComplexity)
}

// bfsSequenceLength computes the minimal number of button presses using BFS
func bfsSequenceLength(code string, keypad map[string]Position) int {
	targetSequence := strings.ToUpper(code)
	codeLength := len(targetSequence)

	// Precompute the target characters
	targetChars := []rune(targetSequence)

	// Initialize BFS
	startPos := keypad["A"]
	initialState := State{pos: startPos, index: 0, presses: 0}
	queue := []State{initialState}

	// Visited map to prevent revisiting states
	visited := make(map[string]bool)
	visited[stateKey(initialState)] = true

	for len(queue) > 0 {
		currentState := queue[0]
		queue = queue[1:]

		// If all characters are pressed, return the number of presses
		if currentState.index == codeLength {
			return currentState.presses
		}

		// Try all possible directional button presses
		for _, move := range directions {
			// Move the arm
			newX := currentState.pos.x + move.x
			newY := currentState.pos.y + move.y
			newPos, exists := positionAt(newX, newY, keypad)
			if !exists {
				// Invalid move, skip
				continue
			}
			newState := State{
				pos:     newPos,
				index:   currentState.index,
				presses: currentState.presses + 1,
			}
			key := stateKey(newState)
			if !visited[key] {
				visited[key] = true
				queue = append(queue, newState)
			}
		}

		// Press 'A' if at the correct position
		currentButton := getKeyAt(currentState.pos, keypad)
		if currentButton == string(targetChars[currentState.index]) {
			newState := State{
				pos:     currentState.pos,
				index:   currentState.index + 1,
				presses: currentState.presses + 1,
			}
			key := stateKey(newState)
			if !visited[key] {
				visited[key] = true
				queue = append(queue, newState)
			}
		}
	}

	// If the code cannot be typed
	return -1
}

// positionAt checks if a position exists on the keypad and returns the Position
func positionAt(x, y int, keypad map[string]Position) (Position, bool) {
	for _, pos := range keypad {
		if pos.x == x && pos.y == y {
			return pos, true
		}
	}
	return Position{}, false
}

// getKeyAt returns the key at a given position
func getKeyAt(pos Position, keypad map[string]Position) string {
	for key, p := range keypad {
		if p.x == pos.x && p.y == pos.y {
			return key
		}
	}
	return ""
}

// stateKey generates a unique key for a state
func stateKey(state State) string {
	return fmt.Sprintf("%d,%d,%d", state.pos.x, state.pos.y, state.index)
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
