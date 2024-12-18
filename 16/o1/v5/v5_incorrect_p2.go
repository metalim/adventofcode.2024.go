package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"time"
)

// Направления
const (
	East = iota
	South
	West
	North
)

// State представляет текущее состояние: позицию, направление и накопленный счёт
type State struct {
	x, y      int
	direction int
	score     int
}

// StateKey уникально идентифицирует состояние по позиции и направлению
type StateKey struct {
	x, y      int
	direction int
}

// PriorityQueue реализует интерфейс heap.Interface и хранит состояния
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

// Point представляет координаты на карте
type Point struct {
	x, y int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <input_file>")
		return
	}

	filename := os.Args[1]

	// Чтение карты
	grid, start, end, err := readMap(filename)
	if err != nil {
		fmt.Println("Error reading map:", err)
		return
	}

	// Часть 1
	startTime := time.Now()
	minScore, prevMap := dijkstra(grid, start, end)
	part1Time := time.Since(startTime)
	if minScore == -1 {
		fmt.Println("Part 1: No path found")
	} else {
		fmt.Printf("Part 1: Lowest score = %d\nTime taken: %v\n", minScore, part1Time)
	}

	// Часть 2
	startTime = time.Now()
	bestPathTiles := findBestPathTiles(prevMap, end, grid)
	part2Time := time.Since(startTime)
	fmt.Printf("Part 2: Number of tiles on best paths = %d\nTime taken: %v\n", bestPathTiles, part2Time)
}

// readMap читает входной файл и возвращает сетку, точки старта и конца
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

// dijkstra выполняет алгоритм Дейкстры для поиска пути с минимальным счётом
func dijkstra(grid [][]rune, start, end Point) (int, map[StateKey][]StateKey) {
	// Инициализация карты расстояний
	dist := make(map[StateKey]int)
	prev := make(map[StateKey][]StateKey)

	// Инициализация очереди приоритетов
	pq := &PriorityQueue{}
	heap.Init(pq)
	initialState := State{start.x, start.y, East, 0}
	heap.Push(pq, initialState)
	initialKey := StateKey{start.x, start.y, East}
	dist[initialKey] = 0

	// Направления: Восток, Юг, Запад, Север
	dirs := []Point{
		{1, 0},  // East
		{0, 1},  // South
		{-1, 0}, // West
		{0, -1}, // North
	}

	for pq.Len() > 0 {
		current := heap.Pop(pq).(State)
		currentKey := StateKey{current.x, current.y, current.direction}

		// Если достигли конца, продолжаем искать все возможные пути
		if current.x == end.x && current.y == end.y {
			continue
		}

		// Возможные действия: движение вперед, поворот налево, поворот направо

		// Движение вперед
		nextX := current.x + dirs[current.direction].x
		nextY := current.y + dirs[current.direction].y
		if isValid(grid, nextX, nextY) {
			nextKey := StateKey{nextX, nextY, current.direction}
			newScore := current.score + 1
			if existingScore, exists := dist[nextKey]; !exists || newScore < existingScore {
				dist[nextKey] = newScore
				heap.Push(pq, State{nextX, nextY, current.direction, newScore})
				prev[nextKey] = append(prev[nextKey], currentKey)
			} else if newScore == existingScore {
				prev[nextKey] = append(prev[nextKey], currentKey)
			}
		}

		// Поворот налево
		newDirLeft := (current.direction + 3) % 4
		leftKey := StateKey{current.x, current.y, newDirLeft}
		newScoreLeft := current.score + 1000
		if existingScore, exists := dist[leftKey]; !exists || newScoreLeft < existingScore {
			dist[leftKey] = newScoreLeft
			heap.Push(pq, State{current.x, current.y, newDirLeft, newScoreLeft})
			prev[leftKey] = append(prev[leftKey], currentKey)
		} else if newScoreLeft == existingScore {
			prev[leftKey] = append(prev[leftKey], currentKey)
		}

		// Поворот направо
		newDirRight := (current.direction + 1) % 4
		rightKey := StateKey{current.x, current.y, newDirRight}
		newScoreRight := current.score + 1000
		if existingScore, exists := dist[rightKey]; !exists || newScoreRight < existingScore {
			dist[rightKey] = newScoreRight
			heap.Push(pq, State{current.x, current.y, newDirRight, newScoreRight})
			prev[rightKey] = append(prev[rightKey], currentKey)
		} else if newScoreRight == existingScore {
			prev[rightKey] = append(prev[rightKey], currentKey)
		}
	}

	// Определение минимального счёта для достижения конца
	minScore := -1
	for d := 0; d < 4; d++ {
		key := StateKey{end.x, end.y, d}
		if score, exists := dist[key]; exists {
			if minScore == -1 || score < minScore {
				minScore = score
			}
		}
	}

	return minScore, prev
}

// isValid проверяет, находится ли позиция внутри границ карты и не является стеной
func isValid(grid [][]rune, x, y int) bool {
	if y < 0 || y >= len(grid) || x < 0 || x >= len(grid[0]) {
		return false
	}
	return grid[y][x] != '#'
}

// findBestPathTiles находит все уникальные плитки, участвующие в любых оптимальных путях
func findBestPathTiles(prev map[StateKey][]StateKey, end Point, grid [][]rune) int {
	bestTiles := make(map[Point]bool)
	visited := make(map[StateKey]bool)

	// Сначала найдем все состояния, соответствующие конечной позиции
	endKeys := []StateKey{}
	for key := range prev {
		if key.x == end.x && key.y == end.y {
			endKeys = append(endKeys, key)
		}
	}

	// Если нет предшественников для конца, добавляем произвольное состояние
	if len(endKeys) == 0 {
		endKeys = append(endKeys, StateKey{end.x, end.y, East})
	}

	// Используем BFS для обратного обхода предшественников
	queue := endKeys
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current] {
			continue
		}
		visited[current] = true
		bestTiles[Point{current.x, current.y}] = true

		predecessors, exists := prev[current]
		if !exists {
			continue
		}

		for _, pred := range predecessors {
			if !visited[pred] {
				queue = append(queue, pred)
			}
		}
	}

	return len(bestTiles)
}
