/*
вообще неправильные ответы:
➜ go run ./o1 sample.txt
Часть 1: 9605 (время выполнения: 332.833µs)
Часть 2: 19631 (время выполнения: 348.875µs)

Для сравнения какие должны быть:
➜ go run . sample.txt
Part 1: 10092           in 113.667µs
Part 2: 9021            in 156.125µs
*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// Direction представляет направление движения
type Direction struct {
	dx int
	dy int
}

// MapGrid представляет карту склада
type MapGrid struct {
	grid     [][]rune
	robotX   int
	robotY   int
	width    int
	height   int
	boxes    map[string]bool
	walls    map[string]bool
	isScaled bool
}

// NewMapGrid инициализирует новую карту склада
func NewMapGrid(grid [][]rune, isScaled bool) *MapGrid {
	m := &MapGrid{
		grid:     grid,
		boxes:    make(map[string]bool),
		walls:    make(map[string]bool),
		isScaled: isScaled,
	}
	m.height = len(grid)
	if m.height > 0 {
		m.width = len(grid[0])
	}
	for y, row := range grid {
		for x, cell := range row {
			key := fmt.Sprintf("%d,%d", y, x)
			if cell == 'O' || cell == '[' {
				m.boxes[key] = true
			} else if cell == '#' {
				m.walls[key] = true
			} else if cell == '@' {
				m.robotY = y
				m.robotX = x
			}
		}
	}
	return m
}

// MoveRobot перемещает робота в указанном направлении
func (m *MapGrid) MoveRobot(dir Direction) {
	newY := m.robotY + dir.dy
	newX := m.robotX + dir.dx
	targetKey := fmt.Sprintf("%d,%d", newY, newX)

	// Проверка на стену
	if m.walls[targetKey] {
		return
	}

	// Проверка на коробку
	if m.boxes[targetKey] {
		if m.isScaled {
			// В Части 2 коробка занимает два символа: [ и ]
			// Проверяем, является ли текущая коробка левой частью
			if m.grid[newY][newX] != '[' {
				// Правая часть коробки не должна двигаться отдельно
				return
			}
			// Новые позиции для обеих частей коробки
			boxNewY1 := newY + dir.dy
			boxNewX1 := newX + dir.dx
			boxNewY2 := newY + dir.dy
			boxNewX2 := newX + 1 + dir.dx

			boxNewKey1 := fmt.Sprintf("%d,%d", boxNewY1, boxNewX1)
			boxNewKey2 := fmt.Sprintf("%d,%d", boxNewY2, boxNewX2)

			// Проверяем, могут ли обе части коробки быть перемещены
			if m.walls[boxNewKey1] || m.walls[boxNewKey2] || m.boxes[boxNewKey1] || m.boxes[boxNewKey2] {
				return
			}

			// Перемещаем обе части коробки
			delete(m.boxes, targetKey)
			delete(m.boxes, fmt.Sprintf("%d,%d", newY, newX+1))
			m.boxes[boxNewKey1] = true
			m.boxes[boxNewKey2] = true

			// Обновляем сетку
			m.grid[newY][newX] = '.'   // Освобождаем старую левую часть
			m.grid[newY][newX+1] = '.' // Освобождаем старую правую часть
			m.grid[boxNewY1][boxNewX1] = '['
			m.grid[boxNewY2][boxNewX2] = ']'

		} else {
			// В Части 1 коробка занимает один символ
			// Новая позиция для коробки
			boxNewY := newY + dir.dy
			boxNewX := newX + dir.dx
			boxNewKey := fmt.Sprintf("%d,%d", boxNewY, boxNewX)

			// Проверяем, может ли коробка быть перемещена
			if m.walls[boxNewKey] || m.boxes[boxNewKey] {
				return
			}

			// Перемещаем коробку
			delete(m.boxes, targetKey)
			m.boxes[boxNewKey] = true

			// Обновляем сетку
			m.grid[newY][newX] = 'O'
			m.grid[m.robotY][m.robotX] = '.'
			m.grid[newY][newX] = '@'
		}
	}

	// Перемещаем робота
	m.grid[m.robotY][m.robotX] = '.' // Освобождаем старую позицию
	m.robotY = newY
	m.robotX = newX
	m.grid[m.robotY][m.robotX] = '@' // Устанавливаем новую позицию робота
}

// SumGPSCoordinates вычисляет сумму GPS координат всех коробок
func (m *MapGrid) SumGPSCoordinates() int {
	sum := 0
	visited := make(map[string]bool)

	for key := range m.boxes {
		if visited[key] {
			continue
		}
		y, x := 0, 0
		fmt.Sscanf(key, "%d,%d", &y, &x)
		cell := m.grid[y][x]
		if cell == 'O' {
			sum += 100*y + x
			visited[key] = true
		} else if cell == '[' {
			// Для коробки из двух символов считаем только левую часть
			sum += 100*y + x
			visited[key] = true
			// Также отмечаем правую часть как посещенную
			rightKey := fmt.Sprintf("%d,%d", y, x+1)
			visited[rightKey] = true
		}
	}
	return sum
}

// Clone создает глубокую копию карты склада
func (m *MapGrid) Clone() *MapGrid {
	newGrid := make([][]rune, m.height)
	for y := range m.grid {
		newGrid[y] = make([]rune, m.width)
		copy(newGrid[y], m.grid[y])
	}
	newM := &MapGrid{
		grid:     newGrid,
		robotX:   m.robotX,
		robotY:   m.robotY,
		width:    m.width,
		height:   m.height,
		boxes:    make(map[string]bool),
		walls:    make(map[string]bool),
		isScaled: m.isScaled,
	}
	for k := range m.boxes {
		newM.boxes[k] = true
	}
	for k := range m.walls {
		newM.walls[k] = true
	}
	return newM
}

// ScaleMap удваивает ширину карты согласно Части 2
func (m *MapGrid) ScaleMap() *MapGrid {
	newWidth := m.width * 2
	newGrid := make([][]rune, m.height)
	for y := range m.grid {
		newGrid[y] = make([]rune, newWidth)
		for x := 0; x < m.width; x++ {
			cell := m.grid[y][x]
			switch cell {
			case '#':
				newGrid[y][2*x] = '#'
				newGrid[y][2*x+1] = '#'
			case 'O':
				newGrid[y][2*x] = 'O'
				newGrid[y][2*x+1] = 'O'
			case '.':
				newGrid[y][2*x] = '.'
				newGrid[y][2*x+1] = '.'
			case '@':
				newGrid[y][2*x] = '@'
				newGrid[y][2*x+1] = '.'
			case '[':
				newGrid[y][2*x] = '['
				newGrid[y][2*x+1] = ']'
			case ']':
				// Обычно не должно быть отдельных ']' без '['
				newGrid[y][2*x] = '['
				newGrid[y][2*x+1] = ']'
			default:
				newGrid[y][2*x] = cell
				newGrid[y][2*x+1] = cell
			}
		}
	}

	// Обновляем позиции коробок и стен
	scaledMap := NewMapGrid(newGrid, true)
	return scaledMap
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: go run main.go <файл_входных_данных>")
		return
	}

	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Ошибка при открытии файла: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var mapLines []string
	var moveLines []string
	mapEnded := false

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			mapEnded = true
			continue
		}
		if !mapEnded {
			mapLines = append(mapLines, line)
		} else {
			moveLines = append(moveLines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Ошибка при чтении файла: %v\n", err)
		return
	}

	// Объединяем линии движения в одну строку, удаляя новые строки
	moveSeq := strings.Join(moveLines, "")
	moveSeq = strings.ReplaceAll(moveSeq, "\n", "")
	moveSeq = strings.ReplaceAll(moveSeq, "\r", "")

	// Создаём карту
	grid := make([][]rune, len(mapLines))
	for i, line := range mapLines {
		grid[i] = []rune(line)
	}
	originalMap := NewMapGrid(grid, false)

	// Часть 1
	start1 := time.Now()
	map1 := originalMap.Clone()
	directions := parseMoves(moveSeq)
	for _, dir := range directions {
		map1.MoveRobot(dir)
	}
	sum1 := map1.SumGPSCoordinates()
	elapsed1 := time.Since(start1)
	fmt.Printf("Часть 1: %d (время выполнения: %v)\n", sum1, elapsed1)

	// Часть 2
	start2 := time.Now()
	scaledMap := originalMap.ScaleMap()
	directions2 := parseMoves(moveSeq)
	for _, dir := range directions2 {
		scaledMap.MoveRobot(dir)
	}
	sum2 := scaledMap.SumGPSCoordinates()
	elapsed2 := time.Since(start2)
	fmt.Printf("Часть 2: %d (время выполнения: %v)\n", sum2, elapsed2)
}

// parseMoves преобразует строку движений в срез направлений
func parseMoves(moves string) []Direction {
	var directions []Direction
	for _, ch := range moves {
		switch ch {
		case '<':
			directions = append(directions, Direction{dx: -1, dy: 0})
		case '>':
			directions = append(directions, Direction{dx: 1, dy: 0})
		case '^':
			directions = append(directions, Direction{dx: 0, dy: -1})
		case 'v':
			directions = append(directions, Direction{dx: 0, dy: 1})
		}
	}
	return directions
}
