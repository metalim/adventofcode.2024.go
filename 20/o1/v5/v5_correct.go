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

// Directions для перемещения: вверх, вниз, влево, вправо
var directions = []Position{
	{X: 0, Y: -1}, // Вверх
	{X: 0, Y: 1},  // Вниз
	{X: -1, Y: 0}, // Влево
	{X: 1, Y: 0},  // Вправо
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

// BFS выполняет поиск в ширину на карте, учитывая, можно ли проходить через стены
// Если allowWalls=true, разрешено проходить через стены
func BFS(grid [][]rune, start Position, allowWalls bool) [][]int {
	height := len(grid)
	width := len(grid[0])
	dist := make([][]int, height)
	for i := range dist {
		dist[i] = make([]int, width)
		for j := range dist[i] {
			dist[i][j] = -1 // -1 означает, что позиция ещё не посещена
		}
	}

	queue := []Position{start}
	dist[start.Y][start.X] = 0

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		currentDist := dist[current.Y][current.X]

		for _, dir := range directions {
			next := Position{X: current.X + dir.X, Y: current.Y + dir.Y}
			if next.Y < 0 || next.Y >= height || next.X < 0 || next.X >= width {
				continue
			}
			if dist[next.Y][next.X] != -1 {
				continue
			}
			if grid[next.Y][next.X] == '#' && !allowWalls {
				continue
			}
			dist[next.Y][next.X] = currentDist + 1
			queue = append(queue, next)
		}
	}

	return dist
}

// BFSCheat выполняет BFS с возможностью проходить через стены, ограничивая количество шагов
// Возвращает карту с количеством шагов до каждой позиции при использовании чита
func BFSCheat(grid [][]rune, start Position, maxCheatDuration int) map[Position]int {
	height := len(grid)
	width := len(grid[0])
	visited := make([][]bool, height)
	for i := range visited {
		visited[i] = make([]bool, width)
	}
	positions := make(map[Position]int)

	type Node struct {
		Pos   Position
		Steps int
	}

	queue := []Node{{Pos: start, Steps: 0}}
	visited[start.Y][start.X] = true
	positions[start] = 0

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.Steps >= maxCheatDuration {
			continue
		}

		for _, dir := range directions {
			next := Position{X: current.Pos.X + dir.X, Y: current.Pos.Y + dir.Y}
			if next.Y < 0 || next.Y >= height || next.X < 0 || next.X >= width {
				continue
			}
			if visited[next.Y][next.X] {
				continue
			}
			// При использовании чита, можно проходить через стены
			visited[next.Y][next.X] = true
			positions[next] = current.Steps + 1
			queue = append(queue, Node{Pos: next, Steps: current.Steps + 1})
		}
	}

	return positions
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

	// Вычисление расстояний без использования чита
	startTime := time.Now()
	distFromStart := BFS(grid, start, false)
	distFromEnd := BFS(grid, end, false)
	normalPath := distFromStart[end.Y][end.X]
	if normalPath == -1 {
		fmt.Println("Путь от Старт до Конца не найден.")
		return
	}
	fmt.Printf("Нормальное время прохождения трассы: %d пикосекунд.\n", normalPath)
	fmt.Printf("Время выполнения вычисления расстояний: %v\n", time.Since(startTime))

	// Функция для подсчета читов
	countCheats := func(cheatDuration int) int {
		startTime := time.Now()
		count := 0

		// Итерация по всем возможным позициям для начала чита
		for y := 0; y < len(grid); y++ {
			for x := 0; x < len(grid[0]); x++ {
				startCheat := Position{X: x, Y: y}
				P1 := distFromStart[y][x]
				if P1 == -1 {
					continue // Позиция недостижима без использования чита
				}

				// Получение всех endCheat позиций, достижимых из startCheat с использованием чита
				endCheats := BFSCheat(grid, startCheat, cheatDuration)
				for endCheat, cheatSteps := range endCheats {
					P3 := distFromEnd[endCheat.Y][endCheat.X]
					if P3 == -1 {
						continue // Позиция endCheat недостижима до конца без использования чита
					}
					P_cheat := P1 + cheatSteps + P3
					saving := normalPath - P_cheat
					if saving >= 100 {
						count++
					}
				}
			}
		}

		fmt.Printf("Время выполнения подсчета для cheatDuration=%d: %v\n", cheatDuration, time.Since(startTime))
		return count
	}

	// Часть Первая: cheatDuration=2
	startTime = time.Now()
	countPartOne := countCheats(2)
	fmt.Printf("Ответ Часть Первая: %d читов сохраняют как минимум 100 пикосекунд.\n", countPartOne)
	fmt.Printf("Общее время выполнения Часть Первая: %v\n", time.Since(startTime))

	// Часть Вторая: cheatDuration=20
	startTime = time.Now()
	countPartTwo := countCheats(20)
	fmt.Printf("Ответ Часть Вторая: %d читов сохраняют как минимум 100 пикосекунд.\n", countPartTwo)
	fmt.Printf("Общее время выполнения Часть Вторая: %v\n", time.Since(startTime))
}
