package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type cellType int

const (
	cellWall cellType = iota
	cellEmpty
	cellBox
	cellRobot
)

func readInput(filename string) ([]string, string) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var lines []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	if err := sc.Err(); err != nil {
		panic(err)
	}

	// Найдем, где заканчивается карта
	width := len(lines[0])
	mapEnd := 0
	for i, l := range lines {
		if len(l) != width || strings.ContainsAny(l, "^v<>") {
			mapEnd = i
			break
		}
	}
	if mapEnd == 0 {
		mapEnd = len(lines)
	}

	mapLines := lines[:mapEnd]
	moveLines := lines[mapEnd:]
	var movesBuilder strings.Builder
	for _, ml := range moveLines {
		movesBuilder.WriteString(strings.TrimSpace(ml))
	}
	moves := movesBuilder.String()

	return mapLines, moves
}

type warehouse struct {
	grid           [][]rune
	h, w           int
	robotR, robotC int
}

func findRobot(m [][]rune) (int, int) {
	for r := 0; r < len(m); r++ {
		for c := 0; c < len(m[r]); c++ {
			if m[r][c] == '@' {
				return r, c
			}
		}
	}
	panic("no robot")
}

func dirDelta(d rune) (int, int) {
	switch d {
	case '^':
		return -1, 0
	case 'v':
		return 1, 0
	case '<':
		return 0, -1
	case '>':
		return 0, 1
	}
	return 0, 0
}

func tryMovePart1(w *warehouse, dr, dc int) {
	nr, nc := w.robotR+dr, w.robotC+dc
	if nr < 0 || nr >= w.h || nc < 0 || nc >= w.w {
		return
	}
	if w.grid[nr][nc] == '#' {
		return
	}
	if w.grid[nr][nc] == '.' {
		w.grid[w.robotR][w.robotC] = '.'
		w.grid[nr][nc] = '@'
		w.robotR, w.robotC = nr, nc
		return
	}

	if w.grid[nr][nc] == 'O' {
		boxPositions := make([][2]int, 0)
		cr, cc := nr, nc
		for {
			if cr < 0 || cr >= w.h || cc < 0 || cc >= w.w {
				break
			}
			if w.grid[cr][cc] == 'O' {
				boxPositions = append(boxPositions, [2]int{cr, cc})
				cr += dr
				cc += dc
			} else {
				break
			}
		}
		if cr < 0 || cr >= w.h || cc < 0 || cc >= w.w {
			return
		}
		if w.grid[cr][cc] != '.' {
			return
		}
		// Очищаем старые позиции
		w.grid[w.robotR][w.robotC] = '.'
		for _, pos := range boxPositions {
			w.grid[pos[0]][pos[1]] = '.'
		}
		// Итоговое состояние:
		// Робот становится на место первой коробки
		w.grid[boxPositions[0][0]][boxPositions[0][1]] = '@'
		// Цепочка коробок сдвигается вперёд
		for i := 0; i < len(boxPositions)-1; i++ {
			w.grid[boxPositions[i+1][0]][boxPositions[i+1][1]] = 'O'
		}
		// Последняя коробка становится на свободное место
		w.grid[cr][cc] = 'O'

		w.robotR, w.robotC = boxPositions[0][0], boxPositions[0][1]
	}
}

func simulatePart1(w *warehouse, moves string) {
	for _, mv := range moves {
		dr, dc := dirDelta(mv)
		tryMovePart1(w, dr, dc)
	}
}

func sumGPSPart1(w *warehouse) int {
	sum := 0
	for r := 0; r < w.h; r++ {
		for c := 0; c < w.w; c++ {
			if w.grid[r][c] == 'O' {
				sum += 100*r + c
			}
		}
	}
	return sum
}

// Часть 2
func scaleMap(lines []string) []string {
	var scaled []string
	for _, l := range lines {
		var sb strings.Builder
		for _, ch := range l {
			switch ch {
			case '#':
				sb.WriteString("##")
			case 'O':
				sb.WriteString("[]")
			case '.':
				sb.WriteString("..")
			case '@':
				sb.WriteString("@.")
			}
		}
		scaled = append(scaled, sb.String())
	}
	return scaled
}

func findRobotPart2(grid [][]rune) (int, int) {
	for r := 0; r < len(grid); r++ {
		for c := 0; c+1 < len(grid[r]); c++ {
			if grid[r][c] == '@' && grid[r][c+1] == '.' {
				return r, c
			}
		}
	}
	panic("no robot in part2")
}

func isWallPart2(grid [][]rune, r, c int) bool {
	if r < 0 || r >= len(grid) || c < 0 || c+1 >= len(grid[r]) {
		return true
	}
	return grid[r][c] == '#' && grid[r][c+1] == '#'
}

func isEmptyPart2(grid [][]rune, r, c int) bool {
	if r < 0 || r >= len(grid) || c < 0 || c+1 >= len(grid[r]) {
		return false
	}
	return grid[r][c] == '.' && grid[r][c+1] == '.'
}

func isBoxPart2(grid [][]rune, r, c int) bool {
	if r < 0 || r >= len(grid) || c < 0 || c+1 >= len(grid[r]) {
		return false
	}
	return grid[r][c] == '[' && grid[r][c+1] == ']'
}

func isRobotPart2(grid [][]rune, r, c int) bool {
	if r < 0 || r >= len(grid) || c < 0 || c+1 >= len(grid[r]) {
		return false
	}
	return grid[r][c] == '@' && grid[r][c+1] == '.'
}

func clearCellPart2(grid [][]rune, r, c int) {
	grid[r][c] = '.'
	grid[r][c+1] = '.'
}

func setRobotPart2(grid [][]rune, r, c int) {
	grid[r][c] = '@'
	grid[r][c+1] = '.'
}

func setBoxPart2(grid [][]rune, r, c int) {
	grid[r][c] = '['
	grid[r][c+1] = ']'
}

func tryMovePart2(grid [][]rune, rR, cR, dr, dc int) (int, int) {
	nr, nc := rR+dr, cR+dc
	if isWallPart2(grid, nr, nc) {
		return rR, cR
	}
	if isRobotPart2(grid, nr, nc) {
		return rR, cR
	}
	if isEmptyPart2(grid, nr, nc) {
		clearCellPart2(grid, rR, cR)
		setRobotPart2(grid, nr, nc)
		return nr, nc
	}
	if isBoxPart2(grid, nr, nc) {
		// Найдём цепочку коробок
		boxes := make([][2]int, 0)
		cr, cc := nr, nc
		for {
			if isBoxPart2(grid, cr, cc) {
				boxes = append(boxes, [2]int{cr, cc})
				cr += dr
				cc += dc
			} else {
				break
			}
		}
		if !isEmptyPart2(grid, cr, cc) {
			return rR, cR
		}

		// Очистим прежние позиции
		clearCellPart2(grid, rR, cR)
		for _, pos := range boxes {
			clearCellPart2(grid, pos[0], pos[1])
		}
		// Расставляем после сдвига:
		// Робот на место первой коробки
		setRobotPart2(grid, boxes[0][0], boxes[0][1])
		// Каждая коробка сдвигается на позицию следующей
		for i := 0; i < len(boxes)-1; i++ {
			setBoxPart2(grid, boxes[i+1][0], boxes[i+1][1])
		}
		// Последняя коробка на свободное место (cr,cc)
		setBoxPart2(grid, cr, cc)

		return boxes[0][0], boxes[0][1]
	}
	return rR, cR
}

func simulatePart2(grid [][]rune, moves string) (int, int) {
	rR, cR := findRobotPart2(grid)
	for _, mv := range moves {
		dr, dc := dirDelta(mv)
		nr, nc := tryMovePart2(grid, rR, cR, dr, dc)
		rR, cR = nr, nc
	}
	return rR, cR
}

func sumGPSPart2(grid [][]rune) int {
	sum := 0
	for r := 0; r < len(grid); r++ {
		for c := 0; c+1 < len(grid[r]); c++ {
			if grid[r][c] == '[' && grid[r][c+1] == ']' {
				sum += 100*r + c
			}
		}
	}
	return sum
}

func main() {
	if len(os.Args) < 2 {
		panic("usage: go run . input.txt")
	}
	filename := os.Args[1]
	mapLines, moves := readInput(filename)

	var grid1 [][]rune
	for _, l := range mapLines {
		grid1 = append(grid1, []rune(l))
	}
	w1 := &warehouse{
		grid: grid1,
		h:    len(grid1),
		w:    len(grid1[0]),
	}
	w1.robotR, w1.robotC = findRobot(w1.grid)

	start := time.Now()
	simulatePart1(w1, moves)
	part1Sum := sumGPSPart1(w1)
	elapsed1 := time.Since(start)

	scaledLines := scaleMap(mapLines)
	var grid2 [][]rune
	for _, l := range scaledLines {
		grid2 = append(grid2, []rune(l))
	}

	start = time.Now()
	simulatePart2(grid2, moves)
	part2Sum := sumGPSPart2(grid2)
	elapsed2 := time.Since(start)

	fmt.Printf("Part 1: %d           in %v\n", part1Sum, elapsed1)
	fmt.Printf("Part 2: %d            in %v\n", part2Sum, elapsed2)
}
