package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Position represents a coordinate in the grid
type Position struct {
	x, y int
}

// Directions for movement: up, down, left, right
var directions = []Position{
	{0, -1}, // Up
	{0, 1},  // Down
	{-1, 0}, // Left
	{1, 0},  // Right
}

// Grid size
const (
	GRID_SIZE_X = 71
	GRID_SIZE_Y = 71
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide the input file as a command-line argument.")
		return
	}

	inputFile := os.Args[1]
	bytes, err := readInput(inputFile)
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		return
	}

	// Part One
	startPartOne := time.Now()
	corrupted := make([][]bool, GRID_SIZE_Y)
	for i := range corrupted {
		corrupted[i] = make([]bool, GRID_SIZE_X)
	}

	// Simulate first 1024 bytes
	for i := 0; i < len(bytes) && i < 1024; i++ {
		pos := bytes[i]
		if pos.x >= 0 && pos.x < GRID_SIZE_X && pos.y >= 0 && pos.y < GRID_SIZE_Y {
			corrupted[pos.y][pos.x] = true
		}
	}
	steps, err := bfs(corrupted, Position{0, 0}, Position{GRID_SIZE_X - 1, GRID_SIZE_Y - 1})
	elapsedPartOne := time.Since(startPartOne)
	if err != nil {
		fmt.Printf("Part One: No path found. Error: %v\nTime taken: %v\n", err, elapsedPartOne)
	} else {
		fmt.Printf("%d\n%v\n", steps, elapsedPartOne)
	}

	// Part Two
	startPartTwo := time.Now()
	corruptedPartTwo := make([][]bool, GRID_SIZE_Y)
	for i := range corruptedPartTwo {
		corruptedPartTwo[i] = make([]bool, GRID_SIZE_X)
	}

	var firstBlockingByte Position
	pathExists := true
	for _, pos := range bytes { // Changed from for i, pos := range bytes to for _, pos := range bytes
		if pos.x >= 0 && pos.x < GRID_SIZE_X && pos.y >= 0 && pos.y < GRID_SIZE_Y {
			corruptedPartTwo[pos.y][pos.x] = true
		}

		// Check if start or end is corrupted
		if corruptedPartTwo[0][0] || corruptedPartTwo[GRID_SIZE_Y-1][GRID_SIZE_X-1] {
			firstBlockingByte = pos
			pathExists = false
			break
		}

		// Check if path exists
		_, err := bfs(corruptedPartTwo, Position{0, 0}, Position{GRID_SIZE_X - 1, GRID_SIZE_Y - 1})
		if err != nil {
			firstBlockingByte = pos
			pathExists = false
			break
		}

		// Optional: Print progress every 100 bytes
		// if (i+1)%100 == 0 {
		// 	fmt.Printf("Processed %d bytes...\n", i+1)
		// }
	}

	elapsedPartTwo := time.Since(startPartTwo)
	if pathExists {
		fmt.Printf("No blocking byte found within the input.\n%v\n", elapsedPartTwo)
	} else {
		fmt.Printf("%d,%d\n%v\n", firstBlockingByte.x, firstBlockingByte.y, elapsedPartTwo)
	}
}

// readInput reads the input file and returns a slice of Positions
func readInput(filename string) ([]Position, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var bytes []Position
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line format: %s", line)
		}
		x, err1 := strconv.Atoi(parts[0])
		y, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("invalid coordinates in line: %s", line)
		}
		bytes = append(bytes, Position{x, y})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return bytes, nil
}

// bfs performs Breadth-First Search to find the shortest path
func bfs(grid [][]bool, start, end Position) (int, error) {
	if grid[start.y][start.x] || grid[end.y][end.x] {
		return -1, fmt.Errorf("start or end is blocked")
	}

	visited := make([][]bool, GRID_SIZE_Y)
	for i := range visited {
		visited[i] = make([]bool, GRID_SIZE_X)
	}
	queue := []Position{start}
	visited[start.y][start.x] = true

	steps := 0
	for len(queue) > 0 {
		nextQueue := []Position{}
		for _, current := range queue {
			if current == end {
				return steps, nil
			}
			for _, dir := range directions {
				newX := current.x + dir.x
				newY := current.y + dir.y
				if newX < 0 || newX >= GRID_SIZE_X || newY < 0 || newY >= GRID_SIZE_Y {
					continue
				}
				if !grid[newY][newX] && !visited[newY][newX] {
					visited[newY][newX] = true
					nextQueue = append(nextQueue, Position{newX, newY})
				}
			}
		}
		queue = nextQueue
		steps++
	}

	return -1, fmt.Errorf("no path found")
}
