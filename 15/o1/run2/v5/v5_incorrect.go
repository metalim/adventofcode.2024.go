package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Типы плиток на карте
type Tile int

const (
	Wall Tile = iota
	Empty
	Box
	Robot
)

// Координаты на карте
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

	// Разделяем карту и ходы по первой пустой строке
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
	finalMap1 := simulate(copyMap(originalMap), robotPos, movesStr, false)
	res1 := calcSum(finalMap1, false)
	dur1 := time.Since(start1)

	// Решаем часть 2 (масштабирование карты)
	expandedMap, expandedRobotPos := expandMap(originalMap)
	start2 := time.Now()
	finalMap2 := simulate(copyMap(expandedMap), expandedRobotPos, movesStr, true)
	res2 := calcSum(finalMap2, true)
	dur2 := time.Since(start2)

	// Выводим результаты
	fmt.Printf("Part 1: %d\tin %v\n", res1, dur1)
	fmt.Printf("Part 2: %d\tin %v\n", res2, dur2)

	// Отладочный вывод: список координат коробок для части 2
	fmt.Println("\n[DEBUG] Box positions after Part 2 simulation:")
	boxPositions := getBoxPositions(finalMap2, true)
	for _, pos := range boxPositions {
		fmt.Printf("Box at row %d, column %d (GPS: %d)\n", pos.R, pos.C, 100*pos.R+pos.C)
	}
}

// Функция для копирования карты
func copyMap(grid [][]Tile) [][]Tile {
	h := len(grid)
	w := len(grid[0])
	newGrid := make([][]Tile, h)
	for r := 0; r < h; r++ {
		newGrid[r] = make([]Tile, w)
		copy(newGrid[r], grid[r])
	}
	return newGrid
}

// Парсинг карты из строк
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
			default:
				grid[r][c] = Empty
			}
		}
	}
	return grid, robot
}

// Масштабирование карты для части 2
func expandMap(original [][]Tile) ([][]Tile, Coord) {
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
				newMap[r][col] = Box
				newMap[r][col+1] = Box
				col += 2
			case Robot:
				// @.
				newMap[r][col] = Robot
				newMap[r][col+1] = Empty
				robot = Coord{r, col}
				col += 2
			default:
				// Заполняем пустыми, если встретили неизвестный символ
				newMap[r][col] = Empty
				newMap[r][col+1] = Empty
				col += 2
			}
		}
	}
	return newMap, robot
}

// Симуляция перемещений робота
func simulate(grid [][]Tile, robot Coord, moves string, expanded bool) [][]Tile {
	// for idx, m := range moves {
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

		// Отладочный вывод после каждого шага (можно раскомментировать для подробной отладки)
		/*
			fmt.Printf("\n[DEBUG] After move %d (%c):\n", idx+1, m)
			printMap(grid, expanded)
		*/
	}
	return grid
}

// Попытка перемещения робота
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
		// Пытаемся сдвинуть коробки
		if pushBoxes(grid, nr, nc, dr, dc, expanded) {
			// Если получилось сдвинуть коробки, двигаем робота
			grid[R][C] = Empty
			grid[nr][nc] = Robot
			robot.R, robot.C = nr, nc
		}
	}
}

// Функция для сдвига коробок
func pushBoxes(grid [][]Tile, r, c, dr, dc int, expanded bool) bool {
	h, w := len(grid), len(grid[0])
	if expanded {
		// Проверяем, является ли (r,c) первой коробкой пары
		if c+1 >= w || grid[r][c+1] != Box {
			// Возможно, это вторая коробка пары, попробуем сдвинуть с начала пары
			if c-1 >= 0 && grid[r][c-1] == Box {
				// Сдвигаем пару начиная с (r, c-1)
				return pushBoxes(grid, r, c-1, dr, dc, expanded)
			}
			// Это одиночная коробка или не часть пары
			return false
		}

		// Теперь (r,c) и (r,c+1) - пара коробок
		tr1, tc1 := r+dr, c+dc
		tr2, tc2 := r+dr, c+dc+1

		// Проверяем границы
		if tr1 < 0 || tr1 >= h || tc1 < 0 || tc1 >= w || tr2 < 0 || tc2 >= w {
			return false
		}

		// Проверяем, свободны ли целевые позиции
		if grid[tr1][tc1] == Empty && grid[tr2][tc2] == Empty {
			// Сдвигаем обе коробки
			grid[tr1][tc1] = Box
			grid[tr2][tc2] = Box
			grid[r][c] = Empty
			grid[r][c+1] = Empty
			// Отладочный вывод
			fmt.Printf("[DEBUG] Shifted Box pair from (%d, %d)-(%d, %d) to (%d, %d)-(%d, %d)\n", r, c, r, c+1, tr1, tc1, tr2, tc2)
			return true
		}

		// Если целевые позиции заняты коробками, пытаемся сдвинуть их рекурсивно
		if grid[tr1][tc1] == Box && grid[tr2][tc2] == Box {
			if pushBoxes(grid, tr1, tc1, dr, dc, expanded) {
				// После успешного сдвига существующих коробок, сдвигаем текущую пару
				grid[tr1][tc1] = Box
				grid[tr2][tc2] = Box
				grid[r][c] = Empty
				grid[r][c+1] = Empty
				// Отладочный вывод
				fmt.Printf("[DEBUG] Recursively shifted Box pair from (%d, %d)-(%d, %d) to (%d, %d)-(%d, %d)\n", r, c, r, c+1, tr1, tc1, tr2, tc2)
				return true
			}
		}

		// В противном случае, не можем сдвинуть коробки
		return false
	} else {
		// В обычном режиме коробки одиночные
		// Определяем цепочку коробок
		boxes := []Coord{{r, c}}
		cr, cc := r+dr, c+dc
		for cr >= 0 && cr < h && cc >= 0 && cc < w && grid[cr][cc] == Box {
			boxes = append(boxes, Coord{cr, cc})
			cr += dr
			cc += dc
		}

		// Проверяем, можно ли сдвинуть последний коробку
		if cr < 0 || cr >= h || cc < 0 || cc >= w {
			return false
		}
		if grid[cr][cc] == Wall || grid[cr][cc] == Robot {
			return false
		}
		if grid[cr][cc] == Box {
			if !pushBoxes(grid, cr, cc, dr, dc, expanded) {
				return false
			}
		}

		// Сдвигаем коробки с конца
		for i := len(boxes) - 1; i >= 0; i-- {
			br, bc := boxes[i].R, boxes[i].C
			grid[br+dr][bc+dc] = Box
			grid[br][bc] = Empty
			// Отладочный вывод
			fmt.Printf("[DEBUG] Shifted Box from (%d, %d) to (%d, %d)\n", br, bc, br+dr, bc+dc)
		}
		return true
	}
}

// Вычисление суммы GPS координат
func calcSum(grid [][]Tile, expanded bool) int {
	sum := 0
	h, w := len(grid), len(grid[0])
	if !expanded {
		// Каждая коробка - единичный tile
		for r := 0; r < h; r++ {
			for c := 0; c < w; c++ {
				if grid[r][c] == Box {
					sum += 100*r + c
				}
			}
		}
	} else {
		// Коробки двойные по горизонтали
		// Считаем только левую часть каждой пары коробок
		for r := 0; r < h; r++ {
			c := 0
			for c < w {
				if grid[r][c] == Box {
					// Проверяем, является ли это левая часть пары
					if c+1 < w && grid[r][c+1] == Box {
						sum += 100*r + c
						c += 2 // Пропускаем правую часть пары
						continue
					}
				}
				c++
			}
		}
	}
	return sum
}

// Получение списка координат коробок (левых частей в расширенном режиме)
func getBoxPositions(grid [][]Tile, expanded bool) []Coord {
	var positions []Coord
	h, w := len(grid), len(grid[0])
	if !expanded {
		// Каждая коробка - единичный tile
		for r := 0; r < h; r++ {
			for c := 0; c < w; c++ {
				if grid[r][c] == Box {
					positions = append(positions, Coord{r, c})
				}
			}
		}
	} else {
		// Коробки двойные по горизонтали
		for r := 0; r < h; r++ {
			c := 0
			for c < w {
				if grid[r][c] == Box {
					// Проверяем, является ли это левая часть пары
					if c+1 < w && grid[r][c+1] == Box {
						positions = append(positions, Coord{r, c})
						c += 2 // Пропускаем правую часть пары
						continue
					}
				}
				c++
			}
		}
	}
	return positions
}

// Отладочный вывод карты (необязательно, можно раскомментировать в simulate для подробной отладки)
func printMap(grid [][]Tile, expanded bool) {
	for _, row := range grid {
		for _, tile := range row {
			switch tile {
			case Wall:
				fmt.Print("#")
			case Empty:
				fmt.Print(".")
			case Box:
				if expanded {
					fmt.Print("O") // Можно заменить на другой символ, если нужно
				} else {
					fmt.Print("O")
				}
			case Robot:
				fmt.Print("@")
			}
		}
		fmt.Println()
	}
}