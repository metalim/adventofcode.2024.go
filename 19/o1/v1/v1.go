package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run program.go <input_file>")
		return
	}

	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read the first line: towel patterns
	var towelPatterns []string
	if scanner.Scan() {
		line := scanner.Text()
		towelPatterns = parsePatterns(line)
	} else {
		fmt.Println("Input file is empty or invalid.")
		return
	}

	// Skip blank lines
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			break
		}
	}

	// Read designs
	var designs []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			designs = append(designs, line)
		}
	}

	// Part One
	startPartOne := time.Now()
	partOneCount := 0
	for _, design := range designs {
		if canBeFormed(design, towelPatterns) {
			partOneCount++
		}
	}
	durationPartOne := time.Since(startPartOne)
	fmt.Printf("Part One: %d\nTime taken: %v\n", partOneCount, durationPartOne)

	// Part Two
	startPartTwo := time.Now()
	partTwoSum := 0
	for _, design := range designs {
		count := countWays(design, towelPatterns)
		partTwoSum += count
	}
	durationPartTwo := time.Since(startPartTwo)
	fmt.Printf("Part Two: %d\nTime taken: %v\n", partTwoSum, durationPartTwo)
}

// parsePatterns splits the first line by commas and trims spaces
func parsePatterns(line string) []string {
	parts := strings.Split(line, ",")
	var patterns []string
	for _, part := range parts {
		pattern := strings.TrimSpace(part)
		if pattern != "" {
			patterns = append(patterns, pattern)
		}
	}
	return patterns
}

// canBeFormed checks if the design can be formed by concatenating towel patterns
func canBeFormed(design string, patterns []string) bool {
	n := len(design)
	dp := make([]bool, n+1)
	dp[0] = true

	for i := 1; i <= n; i++ {
		for _, p := range patterns {
			l := len(p)
			if i >= l && design[i-l:i] == p && dp[i-l] {
				dp[i] = true
				break
			}
		}
	}

	return dp[n]
}

// countWays returns the number of ways the design can be formed by concatenating towel patterns
func countWays(design string, patterns []string) int {
	n := len(design)
	dp := make([]int, n+1)
	dp[0] = 1

	for i := 1; i <= n; i++ {
		for _, p := range patterns {
			l := len(p)
			if i >= l && design[i-l:i] == p {
				dp[i] += dp[i-l]
			}
		}
	}

	return dp[n]
}
