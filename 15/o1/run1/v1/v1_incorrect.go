package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// Direction represents movement directions
type Direction struct {
	dx int
	dy int
}

// MapGrid represents the warehouse map
type MapGrid struct {
	grid   [][]rune
	robotX int
	robotY int
	width  int
	height int
	boxes  map[string]bool
	walls  map[string]bool
}

// NewMapGrid initializes a new MapGrid
func NewMapGrid(grid [][]rune) *MapGrid {
	m := &MapGrid{
		grid:  grid,
		boxes: make(map[string]bool),
		walls: make(map[string]bool),
	}
	m.height = len(grid)
	if m.height > 0 {
		m.width = len(grid[0])
	}
	for y, row := range grid {
		for x, cell := range row {
			key := fmt.Sprintf("%d,%d", y, x)
			if cell == 'O' {
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

// MoveRobot moves the robot based on the direction
func (m *MapGrid) MoveRobot(dir Direction) {
	newY := m.robotY + dir.dy
	newX := m.robotX + dir.dx
	targetKey := fmt.Sprintf("%d,%d", newY, newX)

	// Check if target is a wall
	if m.walls[targetKey] {
		return
	}

	// Check if target has a box
	if m.boxes[targetKey] {
		// Calculate position to push the box
		boxNewY := newY + dir.dy
		boxNewX := newX + dir.dx
		boxNewKey := fmt.Sprintf("%d,%d", boxNewY, boxNewX)

		// Check if box can be pushed
		if m.walls[boxNewKey] || m.boxes[boxNewKey] {
			return
		}

		// Push the box
		delete(m.boxes, targetKey)
		m.boxes[boxNewKey] = true
	}

	// Move the robot
	m.robotY = newY
	m.robotX = newX
}

// SumGPSCoordinates calculates the sum of GPS coordinates of all boxes
func (m *MapGrid) SumGPSCoordinates() int {
	sum := 0
	for key := range m.boxes {
		var y, x int
		fmt.Sscanf(key, "%d,%d", &y, &x)
		sum += 100*y + x
	}
	return sum
}

// Clone creates a deep copy of the MapGrid
func (m *MapGrid) Clone() *MapGrid {
	newGrid := make([][]rune, m.height)
	for y := range m.grid {
		newGrid[y] = make([]rune, m.width)
		copy(newGrid[y], m.grid[y])
	}
	newM := &MapGrid{
		grid:   newGrid,
		robotX: m.robotX,
		robotY: m.robotY,
		width:  m.width,
		height: m.height,
		boxes:  make(map[string]bool),
		walls:  make(map[string]bool),
	}
	for k := range m.boxes {
		newM.boxes[k] = true
	}
	for k := range m.walls {
		newM.walls[k] = true
	}
	return newM
}

// ScaleMap doubles the width of the map as per Part Two
func (m *MapGrid) ScaleMap() *MapGrid {
	newWidth := m.width * 2
	newGrid := make([][]rune, m.height)
	for y := range m.grid {
		newGrid[y] = make([]rune, newWidth)
		for x := range m.grid[y] {
			cell := m.grid[y][x]
			if cell == '#' {
				newGrid[y][2*x] = '#'
				newGrid[y][2*x+1] = '#'
			} else if cell == 'O' {
				newGrid[y][2*x] = '['
				newGrid[y][2*x+1] = ']'
			} else if cell == '.' {
				newGrid[y][2*x] = '.'
				newGrid[y][2*x+1] = '.'
			} else if cell == '@' {
				newGrid[y][2*x] = '@'
				newGrid[y][2*x+1] = '.'
			} else {
				newGrid[y][2*x] = cell
				newGrid[y][2*x+1] = cell
			}
		}
	}
	// Update walls and boxes
	newM := &MapGrid{
		grid:  newGrid,
		boxes: make(map[string]bool),
		walls: make(map[string]bool),
	}
	newM.height = m.height
	newM.width = newWidth
	for y, row := range newGrid {
		for x, cell := range row {
			key := fmt.Sprintf("%d,%d", y, x)
			if cell == 'O' || cell == '[' || cell == ']' {
				newM.boxes[key] = true
			} else if cell == '#' {
				newM.walls[key] = true
			} else if cell == '@' {
				newM.robotY = y
				newM.robotX = x
			}
		}
	}
	return newM
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
	originalMap := NewMapGrid(grid)

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
