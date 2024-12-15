package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Структуры для описания карты и состояния
type Tile int

const (
	Wall Tile = iota
	Empty
	Box
	Robot
)

// Для второго варианта: расширенный tile-тип
// В расширенной карте коробки - пара символов "[]", а робот - "@." и т.д.
// Но для внутреннего представления логика та же, изменится только способ создания карты.
type Coord struct {
	R, C int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go input.txt")
		return
	}
	inputFile := os.Args[1]
	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	lines := strings.Split(strings.ReplaceAll(string(data), "\r", ""), "\n")

	// Считаем входные данные: карта и последовательность движений
	// Карта окружена стенами. Неизвестно точное кол-во строк. Ищем раздел: первая пустая строка - после карты идут ходы.
	var mapLines []string
	var movesLines []string
	readingMap := true
	for _, l := range lines {
		if strings.TrimSpace(l) == "" {
			readingMap = false
			continue
		}
		if readingMap {
			mapLines = append(mapLines, l)
		} else {
			movesLines = append(movesLines, l)
		}
	}
	movesStr := strings.ReplaceAll(strings.Join(movesLines, ""), "\n", "")
	movesStr = strings.ReplaceAll(movesStr, " ", "")

	// Парсим карту
	originalMap, robotPos := parseMap(mapLines)
	// Решаем часть 1
	start1 := time.Now()
	finalMap1 := simulate(originalMap, robotPos, movesStr, false)
	res1 := calcSum(finalMap1, false)
	dur1 := time.Since(start1)
	fmt.Println(res1)
	fmt.Println(dur1)

	// Решаем часть 2 (масштабирование карты)
	expandedMap, expandedRobotPos := expandMap(originalMap)
	start2 := time.Now()
	finalMap2 := simulate(expandedMap, expandedRobotPos, movesStr, true)
	res2 := calcSum(finalMap2, true)
	dur2 := time.Since(start2)
	fmt.Println(res2)
	fmt.Println(dur2)
}

func parseMap(lines []string) ([][]Tile, Coord) {
	h := len(lines)
	w := len(lines[0])
	grid := make([][]Tile, h)
	var robot Coord
	for r := 0; r < h; r++ {
		grid[r] = make([]Tile, w)
		for c := 0; c < w; c++ {
			switch lines[r][c] {
			case '#':
				grid[r][c] = Wall
			case '.':
				grid[r][c] = Empty
			case 'O':
				grid[r][c] = Box
			case '@':
				grid[r][c] = Robot
				robot = Coord{r, c}
			}
		}
	}
	return grid, robot
}

// При расширении:
// # -> ##
// O -> []
// . -> ..
// @ -> @.
// Высота не меняется, ширина - в 2 раза больше.
// Робот занимает одну клетку, но в итоговой строке он будет '@' и справа '.' для заполнения.
func expandMap(original [][]Tile) ([][]Tile, Coord) {
	// Определяем размеры
	h := len(original)
	w := len(original[0])
	newW := w * 2
	newMap := make([][]Tile, h)
	var robot Coord
	for r := 0; r < h; r++ {
		newMap[r] = make([]Tile, newW)
		col := 0
		for c := 0; c < w; c++ {
			switch original[r][c] {
			case Wall:
				// ##
				newMap[r][col] = Wall
				newMap[r][col+1] = Wall
				col += 2
			case Empty:
				// ..
				newMap[r][col] = Empty
				newMap[r][col+1] = Empty
				col += 2
			case Box:
				// []
				// Для удобства считаем левую ячейку как Box, правую тоже Box, так будет проще.
				newMap[r][col] = Box
				newMap[r][col+1] = Box
				col += 2
			case Robot:
				// @.
				newMap[r][col] = Robot
				newMap[r][col+1] = Empty
				robot = Coord{r, col} // робот находится в левой клетке из двух
				col += 2
			}
		}
	}
	return newMap, robot
}

func simulate(grid [][]Tile, robot Coord, moves string, expanded bool) [][]Tile {
	// Для каждого хода двигаем робота с учетом правил
	for _, m := range moves {
		dr, dc := 0, 0
		switch m {
		case '^':
			dr = -1
		case 'v':
			dr = 1
		case '<':
			dc = -1
		case '>':
			dc = 1
		}
		tryMove(grid, &robot, dr, dc, expanded)
	}
	return grid
}

// Попытка сдвига робота в направлении (dr,dc).
// Если целевая клетка пуста - просто двигаем.
// Если в целевой клетке box - пытаемся "цепочкой" сдвинуть все box в том направлении.
// Если упирается в стену - никто не двигается.
func tryMove(grid [][]Tile, robot *Coord, dr, dc int, expanded bool) {
	R, C := robot.R, robot.C
	nr, nc := R+dr, C+dc
	if nr < 0 || nr >= len(grid) || nc < 0 || nc >= len(grid[0]) {
		return
	}
	if grid[nr][nc] == Wall {
		return
	}
	if grid[nr][nc] == Empty {
		// Двигаем робота
		grid[R][C] = Empty
		grid[nr][nc] = Robot
		robot.R, robot.C = nr, nc
		return
	}
	if grid[nr][nc] == Box {
		// Пытаемся сдвинуть цепочку ящиков
		if pushBoxes(grid, nr, nc, dr, dc) {
			// Если получилось сдвинуть ящики, двигаем робота
			grid[R][C] = Empty
			grid[nr][nc] = Robot
			robot.R, robot.C = nr, nc
		}
	}
}

// Рекурсивный/цепной сдвиг ящиков
func pushBoxes(grid [][]Tile, r, c, dr, dc int) bool {
	h, w := len(grid), len(grid[0])
	// Найдем цепочку ящиков в направлении dr,dc
	boxes := make([]Coord, 0)
	cr, cc := r, c
	for {
		if cr < 0 || cr >= h || cc < 0 || cc >= w {
			return false
		}
		if grid[cr][cc] != Box {
			break
		}
		boxes = append(boxes, Coord{cr, cc})
		cr += dr
		cc += dc
	}
	// cr,cc - следующая клетка за коробками
	if cr < 0 || cr >= h || cc < 0 || cc >= w {
		return false
	}
	if grid[cr][cc] == Wall || grid[cr][cc] == Robot {
		return false
	}
	if grid[cr][cc] == Box {
		// Еще коробки, продолжаем
		if !pushBoxes(grid, cr, cc, dr, dc) {
			return false
		}
	}
	// Теперь можно сдвинуть все коробки
	// Сдвигаем с конца цепочки, чтобы не перезаписать
	for i := len(boxes) - 1; i >= 0; i-- {
		br, bc := boxes[i].R, boxes[i].C
		grid[br+dr][bc+dc] = Box
		grid[br][bc] = Empty
	}
	return true
}

// Вычисление суммы GPS координат
// Для первой части: GPS = 100*r + c, где r,c - дистанция от верхнего/левого края карты.
// Для второй части: для коробок шириной 2 клетки, координаты считаются по ближнему краю коробки.
// Так как мы оставили внутреннее представление одинаковым, для второй части коробка занимает две клетки.
// Нам нужно идентифицировать "левую" половину коробки. Любая пара Box идет подряд по горизонтали.
// Для второй части считаем координаты по "левой" клетке пары.
// Для первой части - каждая Box - отдельна.
func calcSum(grid [][]Tile, expanded bool) int {
	sum := 0
	h, w := len(grid), len(grid[0])
	if !expanded {
		// Каждая коробка - единичный tile
		// GPS: distance from top = row index, distance from left = column index
		// sum(100*r + c)
		for r := 0; r < h; r++ {
			for c := 0; c < w; c++ {
				if grid[r][c] == Box {
					sum += 100*r + c
				}
			}
		}
	} else {
		// Коробки двойные по горизонтали
		// Ищем пары Box[]. Левая часть - координаты.
		// Считаем только левую часть каждой коробки-пары.
		// Предполагаем, что все коробки корректны (парные).
		for r := 0; r < h; r++ {
			c := 0
			for c < w {
				if grid[r][c] == Box {
					// Проверим c+1
					if c+1 < w && grid[r][c+1] == Box {
						// Это коробка
						sum += 100*r + c
						c += 2
					} else {
						// Если вдруг одиночный Box (теоретически не должно быть), считаем как есть
						sum += 100*r + c
						c++
					}
				} else {
					c++
				}
			}
		}
	}
	return sum
}
