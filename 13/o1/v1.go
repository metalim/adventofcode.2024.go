/*
Напиши код на Go для решения задачи.
Входные данные в файле указываемом аргументом в командной строке.
Выведи ответ и время решения после решения каждой части.
Каждая часть должна решаться за несколько секунд максимум. Вторая часть задачи МОЖЕТ требовать особого подхода и не решаться перебором вариантов.
Если программа не сработает, я вставлю вывод и возможные комментарии. В ответ просто выдай исправленную версию.

https://adventofcode.com/2024/day/13

--- Day 13: Claw Contraption ---
Next up: the lobby of a resort on a tropical island. The Historians take a moment to admire the hexagonal floor tiles before spreading out.

Fortunately, it looks like the resort has a new arcade! Maybe you can win some prizes from the claw machines?

The claw machines here are a little unusual. Instead of a joystick or directional buttons to control the claw, these machines have two buttons labeled A and B. Worse, you can't just put in a token and play; it costs 3 tokens to push the A button and 1 token to push the B button.

With a little experimentation, you figure out that each machine's buttons are configured to move the claw a specific amount to the right (along the X axis) and a specific amount forward (along the Y axis) each time that button is pressed.

Each machine contains one prize; to win the prize, the claw must be positioned exactly above the prize on both the X and Y axes.

You wonder: what is the smallest number of tokens you would have to spend to win as many prizes as possible? You assemble a list of every machine's button behavior and prize location (your puzzle input). For example:

Button A: X+94, Y+34
Button B: X+22, Y+67
Prize: X=8400, Y=5400

Button A: X+26, Y+66
Button B: X+67, Y+21
Prize: X=12748, Y=12176

Button A: X+17, Y+86
Button B: X+84, Y+37
Prize: X=7870, Y=6450

Button A: X+69, Y+23
Button B: X+27, Y+71
Prize: X=18641, Y=10279
This list describes the button configuration and prize location of four different claw machines.

For now, consider just the first claw machine in the list:

Pushing the machine's A button would move the claw 94 units along the X axis and 34 units along the Y axis.
Pushing the B button would move the claw 22 units along the X axis and 67 units along the Y axis.
The prize is located at X=8400, Y=5400; this means that from the claw's initial position, it would need to move exactly 8400 units along the X axis and exactly 5400 units along the Y axis to be perfectly aligned with the prize in this machine.
The cheapest way to win the prize is by pushing the A button 80 times and the B button 40 times. This would line up the claw along the X axis (because 80*94 + 40*22 = 8400) and along the Y axis (because 80*34 + 40*67 = 5400). Doing this would cost 80*3 tokens for the A presses and 40*1 for the B presses, a total of 280 tokens.

For the second and fourth claw machines, there is no combination of A and B presses that will ever win a prize.

For the third claw machine, the cheapest way to win the prize is by pushing the A button 38 times and the B button 86 times. Doing this would cost a total of 200 tokens.

So, the most prizes you could possibly win is two; the minimum tokens you would have to spend to win all (two) prizes is 480.

You estimate that each button would need to be pressed no more than 100 times to win a prize. How else would someone be expected to play?

Figure out how to win as many prizes as possible. What is the fewest tokens you would have to spend to win all possible prizes?

--- Part Two ---
As you go to win the first prize, you discover that the claw is nowhere near where you expected it would be. Due to a unit conversion error in your measurements, the position of every prize is actually 10000000000000 higher on both the X and Y axis!

Add 10000000000000 to the X and Y position of every prize. After making this change, the example above would now look like this:

Button A: X+94, Y+34
Button B: X+22, Y+67
Prize: X=10000000008400, Y=10000000005400

Button A: X+26, Y+66
Button B: X+67, Y+21
Prize: X=10000000012748, Y=10000000012176

Button A: X+17, Y+86
Button B: X+84, Y+37
Prize: X=10000000007870, Y=10000000006450

Button A: X+69, Y+23
Button B: X+27, Y+71
Prize: X=10000000018641, Y=10000000010279
Now, it is only possible to win a prize on the second and fourth claw machines. Unfortunately, it will take many more than 100 presses to do so.

Using the corrected prize coordinates, figure out how to win as many prizes as possible. What is the fewest tokens you would have to spend to win all possible prizes?

*/

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
