package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"time"
)

// Direction constants
const (
	East = iota
	South
	West
	North
)

// State represents the current position, direction, and accumulated score
type State struct {
	x, y      int
	direction int
	score     int
}

// PriorityQueue implements heap.Interface and holds States
type PriorityQueue []State

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].score < pq[j].score
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(State))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// Point represents a coordinate on the map
type Point struct {
	x, y int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <input_file>")
		return
	}

	filename := os.Args[1]

	// Read the map
	grid, start, end, err := readMap(filename)
	if err != nil {
		fmt.Println("Error reading map:", err)
		return
	}

	// Part 1
	startTime := time.Now()
	minScore, prev := dijkstra(grid, start, end)
	part1Time := time.Since(startTime)
	if minScore == -1 {
		fmt.Println("Part 1: No path found")
	} else {
		fmt.Printf("Part 1: Lowest score = %d\nTime taken: %v\n", minScore, part1Time)
	}

	// Part 2
	startTime = time.Now()
	bestPathTiles := findBestPathTiles(prev, end, grid)
	part2Time := time.Since(startTime)
	fmt.Printf("Part 2: Number of tiles on best paths = %d\nTime taken: %v\n", bestPathTiles, part2Time)
}

// readMap reads the input file and returns the grid, start and end points
func readMap(filename string) ([][]rune, Point, Point, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, Point{}, Point{}, err
	}
	defer file.Close()

	var grid [][]rune
	var start, end Point
	scanner := bufio.NewScanner(file)
	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		row := []rune(line)
		for x, char := range row {
			if char == 'S' {
				start = Point{x, y}
			} else if char == 'E' {
				end = Point{x, y}
			}
		}
		grid = append(grid, row)
		y++
	}

	if err := scanner.Err(); err != nil {
		return nil, Point{}, Point{}, err
	}

	return grid, start, end, nil
}

// dijkstra performs Dijkstra's algorithm to find the minimum score path
func dijkstra(grid [][]rune, start, end Point) (int, map[Point][]State) {
	rows := len(grid)
	cols := len(grid[0])

	// Initialize distance map
	dist := make(map[Point]map[int]int)
	prev := make(map[Point][]State)

	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			p := Point{x, y}
			dist[p] = make(map[int]int)
			for d := 0; d < 4; d++ {
				dist[p][d] = int(^uint(0) >> 1) // Infinity
			}
		}
	}

	// Initialize priority queue
	pq := &PriorityQueue{}
	heap.Init(pq)
	initialState := State{start.x, start.y, East, 0}
	heap.Push(pq, initialState)
	dist[start][East] = 0

	// Directions: East, South, West, North
	dirs := []Point{
		{1, 0},  // East
		{0, 1},  // South
		{-1, 0}, // West
		{0, -1}, // North
	}

	for pq.Len() > 0 {
		current := heap.Pop(pq).(State)
		currentPoint := Point{current.x, current.y}

		// If reached end
		if currentPoint == end {
			continue
		}

		// Explore possible actions: move forward, turn left, turn right

		// Move Forward
		nextX := current.x + dirs[current.direction].x
		nextY := current.y + dirs[current.direction].y
		if isValid(grid, nextX, nextY) {
			nextPoint := Point{nextX, nextY}
			newScore := current.score + 1
			if newScore < dist[nextPoint][current.direction] {
				dist[nextPoint][current.direction] = newScore
				heap.Push(pq, State{nextX, nextY, current.direction, newScore})
				prev[nextPoint] = append(prev[nextPoint], current)
			} else if newScore == dist[nextPoint][current.direction] {
				prev[nextPoint] = append(prev[nextPoint], current)
			}
		}

		// Turn Left
		newDir := (current.direction + 3) % 4
		newScore := current.score + 1000
		if newScore < dist[currentPoint][newDir] {
			dist[currentPoint][newDir] = newScore
			heap.Push(pq, State{current.x, current.y, newDir, newScore})
			prev[currentPoint] = append(prev[currentPoint], current)
		} else if newScore == dist[currentPoint][newDir] {
			prev[currentPoint] = append(prev[currentPoint], current)
		}

		// Turn Right
		newDir = (current.direction + 1) % 4
		newScore = current.score + 1000
		if newScore < dist[currentPoint][newDir] {
			dist[currentPoint][newDir] = newScore
			heap.Push(pq, State{current.x, current.y, newDir, newScore})
			prev[currentPoint] = append(prev[currentPoint], current)
		} else if newScore == dist[currentPoint][newDir] {
			prev[currentPoint] = append(prev[currentPoint], current)
		}
	}

	// Find the minimum score to reach end
	minScore := int(^uint(0) >> 1)
	for d := 0; d < 4; d++ {
		if dist[end][d] < minScore {
			minScore = dist[end][d]
		}
	}

	if minScore == int(^uint(0)>>1) {
		return -1, prev // If no path found
	}

	return minScore, prev
}

// isValid checks if the next position is within bounds and not a wall
func isValid(grid [][]rune, x, y int) bool {
	if y < 0 || y >= len(grid) || x < 0 || x >= len(grid[0]) {
		return false
	}
	return grid[y][x] != '#'
}

// findBestPathTiles finds all tiles that are part of any best path
func findBestPathTiles(prev map[Point][]State, end Point, grid [][]rune) int {
	bestTiles := make(map[Point]bool)
	visited := make(map[Point]bool)
	queue := []Point{end}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current] {
			continue
		}
		visited[current] = true
		bestTiles[current] = true

		predecessors, exists := prev[current]
		if !exists {
			continue
		}

		for _, pred := range predecessors {
			predPoint := Point{pred.x, pred.y}
			if !visited[predPoint] {
				queue = append(queue, predPoint)
			}
		}
	}

	return len(bestTiles)
}
