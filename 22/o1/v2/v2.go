package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

// Constants for pruning
const PruneModulo = 16777216

// Function to generate the next secret number
func nextSecret(secret int) int {
	// Step 1: Multiply by 64, mix, prune
	step1 := secret * 64
	secret ^= step1
	secret %= PruneModulo

	// Step 2: Divide by 32, round down, mix, prune
	step2 := secret / 32
	secret ^= step2
	secret %= PruneModulo

	// Step 3: Multiply by 2048, mix, prune
	step3 := secret * 2048
	secret ^= step3
	secret %= PruneModulo

	return secret
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <input_file>")
		return
	}

	inputFile := os.Args[1]

	// Read initial secret numbers from the input file
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	var initialSecrets []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		secret, err := strconv.Atoi(line)
		if err != nil {
			fmt.Printf("Invalid secret number '%s': %v\n", line, err)
			return
		}
		initialSecrets = append(initialSecrets, secret)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	// Part One
	startPartOne := time.Now()
	sumPartOne := 0
	for _, secret := range initialSecrets {
		current := secret
		for i := 0; i < 2000; i++ {
			current = nextSecret(current)
		}
		sumPartOne += current
	}
	durationPartOne := time.Since(startPartOne)
	fmt.Printf("Part One Answer: %d\nTime Taken: %v\n", sumPartOne, durationPartOne)

	// Part Two
	startPartTwo := time.Now()
	// Map to store sequence to total bananas
	sequenceBananas := make(map[[4]int]int)

	for _, secret := range initialSecrets {
		current := secret
		prices := make([]int, 2001) // Including initial price
		prices[0] = current % 10
		for i := 1; i <= 2000; i++ {
			current = nextSecret(current)
			prices[i] = current % 10
		}

		// Calculate changes
		changes := make([]int, 2000)
		for i := 1; i <= 2000; i++ {
			changes[i-1] = prices[i] - prices[i-1]
		}

		// Iterate through changes to find all possible sequences of 4
		seenSequences := make(map[[4]int]bool)
		for i := 0; i <= 2000-4; i++ {
			var seq [4]int
			for j := 0; j < 4; j++ {
				seq[j] = changes[i+j]
			}
			// If this is the first occurrence of the sequence for this buyer
			if !seenSequences[seq] {
				sequenceBananas[seq] += prices[i+4]
				seenSequences[seq] = true
			}
		}
	}

	// Find the sequence with the maximum total bananas
	maxBananas := 0
	// bestSequence удален, так как не используется
	for _, bananas := range sequenceBananas {
		if bananas > maxBananas {
			maxBananas = bananas
		}
	}

	durationPartTwo := time.Since(startPartTwo)
	fmt.Printf("Part Two Answer: %d\nTime Taken: %v\n", maxBananas, durationPartTwo)
}
