package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

// Constants for grid size
const (
	Width  = 101
	Height = 103
)

// Robot represents a robot with position and velocity
type Robot struct {
	x  int
	y  int
	dx int
	dy int
}

// parseLine parses a line of input and returns a Robot
func parseLine(line string) (Robot, error) {
	// Example line: p=0,4 v=3,-3
	re := regexp.MustCompile(`p=(-?\d+),(-?\d+) v=(-?\d+),(-?\d+)`)
	matches := re.FindStringSubmatch(line)
	if matches == nil || len(matches) != 5 {
		return Robot{}, fmt.Errorf("invalid line format: %s", line)
	}
	x, err := strconv.Atoi(matches[1])
	if err != nil {
		return Robot{}, err
	}
	y, err := strconv.Atoi(matches[2])
	if err != nil {
		return Robot{}, err
	}
	dx, err := strconv.Atoi(matches[3])
	if err != nil {
		return Robot{}, err
	}
	dy, err := strconv.Atoi(matches[4])
	if err != nil {
		return Robot{}, err
	}
	return Robot{x: x, y: y, dx: dx, dy: dy}, nil
}

// readInput reads the input file and returns a slice of Robots
func readInput(filename string) ([]Robot, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	robots := []Robot{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		robot, err := parseLine(line)
		if err != nil {
			return nil, err
		}
		robots = append(robots, robot)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return robots, nil
}

// simulate simulates the movement of robots for n seconds
func simulate(robots []Robot, n int) []Robot {
	simulated := make([]Robot, len(robots))
	copy(simulated, robots)
	for i := 0; i < n; i++ {
		for j := 0; j < len(simulated); j++ {
			simulated[j].x = (simulated[j].x + simulated[j].dx) % Width
			if simulated[j].x < 0 {
				simulated[j].x += Width
			}
			simulated[j].y = (simulated[j].y + simulated[j].dy) % Height
			if simulated[j].y < 0 {
				simulated[j].y += Height
			}
		}
	}
	return simulated
}

// calculateSafetyFactor calculates the safety factor after 100 seconds
func calculateSafetyFactor(robots []Robot) int {
	// Define quadrants
	// Quadrant 1: left half (x < Width/2), top half (y < Height/2)
	// Quadrant 2: right half (x >= Width/2), top half (y < Height/2)
	// Quadrant 3: left half (x < Width/2), bottom half (y >= Height/2)
	// Quadrant 4: right half (x >= Width/2), bottom half (y >= Height/2)
	midX := Width / 2
	midY := Height / 2

	q1, q2, q3, q4 := 0, 0, 0, 0
	for _, robot := range robots {
		// Ignore robots exactly on the middle lines
		if robot.x == midX || robot.y == midY {
			continue
		}
		if robot.x < midX && robot.y < midY {
			q1++
		} else if robot.x >= midX && robot.y < midY {
			q2++
		} else if robot.x < midX && robot.y >= midY {
			q3++
		} else if robot.x >= midX && robot.y >= midY {
			q4++
		}
	}
	return q1 * q2 * q3 * q4
}

// boundingArea calculates the area of the bounding rectangle of all robots
func boundingArea(robots []Robot) (int, int, int, int) {
	minX, maxX := robots[0].x, robots[0].x
	minY, maxY := robots[0].y, robots[0].y
	for _, robot := range robots[1:] {
		if robot.x < minX {
			minX = robot.x
		}
		if robot.x > maxX {
			maxX = robot.x
		}
		if robot.y < minY {
			minY = robot.y
		}
		if robot.y > maxY {
			maxY = robot.y
		}
	}
	width := maxX - minX + 1
	height := maxY - minY + 1
	return width, height, minX, minY
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <input_file>")
		return
	}
	filename := os.Args[1]
	robots, err := readInput(filename)
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		return
	}

	// Part One
	startPart1 := time.Now()
	robotsAfter100 := simulate(robots, 100)
	safetyFactor := calculateSafetyFactor(robotsAfter100)
	durationPart1 := time.Since(startPart1)
	fmt.Printf("Part One Answer: %d\n", safetyFactor)
	fmt.Printf("Time taken for Part One: %v\n", durationPart1)

	// Part Two
	startPart2 := time.Now()
	// We'll simulate until the bounding area starts increasing
	// Assuming that the minimum area corresponds to the Easter egg
	prevArea := 1<<63 - 1 // Initialize to max int
	seconds := 0
	for {
		robots = simulate(robots, 1)
		seconds++
		width, height, _, _ := boundingArea(robots)
		area := width * height
		if area < prevArea {
			prevArea = area
		} else {
			// The previous second was the minimum
			seconds--
			break
		}
		// To prevent infinite loops, set a reasonable upper limit
		if seconds > 200000 {
			fmt.Println("Could not find the Easter egg within 200,000 seconds.")
			break
		}
	}
	durationPart2 := time.Since(startPart2)
	fmt.Printf("Part Two Answer: %d seconds\n", seconds)
	fmt.Printf("Time taken for Part Two: %v\n", durationPart2)
}
