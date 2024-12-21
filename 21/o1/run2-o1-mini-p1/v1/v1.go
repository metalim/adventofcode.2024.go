package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type Position struct {
	x, y int
}

func main() {
	start := time.Now()
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

	// Define the numeric keypad layout
	keypad := map[string]Position{
		"7": {0, 0}, "8": {1, 0}, "9": {2, 0},
		"4": {0, 1}, "5": {1, 1}, "6": {2, 1},
		"1": {0, 2}, "2": {1, 2}, "3": {2, 2},
		"0": {1, 3}, "A": {2, 3},
	}

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
	elapsed := time.Since(start)
	fmt.Printf("Input parsed in %v\n", elapsed)

	totalComplexity := 0
	start = time.Now()
	for _, code := range codes {
		sequenceLength := computeSequenceLength(code, keypad)
		numericValue := getNumericValue(code)
		complexity := sequenceLength * numericValue
		totalComplexity += complexity
		fmt.Printf("Code: %s, Sequence Length: %d, Numeric Value: %d, Complexity: %d\n",
			code, sequenceLength, numericValue, complexity)
	}
	elapsed = time.Since(start)
	fmt.Printf("Processing completed in %v\n", elapsed)
	fmt.Printf("Total Complexity: %d\n", totalComplexity)
}

func computeSequenceLength(code string, keypad map[string]Position) int {
	// Initial position is 'A'
	currentPos := keypad["A"]
	totalPresses := 0

	for _, char := range code {
		target := string(char)
		targetPos, exists := keypad[target]
		if !exists {
			// Invalid character, skip
			continue
		}
		// Calculate moves: up/down/left/right
		dx := abs(targetPos.x - currentPos.x)
		dy := abs(targetPos.y - currentPos.y)
		totalPresses += dx + dy + 1 // moves + 'A' press
		currentPos = targetPos
	}
	return totalPresses
}

func getNumericValue(code string) int {
	num := 0
	for _, char := range code {
		if char >= '0' && char <= '9' {
			num = num*10 + int(char-'0')
		}
	}
	return num
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
