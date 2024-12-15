/*
Напиши код на Go для решения задачи.
Входные данные в файле указываемом аргументом в командной строке.
Выведи ответ и время решения после решения каждой части.
Каждая часть должна решаться за несколько секунд максимум. Вторая часть задачи МОЖЕТ требовать особого подхода и не решаться перебором вариантов.
Если программа не сработает, я вставлю вывод и возможные комментарии. В ответ просто выдай исправленную версию.

https://adventofcode.com/2024/day/15

--- Day 15: Warehouse Woes ---
You appear back inside your own mini submarine! Each Historian drives their mini submarine in a different direction; maybe the Chief has his own submarine down here somewhere as well?

You look up to see a vast school of lanternfish swimming past you. On closer inspection, they seem quite anxious, so you drive your mini submarine over to see if you can help.

Because lanternfish populations grow rapidly, they need a lot of food, and that food needs to be stored somewhere. That's why these lanternfish have built elaborate warehouse complexes operated by robots!

These lanternfish seem so anxious because they have lost control of the robot that operates one of their most important warehouses! It is currently running amok, pushing around boxes in the warehouse with no regard for lanternfish logistics or lanternfish inventory management strategies.

Right now, none of the lanternfish are brave enough to swim up to an unpredictable robot so they could shut it off. However, if you could anticipate the robot's movements, maybe they could find a safe option.

The lanternfish already have a map of the warehouse and a list of movements the robot will attempt to make (your puzzle input). The problem is that the movements will sometimes fail as boxes are shifted around, making the actual movements of the robot difficult to predict.

For example:

##########
#..O..O.O#
#......O.#
#.OO..O.O#
#..O@..O.#
#O#..O...#
#O..O..O.#
#.OO.O.OO#
#....O...#
##########

<vv>^<v^>v>^vv^v>v<>v^v<v<^vv<<<^><<><>>v<vvv<>^v^>^<<<><<v<<<v^vv^v>^
vvv<<^>^v^^><<>>><>^<<><^vv^^<>vvv<>><^^v>^>vv<>v<<<<v<^v>^<^^>>>^<v<v
><>vv>v^v^<>><>>>><^^>vv>v<^^^>>v^v^<^^>v^^>v^<^v>v<>>v^v^<v>v^^<^^vv<
<<v<^>>^^^^>>>v^<>vvv^><v<<<>^^^vv^<vvv>^>v<^^^^v<>^>vvvv><>>v^<<^^^^^
^><^><>>><>^^<<^^v>>><^<v>^<vv>>v>>>^v><>^v><<<<v>>v<v<v>vvv>^<><<>^><
^>><>^v<><^vvv<^^<><v<<<<<><^v<<<><<<^^<v<^^^><^>>^<v^><<<^>>^v<v^v<v^
>^>>^v>vv>^<<^v<>><<><<v<<v><>v<^vv<<<>^^v^>^^>>><<^v>>v^v><^^>>^<>vv^
<><^^>^^^<><vvvvv^v<v<<>^v<v>v<<^><<><<><<<^^<<<^<<>><<><^^^>^^<>^>v<>
^^>vv<^v^v<vv>^<><v<^v>^^^>>>^^vvv^>vvv<>>>^<^>>>>>^<<^v>^vvv<>^<><<v>
v^^>>><<^^<>>^v^<v^vv<>v^<<>^<^v^v><^<<<><<^<v><v<>vv>>v><v^<vv<>v^<<^
As the robot (@) attempts to move, if there are any boxes (O) in the way, the robot will also attempt to push those boxes. However, if this action would cause the robot or a box to move into a wall (#), nothing moves instead, including the robot. The initial positions of these are shown on the map at the top of the document the lanternfish gave you.

The rest of the document describes the moves (^ for up, v for down, < for left, > for right) that the robot will attempt to make, in order. (The moves form a single giant sequence; they are broken into multiple lines just to make copy-pasting easier. Newlines within the move sequence should be ignored.)

Here is a smaller example to get started:

########
#..O.O.#
##@.O..#
#...O..#
#.#.O..#
#...O..#
#......#
########

<^^>>>vv<v>>v<<
Were the robot to attempt the given sequence of moves, it would push around the boxes as follows:

Initial state:
########
#..O.O.#
##@.O..#
#...O..#
#.#.O..#
#...O..#
#......#
########

Move <:
########
#..O.O.#
##@.O..#
#...O..#
#.#.O..#
#...O..#
#......#
########

Move ^:
########
#.@O.O.#
##..O..#
#...O..#
#.#.O..#
#...O..#
#......#
########

Move ^:
########
#.@O.O.#
##..O..#
#...O..#
#.#.O..#
#...O..#
#......#
########

Move >:
########
#..@OO.#
##..O..#
#...O..#
#.#.O..#
#...O..#
#......#
########

Move >:
########
#...@OO#
##..O..#
#...O..#
#.#.O..#
#...O..#
#......#
########

Move >:
########
#...@OO#
##..O..#
#...O..#
#.#.O..#
#...O..#
#......#
########

Move v:
########
#....OO#
##..@..#
#...O..#
#.#.O..#
#...O..#
#...O..#
########

Move v:
########
#....OO#
##..@..#
#...O..#
#.#.O..#
#...O..#
#...O..#
########

Move <:
########
#....OO#
##.@...#
#...O..#
#.#.O..#
#...O..#
#...O..#
########

Move v:
########
#....OO#
##.....#
#..@O..#
#.#.O..#
#...O..#
#...O..#
########

Move >:
########
#....OO#
##.....#
#...@O.#
#.#.O..#
#...O..#
#...O..#
########

Move >:
########
#....OO#
##.....#
#....@O#
#.#.O..#
#...O..#
#...O..#
########

Move v:
########
#....OO#
##.....#
#.....O#
#.#.O@.#
#...O..#
#...O..#
########

Move <:
########
#....OO#
##.....#
#.....O#
#.#O@..#
#...O..#
#...O..#
########

Move <:
########
#....OO#
##.....#
#.....O#
#.#O@..#
#...O..#
#...O..#
########
The larger example has many more moves; after the robot has finished those moves, the warehouse would look like this:

##########
#.O.O.OOO#
#........#
#OO......#
#OO@.....#
#O#.....O#
#O.....OO#
#O.....OO#
#OO....OO#
##########
The lanternfish use their own custom Goods Positioning System (GPS for short) to track the locations of the boxes. The GPS coordinate of a box is equal to 100 times its distance from the top edge of the map plus its distance from the left edge of the map. (This process does not stop at wall tiles; measure all the way to the edges of the map.)

So, the box shown below has a distance of 1 from the top edge of the map and 4 from the left edge of the map, resulting in a GPS coordinate of 100 * 1 + 4 = 104.

#######
#...O..
#......
The lanternfish would like to know the sum of all boxes' GPS coordinates after the robot finishes moving. In the larger example, the sum of all boxes' GPS coordinates is 10092. In the smaller example, the sum is 2028.

Predict the motion of the robot and boxes in the warehouse. After the robot is finished moving, what is the sum of all boxes' GPS coordinates?

--- Part Two ---
The lanternfish use your information to find a safe moment to swim in and turn off the malfunctioning robot! Just as they start preparing a festival in your honor, reports start coming in that a second warehouse's robot is also malfunctioning.

This warehouse's layout is surprisingly similar to the one you just helped. There is one key difference: everything except the robot is twice as wide! The robot's list of movements doesn't change.

To get the wider warehouse's map, start with your original map and, for each tile, make the following changes:

If the tile is #, the new map contains ## instead.
If the tile is O, the new map contains [] instead.
If the tile is ., the new map contains .. instead.
If the tile is @, the new map contains @. instead.
This will produce a new warehouse map which is twice as wide and with wide boxes that are represented by []. (The robot does not change size.)

The larger example from before would now look like this:

####################
##....[]....[]..[]##
##............[]..##
##..[][]....[]..[]##
##....[]@.....[]..##
##[]##....[]......##
##[]....[]....[]..##
##..[][]..[]..[][]##
##........[]......##
####################
Because boxes are now twice as wide but the robot is still the same size and speed, boxes can be aligned such that they directly push two other boxes at once. For example, consider this situation:

#######
#...#.#
#.....#
#..OO@#
#..O..#
#.....#
#######

<vv<<^^<<^^
After appropriately resizing this map, the robot would push around these boxes as follows:

Initial state:
##############
##......##..##
##..........##
##....[][]@.##
##....[]....##
##..........##
##############

Move <:
##############
##......##..##
##..........##
##...[][]@..##
##....[]....##
##..........##
##############

Move v:
##############
##......##..##
##..........##
##...[][]...##
##....[].@..##
##..........##
##############

Move v:
##############
##......##..##
##..........##
##...[][]...##
##....[]....##
##.......@..##
##############

Move <:
##############
##......##..##
##..........##
##...[][]...##
##....[]....##
##......@...##
##############

Move <:
##############
##......##..##
##..........##
##...[][]...##
##....[]....##
##.....@....##
##############

Move ^:
##############
##......##..##
##...[][]...##
##....[]....##
##.....@....##
##..........##
##############

Move ^:
##############
##......##..##
##...[][]...##
##....[]....##
##.....@....##
##..........##
##############

Move <:
##############
##......##..##
##...[][]...##
##....[]....##
##....@.....##
##..........##
##############

Move <:
##############
##......##..##
##...[][]...##
##....[]....##
##...@......##
##..........##
##############

Move ^:
##############
##......##..##
##...[][]...##
##...@[]....##
##..........##
##..........##
##############

Move ^:
##############
##...[].##..##
##...@.[]...##
##....[]....##
##..........##
##..........##
##############
This warehouse also uses GPS to locate the boxes. For these larger boxes, distances are measured from the edge of the map to the closest edge of the box in question. So, the box shown below has a distance of 1 from the top edge of the map and 5 from the left edge of the map, resulting in a GPS coordinate of 100 * 1 + 5 = 105.

##########
##...[]...
##........
In the scaled-up version of the larger example from above, after the robot has finished all of its moves, the warehouse would look like this:

####################
##[].......[].[][]##
##[]...........[].##
##[]........[][][]##
##[]......[]....[]##
##..##......[]....##
##..[]............##
##..@......[].[][]##
##......[][]..[]..##
####################
The sum of these boxes' GPS coordinates is 9021.

Predict the motion of the robot and boxes in this new, scaled-up warehouse. What is the sum of all boxes' final GPS coordinates?

*/

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
