// generated with o1-mini
// sometimes goes into loop

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// Direction represents movement in the maze
type Direction struct {
	dx, dy int
}

var directions = []Direction{
	{dx: 0, dy: -2}, // North
	{dx: 2, dy: 0},  // East
	{dx: 0, dy: 2},  // South
	{dx: -2, dy: 0}, // West
}

func main() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Parse flags
	outputFile := flag.String("o", "", "Output file for the maze")
	flag.Parse()

	// Parse positional arguments
	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Usage: go run maze.go [-o <outputFile>] <size>")
		os.Exit(1)
	}

	// Parse size
	size, err := strconv.Atoi(args[0])
	if err != nil || size < 5 {
		fmt.Println("Size must be an integer >= 5")
		os.Exit(1)
	}

	// Ensure size is odd
	if size%2 == 0 {
		size++
		fmt.Printf("Size is even. Increased to next odd number: %d\n", size)
	}

	// Initialize maze grid
	maze := make([][]rune, size)
	for i := range maze {
		maze[i] = make([]rune, size)
		for j := range maze[i] {
			maze[i][j] = '#'
		}
	}

	// Start DFS from a random cell
	startX := rand.Intn(size/2)*2 + 1
	startY := rand.Intn(size/2)*2 + 1
	maze[startY][startX] = '.'
	dfs(maze, startX, startY, size)

	// Place Start (S) and End (E) with minimum path length
	var sx, sy, ex, ey int
	minPathLength := size / 2 // Минимальная длина пути

	for {
		sx, sy = placeRandomCell(maze, size)
		ex, ey = placeRandomCell(maze, size)
		if sx == ex && sy == ey {
			continue // Разные точки
		}
		pathLength, exists := findPathLength(maze, size, sx, sy, ex, ey)
		if exists && pathLength >= minPathLength {
			break
		}
	}

	maze[sy][sx] = 'S'
	maze[ey][ex] = 'E'

	// Add thick walls
	err = addThickWalls(maze, size, 2, 3, 5, sx, sy, ex, ey)
	if err != nil {
		fmt.Println("Error adding thick walls:", err)
		os.Exit(1)
	}

	// Ensure path exists after adding walls
	pathLength, exists := findPathLength(maze, size, sx, sy, ex, ey)
	if !exists {
		fmt.Println("Failed to ensure a path exists between S and E after adding walls.")
		os.Exit(1)
	}
	if pathLength < minPathLength {
		fmt.Println("Path between S and E is too short after adding walls.")
		os.Exit(1)
	}

	// Generate maze string
	mazeStr := ""
	for _, row := range maze {
		for _, cell := range row {
			mazeStr += string(cell)
		}
		mazeStr += "\n"
	}

	// Output
	if *outputFile != "" {
		err := os.WriteFile(*outputFile, []byte(mazeStr), 0644)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			os.Exit(1)
		}
	} else {
		fmt.Print(mazeStr)
	}
}

// dfs performs depth-first search to generate maze passages
func dfs(maze [][]rune, x, y, size int) {
	dir := rand.Perm(len(directions))
	for _, i := range dir {
		d := directions[i]
		nx, ny := x+d.dx, y+d.dy
		if nx > 0 && nx < size-1 && ny > 0 && ny < size-1 && maze[ny][nx] == '#' {
			maze[ny][nx] = '.'
			maze[y+d.dy/2][x+d.dx/2] = '.'
			dfs(maze, nx, ny, size)
		}
	}
}

// addThickWalls adds a specified number of thick walls to the maze
func addThickWalls(maze [][]rune, size, minWalls, maxWalls, minLength, sx, sy, ex, ey int) error {
	numWalls := rand.Intn(maxWalls-minWalls+1) + minWalls
	attemptsPerWall := 100
	wallsAdded := 0

	for wallsAdded < numWalls {
		success := false
		for attempt := 0; attempt < attemptsPerWall; attempt++ {
			// Randomly choose horizontal or vertical
			orientation := rand.Intn(2) // 0: horizontal, 1: vertical
			length := rand.Intn(size/2) + minLength
			if length > size-4 {
				length = size - 4
			}
			var x, y int
			if orientation == 0 { // horizontal
				x = rand.Intn(size-length-2) + 1
				y = rand.Intn(size-2) + 1
			} else { // vertical
				x = rand.Intn(size-2) + 1
				y = rand.Intn(size-length-2) + 1
			}

			// Check if the wall overlaps with S or E
			overlaps := false
			for l := 0; l < length; l++ {
				px, py := x, y
				if orientation == 0 {
					px = x + l
				} else {
					py = y + l
				}
				// Check a 3x3 area around the wall cell
				for ty := py - 1; ty <= py+1; ty++ {
					for tx := px - 1; tx <= px+1; tx++ {
						if tx == sx && ty == sy || tx == ex && ty == ey {
							overlaps = true
							break
						}
					}
					if overlaps {
						break
					}
				}
				if overlaps {
					break
				}
			}
			if overlaps {
				continue
			}

			// Temporarily place the wall
			originalCells := make([][]rune, size)
			for i := range originalCells {
				originalCells[i] = make([]rune, size)
				copy(originalCells[i], maze[i])
			}

			for l := 0; l < length; l++ {
				px, py := x, y
				if orientation == 0 {
					px = x + l
				} else {
					py = y + l
				}
				// Make the wall thick (3x3)
				for ty := py - 1; ty <= py+1; ty++ {
					for tx := px - 1; tx <= px+1; tx++ {
						if tx >= 0 && tx < size && ty >= 0 && ty < size {
							maze[ty][tx] = '#'
						}
					}
				}
			}

			// Check if path still exists and path length is sufficient
			pathLength, exists := findPathLength(maze, size, sx, sy, ex, ey)
			if exists && pathLength >= (size/2) {
				success = true
				wallsAdded++
				break
			}

			// If not, revert the wall
			for i := range maze {
				copy(maze[i], originalCells[i])
			}
		}
		if !success {
			return fmt.Errorf("failed to place %d thick walls after %d attempts each", numWalls, attemptsPerWall)
		}
	}
	return nil
}

// placeRandomCell places a marker (S or E) on a random passage cell
func placeRandomCell(maze [][]rune, size int) (int, int) {
	for {
		x := rand.Intn(size)
		y := rand.Intn(size)
		if maze[y][x] == '.' {
			return x, y
		}
	}
}

// findPathLength checks if there is a path between (sx, sy) and (ex, ey) using BFS
// It returns the length of the shortest path and whether such a path exists
func findPathLength(maze [][]rune, size, sx, sy, ex, ey int) (int, bool) {
	visited := make([][]bool, size)
	for i := range visited {
		visited[i] = make([]bool, size)
	}
	type Node struct {
		x, y, depth int
	}
	queue := []Node{{sx, sy, 0}}
	visited[sy][sx] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if current.x == ex && current.y == ey {
			return current.depth, true
		}
		for _, d := range []Direction{
			{dx: 0, dy: -1}, // North
			{dx: 1, dy: 0},  // East
			{dx: 0, dy: 1},  // South
			{dx: -1, dy: 0}, // West
		} {
			nx, ny := current.x+d.dx, current.y+d.dy
			if nx >= 0 && nx < size && ny >= 0 && ny < size && !visited[ny][nx] && maze[ny][nx] != '#' {
				visited[ny][nx] = true
				queue = append(queue, Node{nx, ny, current.depth + 1})
			}
		}
	}
	return 0, false
}
