package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Machine represents a claw machine's configuration and prize location
type Machine struct {
	Ax, Ay int64 // Button A movements
	Bx, By int64 // Button B movements
	Px, Py int64 // Prize position
}

// parseMachines parses the input file and returns a slice of Machines
func parseMachines(filename string) ([]Machine, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var machines []Machine
	scanner := bufio.NewScanner(file)

	var lines []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			if len(lines) == 3 {
				machine, err := parseMachineBlock(lines)
				if err != nil {
					return nil, err
				}
				machines = append(machines, machine)
			}
			lines = []string{}
			continue
		}
		lines = append(lines, line)
	}
	// Handle last block if not followed by a blank line
	if len(lines) == 3 {
		machine, err := parseMachineBlock(lines)
		if err != nil {
			return nil, err
		}
		machines = append(machines, machine)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return machines, nil
}

// parseMachineBlock parses a block of three lines into a Machine
func parseMachineBlock(lines []string) (Machine, error) {
	var machine Machine
	var err error

	// Regex patterns to extract numbers
	buttonARegex := regexp.MustCompile(`Button A: X\+(\d+), Y\+(\d+)`)
	buttonBRegex := regexp.MustCompile(`Button B: X\+(\d+), Y\+(\d+)`)
	prizeRegex := regexp.MustCompile(`Prize: X=(\d+), Y=(\d+)`)

	// Parse Button A
	matches := buttonARegex.FindStringSubmatch(lines[0])
	if len(matches) != 3 {
		return machine, fmt.Errorf("invalid Button A format: %s", lines[0])
	}
	machine.Ax, err = strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return machine, err
	}
	machine.Ay, err = strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		return machine, err
	}

	// Parse Button B
	matches = buttonBRegex.FindStringSubmatch(lines[1])
	if len(matches) != 3 {
		return machine, fmt.Errorf("invalid Button B format: %s", lines[1])
	}
	machine.Bx, err = strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return machine, err
	}
	machine.By, err = strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		return machine, err
	}

	// Parse Prize
	matches = prizeRegex.FindStringSubmatch(lines[2])
	if len(matches) != 3 {
		return machine, fmt.Errorf("invalid Prize format: %s", lines[2])
	}
	machine.Px, err = strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return machine, err
	}
	machine.Py, err = strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		return machine, err
	}

	return machine, nil
}

// solvePartOne attempts to find xA and xB <=100 that satisfy the prize position
// Returns the minimum token cost if possible, else returns -1
func solvePartOne(machine Machine) int64 {
	minTokens := int64(-1)
	for xA := int64(0); xA <= 100; xA++ {
		// Calculate remaining X after pressing A xA times
		remainingX := machine.Px - machine.Ax*xA
		if remainingX < 0 {
			continue
		}
		// Check if Bx divides remainingX
		if machine.Bx == 0 {
			if remainingX != 0 {
				continue
			}
			// Bx is 0 and remainingX is 0, so xB can be any value
			// We need to satisfy the Y equation
			// Ax*xA + Bx*xB = Px => already satisfied
			// Ay*xA + By*xB = Py
			// Solve for xB: (Py - Ay*xA) / By
			if machine.By == 0 {
				if machine.Py-machine.Ay*xA != 0 {
					continue
				}
				// By is 0 and equation is satisfied for any xB
				xB := int64(0)
				tokens := 3*xA + 1*xB
				if minTokens == -1 || tokens < minTokens {
					minTokens = tokens
				}
				continue
			}
			if (machine.Py-machine.Ay*xA)%machine.By != 0 {
				continue
			}
			xB := (machine.Py - machine.Ay*xA) / machine.By
			if xB < 0 || xB > 100 {
				continue
			}
			tokens := 3*xA + 1*xB
			if minTokens == -1 || tokens < minTokens {
				minTokens = tokens
			}
			continue
		}
		if remainingX%machine.Bx != 0 {
			continue
		}
		xB := remainingX / machine.Bx
		if xB < 0 || xB > 100 {
			continue
		}
		// Check Y equation
		if machine.Ay*xA+machine.By*xB != machine.Py {
			continue
		}
		tokens := 3*xA + 1*xB
		if minTokens == -1 || tokens < minTokens {
			minTokens = tokens
		}
	}
	return minTokens
}

// solvePartTwo solves the problem without the xA and xB <=100 constraint
// Returns the minimum token cost if possible, else returns -1
func solvePartTwo(machine Machine) int64 {
	// Solve the system:
	// Ax * xA + Bx * xB = Px
	// Ay * xA + By * xB = Py

	det := machine.Ax*machine.By - machine.Ay*machine.Bx
	if det == 0 {
		// No unique solution
		return -1
	}

	// Calculate determinant for xA and xB
	detX := machine.Px*machine.By - machine.Py*machine.Bx
	detY := machine.Ax*machine.Py - machine.Ay*machine.Px

	// Check if det divides detX and detY
	if detX%det != 0 || detY%det != 0 {
		return -1
	}

	xA := detX / det
	xB := detY / det

	if xA < 0 || xB < 0 {
		return -1
	}

	// Check if the solution satisfies the equations
	if machine.Ax*xA+machine.Bx*xB != machine.Px || machine.Ay*xA+machine.By*xB != machine.Py {
		return -1
	}

	// Calculate tokens
	tokens := 3*xA + 1*xB
	return tokens
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <input_file>")
		return
	}

	filename := os.Args[1]
	machines, err := parseMachines(filename)
	if err != nil {
		fmt.Printf("Error parsing machines: %v\n", err)
		return
	}

	// Part One
	startPartOne := time.Now()
	var totalTokensPartOne int64 = 0
	var countPartOne int64 = 0
	for _, machine := range machines {
		tokens := solvePartOne(machine)
		if tokens != -1 {
			totalTokensPartOne += tokens
			countPartOne++
		}
	}
	durationPartOne := time.Since(startPartOne)
	fmt.Printf("Part One:\n")
	fmt.Printf("Total prizes won: %d\n", countPartOne)
	fmt.Printf("Minimum tokens required: %d\n", totalTokensPartOne)
	fmt.Printf("Time taken: %v\n\n", durationPartOne)

	// Part Two
	// Add 10000000000000 to Px and Py for each machine
	startPartTwo := time.Now()
	var totalTokensPartTwo int64 = 0
	var countPartTwo int64 = 0
	for _, machine := range machines {
		machine.Px += 10000000000000
		machine.Py += 10000000000000
		tokens := solvePartTwo(machine)
		if tokens != -1 {
			totalTokensPartTwo += tokens
			countPartTwo++
		}
	}
	durationPartTwo := time.Since(startPartTwo)
	fmt.Printf("Part Two:\n")
	fmt.Printf("Total prizes won: %d\n", countPartTwo)
	fmt.Printf("Minimum tokens required: %d\n", totalTokensPartTwo)
	fmt.Printf("Time taken: %v\n", durationPartTwo)
}
