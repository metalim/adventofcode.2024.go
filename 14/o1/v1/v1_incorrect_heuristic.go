/*
Напиши код на Go для решения задачи.
Входные данные в файле указываемом аргументом в командной строке.
Выведи ответ и время решения после решения каждой части.
Каждая часть должна решаться за несколько секунд максимум. Вторая часть задачи МОЖЕТ требовать особого подхода и не решаться перебором вариантов.
Если программа не сработает, я вставлю вывод и возможные комментарии. В ответ просто выдай исправленную версию.

https://adventofcode.com/2024/day/14

--- Day 14: Restroom Redoubt ---
One of The Historians needs to use the bathroom; fortunately, you know there's a bathroom near an unvisited location on their list, and so you're all quickly teleported directly to the lobby of Easter Bunny Headquarters.

Unfortunately, EBHQ seems to have "improved" bathroom security again after your last visit. The area outside the bathroom is swarming with robots!

To get The Historian safely to the bathroom, you'll need a way to predict where the robots will be in the future. Fortunately, they all seem to be moving on the tile floor in predictable straight lines.

You make a list (your puzzle input) of all of the robots' current positions (p) and velocities (v), one robot per line. For example:

p=0,4 v=3,-3
p=6,3 v=-1,-3
p=10,3 v=-1,2
p=2,0 v=2,-1
p=0,0 v=1,3
p=3,0 v=-2,-2
p=7,6 v=-1,-3
p=3,0 v=-1,-2
p=9,3 v=2,3
p=7,3 v=-1,2
p=2,4 v=2,-3
p=9,5 v=-3,-3
Each robot's position is given as p=x,y where x represents the number of tiles the robot is from the left wall and y represents the number of tiles from the top wall (when viewed from above). So, a position of p=0,0 means the robot is all the way in the top-left corner.

Each robot's velocity is given as v=x,y where x and y are given in tiles per second. Positive x means the robot is moving to the right, and positive y means the robot is moving down. So, a velocity of v=1,-2 means that each second, the robot moves 1 tile to the right and 2 tiles up.

The robots outside the actual bathroom are in a space which is 101 tiles wide and 103 tiles tall (when viewed from above). However, in this example, the robots are in a space which is only 11 tiles wide and 7 tiles tall.

The robots are good at navigating over/under each other (due to a combination of springs, extendable legs, and quadcopters), so they can share the same tile and don't interact with each other. Visually, the number of robots on each tile in this example looks like this:

1.12.......
...........
...........
......11.11
1.1........
.........1.
.......1...
These robots have a unique feature for maximum bathroom security: they can teleport. When a robot would run into an edge of the space they're in, they instead teleport to the other side, effectively wrapping around the edges. Here is what robot p=2,4 v=2,-3 does for the first few seconds:

Initial state:
...........
...........
...........
...........
..1........
...........
...........

After 1 second:
...........
....1......
...........
...........
...........
...........
...........

After 2 seconds:
...........
...........
...........
...........
...........
......1....
...........

After 3 seconds:
...........
...........
........1..
...........
...........
...........
...........

After 4 seconds:
...........
...........
...........
...........
...........
...........
..........1

After 5 seconds:
...........
...........
...........
.1.........
...........
...........
...........
The Historian can't wait much longer, so you don't have to simulate the robots for very long. Where will the robots be after 100 seconds?

In the above example, the number of robots on each tile after 100 seconds has elapsed looks like this:

......2..1.
...........
1..........
.11........
.....1.....
...12......
.1....1....
To determine the safest area, count the number of robots in each quadrant after 100 seconds. Robots that are exactly in the middle (horizontally or vertically) don't count as being in any quadrant, so the only relevant robots are:

..... 2..1.
..... .....
1.... .....

..... .....
...12 .....
.1... 1....
In this example, the quadrants contain 1, 3, 4, and 1 robot. Multiplying these together gives a total safety factor of 12.

Predict the motion of the robots in your list within a space which is 101 tiles wide and 103 tiles tall. What will the safety factor be after exactly 100 seconds have elapsed?

--- Part Two ---
During the bathroom break, someone notices that these robots seem awfully similar to ones built and used at the North Pole. If they're the same type of robots, they should have a hard-coded Easter egg: very rarely, most of the robots should arrange themselves into a picture of a Christmas tree.

What is the fewest number of seconds that must elapse for the robots to display the Easter egg?

*/

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
