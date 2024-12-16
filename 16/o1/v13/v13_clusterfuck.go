/*
➜ go run ./o1/v12 sample.txt
Part 1: Lowest score = 7036
Time taken: 190.917µs
Part 2: Number of tiles on best paths = 45
Time taken: 18.375µs

Ответ верный.

➜ go run ./o1/v12 input.txt
Part 1: Lowest score = 89460
Time taken: 31.448ms
Part 2: Number of tiles on best paths = 586
Time taken: 220µs

Ответ второй части неверный. Ожидаемое значение: 504
*/

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

	// Часть 1 и Часть 2
	startTime := time.Now()

	// Первый проход: от S до всех плиток
	minScore_S := runDijkstra(grid, start)

	// Второй проход: от E до всех плиток
	minScore_E := runDijkstra(grid, end)

	// Определение минимального счёта пути от S до E
	minScore_S_E := minScore_S[end]

	part1Time := time.Since(startTime)
	if minScore_S_E == -1 {
		fmt.Println("Part 1: No path found")
	} else {
		fmt.Printf("Part 1: Lowest score = %d\nTime taken: %v\n", minScore_S_E, part1Time)
	}

	// Часть 2: Сбор уникальных плиток на лучших путях
	startTime = time.Now()
	bestPathTiles := findBestPathTiles(grid, minScore_S, minScore_E, minScore_S_E)
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

// runDijkstra выполняет алгоритм Дейкстры от заданной точки и возвращает карту минимальных счётов до всех плиток
func runDijkstra(grid [][]rune, start Point) map[Point]int {
	dist := make(map[Point]int)
	pq := &PriorityQueue{}
	heap.Init(pq)

	// Инициализация с четырьмя возможными направлениями из стартовой точки
	initialDirections := []int{East, South, West, North}
	for _, dir := range initialDirections {
		initialState := State{start.x, start.y, dir, 0}
		heap.Push(pq, initialState)
		point := Point{start.x, start.y}
		if existing, exists := dist[point]; !exists || 0 < existing {
			dist[point] = 0
		}
	}

	// Направления: Восток, Юг, Запад, Север
	dirs := []Point{
		{1, 0},  // East
		{0, 1},  // South
		{-1, 0}, // West
		{0, -1}, // North
	}

	for pq.Len() > 0 {
		current := heap.Pop(pq).(State)
		currentPoint := Point{current.x, current.y}

		// Если текущий счёт больше уже найденного, пропускаем
		if existing, exists := dist[currentPoint]; exists && current.score > existing {
			continue
		}

		// Возможные действия: движение вперед, поворот налево, поворот направо

		// Движение вперед
		nextX := current.x + dirs[current.direction].x
		nextY := current.y + dirs[current.direction].y
		if isValid(grid, nextX, nextY) {
			nextPoint := Point{nextX, nextY}
			newScore := current.score + 1
			if existing, exists := dist[nextPoint]; !exists || newScore < existing {
				dist[nextPoint] = newScore
				heap.Push(pq, State{nextX, nextY, current.direction, newScore})
			}
		}

		// Поворот налево
		newDirLeft := (current.direction + 3) % 4
		newScoreLeft := current.score + 1000
		leftPoint := Point{current.x, current.y}
		if existing, exists := dist[leftPoint]; !exists || newScoreLeft < existing {
			dist[leftPoint] = newScoreLeft
			heap.Push(pq, State{current.x, current.y, newDirLeft, newScoreLeft})
		}

		// Поворот направо
		newDirRight := (current.direction + 1) % 4
		newScoreRight := current.score + 1000
		rightPoint := Point{current.x, current.y}
		if existing, exists := dist[rightPoint]; !exists || newScoreRight < existing {
			dist[rightPoint] = newScoreRight
			heap.Push(pq, State{current.x, current.y, newDirRight, newScoreRight})
		}
	}

	return dist
}

// isValid проверяет, находится ли позиция внутри границ карты и не является стеной
func isValid(grid [][]rune, x, y int) bool {
	if y < 0 || y >= len(grid) || x < 0 || x >= len(grid[0]) {
		return false
	}
	return grid[y][x] != '#'
}

// findBestPathTiles находит все уникальные плитки, участвующие в любых оптимальных путях
func findBestPathTiles(grid [][]rune, minScore_S, minScore_E map[Point]int, minScore_S_E int) int {
	bestTiles := make(map[Point]bool)

	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[0]); x++ {
			point := Point{x, y}
			if grid[y][x] == '#' {
				continue
			}
			score_S, exists_S := minScore_S[point]
			score_E, exists_E := minScore_E[point]
			if exists_S && exists_E && (score_S+score_E == minScore_S_E) {
				bestTiles[point] = true
			}
		}
	}

	return len(bestTiles)
}
