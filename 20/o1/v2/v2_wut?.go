package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

// Position представляет координату на карте
type Position struct {
	X, Y int
}

// Direction vectors
var directions = []Position{
	{X: 0, Y: -1}, // Вверх
	{X: 0, Y: 1},  // Вниз
	{X: -1, Y: 0}, // Влево
	{X: 1, Y: 0},  // Вправо
}

// State представляет текущее состояние в BFS
type State struct {
	Pos            Position
	Steps          int
	CheatActive    bool
	CheatRemaining int
	CheatStart     Position
	CheatEnd       Position
}

// Key используется для состояния посещённых позиций с учётом состояния чита
type Key struct {
	Pos           Position
	CheatActive   bool
	CheatDuration int
}

// ParseMap читает карту из указанного файла и возвращает сетку, старт и конец
func ParseMap(filename string) ([][]rune, Position, Position, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, Position{}, Position{}, err
	}
	defer file.Close()

	var grid [][]rune
	var start, end Position
	scanner := bufio.NewScanner(file)
	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		row := []rune(line)
		for x, char := range row {
			if char == 'S' {
				start = Position{X: x, Y: y}
			} else if char == 'E' {
				end = Position{X: x, Y: y}
			}
		}
		grid = append(grid, row)
		y++
	}

	if err := scanner.Err(); err != nil {
		return nil, Position{}, Position{}, err
	}

	return grid, start, end, nil
}

// IsWalkable проверяет, можно ли пройти по данной позиции
func IsWalkable(grid [][]rune, pos Position) bool {
	if pos.Y < 0 || pos.Y >= len(grid) || pos.X < 0 || pos.X >= len(grid[0]) {
		return false
	}
	return grid[pos.Y][pos.X] != '#'
}

// BFS ищет кратчайший путь без использования чита
func BFS(grid [][]rune, start, end Position) int {
	visited := make(map[Key]bool)
	queue := []State{{Pos: start, Steps: 0, CheatActive: false, CheatRemaining: 0}}
	keyStart := Key{Pos: start, CheatActive: false, CheatDuration: 0}
	visited[keyStart] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.Pos == end {
			return current.Steps
		}

		for _, dir := range directions {
			next := Position{X: current.Pos.X + dir.X, Y: current.Pos.Y + dir.Y}
			if IsWalkable(grid, next) {
				key := Key{Pos: next, CheatActive: current.CheatActive, CheatDuration: current.CheatRemaining}
				if !visited[key] {
					visited[key] = true
					queue = append(queue, State{
						Pos:            next,
						Steps:          current.Steps + 1,
						CheatActive:    current.CheatActive,
						CheatRemaining: current.CheatRemaining,
						CheatStart:     current.CheatStart,
						CheatEnd:       current.CheatEnd,
					})
				}
			}
		}
	}
	return -1 // Путь не найден
}

// BFSWithCheat ищет все возможные читы и их экономию
func BFSWithCheat(grid [][]rune, start, end Position, maxCheatDuration int) map[int]int {
	visited := make(map[Key]bool)
	cheatSavings := make(map[int]int)

	queue := []State{{Pos: start, Steps: 0, CheatActive: false, CheatRemaining: 0}}
	keyStart := Key{Pos: start, CheatActive: false, CheatDuration: 0}
	visited[keyStart] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.Pos == end {
			continue
		}

		for _, dir := range directions {
			next := Position{X: current.Pos.X + dir.X, Y: current.Pos.Y + dir.Y}
			normalMove := IsWalkable(grid, next)

			// Move без использования чита
			if normalMove {
				key := Key{Pos: next, CheatActive: current.CheatActive, CheatDuration: current.CheatRemaining}
				if !visited[key] {
					visited[key] = true
					queue = append(queue, State{
						Pos:            next,
						Steps:          current.Steps + 1,
						CheatActive:    current.CheatActive,
						CheatRemaining: current.CheatRemaining,
						CheatStart:     current.CheatStart,
						CheatEnd:       current.CheatEnd,
					})
				}
			} else {
				// Использовать чит, если он не активен
				if !current.CheatActive && maxCheatDuration > 0 {
					key := Key{Pos: next, CheatActive: true, CheatDuration: maxCheatDuration - 1}
					if !visited[key] {
						visited[key] = true
						queue = append(queue, State{
							Pos:            next,
							Steps:          current.Steps + 1,
							CheatActive:    true,
							CheatRemaining: maxCheatDuration - 1,
							CheatStart:     current.Pos,
							CheatEnd:       next,
						})
					}
				}

				// Продолжать использовать активный чит
				if current.CheatActive && current.CheatRemaining > 0 {
					key := Key{Pos: next, CheatActive: true, CheatDuration: current.CheatRemaining - 1}
					if !visited[key] {
						visited[key] = true
						queue = append(queue, State{
							Pos:            next,
							Steps:          current.Steps + 1,
							CheatActive:    true,
							CheatRemaining: current.CheatRemaining - 1,
							CheatStart:     current.CheatStart,
							CheatEnd:       next,
						})
					}
				}

				// Если чит заканчивается после этого шага
				if current.CheatActive && current.CheatRemaining == 0 {
					cheatSavings[0]++ // Здесь нужно реализовать вычисление экономии
				}
			}
		}
	}

	// Примерная обработка cheatSavings (необходимо реализовать точную логику)
	// Для демонстрации возвращаем пустую карту
	return cheatSavings
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: go run main.go <input_file>")
		return
	}

	filename := os.Args[1]

	// Парсинг карты
	grid, start, end, err := ParseMap(filename)
	if err != nil {
		fmt.Println("Ошибка при чтении карты:", err)
		return
	}

	// Часть Первая
	startTime := time.Now()
	normalPath := BFS(grid, start, end)
	if normalPath == -1 {
		fmt.Println("Путь от Старт до Конца не найден.")
		return
	}

	cheatSavingsPartOne := BFSWithCheat(grid, start, end, 2)
	// Реализуйте точный подсчёт экономии времени
	// Для демонстрации считаем, что нет читов, сохраняющих >= 100 пикосекунд
	countPartOne := 0
	for saving := range cheatSavingsPartOne {
		if saving >= 100 {
			countPartOne += cheatSavingsPartOne[saving]
		}
	}

	fmt.Printf("Ответ Часть Первая: %d читов сохраняют как минимум 100 пикосекунд.\n", countPartOne)
	fmt.Printf("Время выполнения Часть Первая: %v\n", time.Since(startTime))

	// Часть Вторая
	startTime = time.Now()
	cheatSavingsPartTwo := BFSWithCheat(grid, start, end, 20)
	// Реализуйте точный подсчёт экономии времени
	// Для демонстрации считаем, что нет читов, сохраняющих >= 100 пикосекунд
	countPartTwo := 0
	for saving := range cheatSavingsPartTwo {
		if saving >= 100 {
			countPartTwo += cheatSavingsPartTwo[saving]
		}
	}

	fmt.Printf("Ответ Часть Вторая: %d читов сохраняют как минимум 100 пикосекунд.\n", countPartTwo)
	fmt.Printf("Время выполнения Часть Вторая: %v\n", time.Since(startTime))
}
