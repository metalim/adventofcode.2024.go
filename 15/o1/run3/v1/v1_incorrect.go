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

// Структура карты
type warehouse struct {
	grid           [][]rune
	h, w           int
	robotR, robotC int
}

// Читает входной файл
func readInput(filename string) ([]string, []string) {
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

	// Отделим карту от последовательности движений
	// Карта окружена стенами (#), её конец можно найти по полной строке ###...
	// или первым встретившимся движением.
	// Но надежнее: найти первую строку без '#' по краям будет сложно. Предположим,
	// что после карты будет пустая строка, а затем ходы.
	var mapLines []string
	var moveLines []string
	phase := 0
	for _, l := range lines {
		if phase == 0 {
			if l == "" {
				phase = 1
			} else {
				mapLines = append(mapLines, l)
			}
		} else {
			moveLines = append(moveLines, l)
		}
	}

	// Удалим пустые строки в конце moveLines, склеим все в один
	var moves []string
	for _, ml := range moveLines {
		if ml != "" {
			moves = append(moves, ml)
		}
	}

	return mapLines, moves
}

// Ищет начальное положение робота
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

// Перемещение: ^, v, <, >
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

// Симуляция для части 1
func simulatePart1(w *warehouse, moves string) {
	for _, mv := range moves {
		dr, dc := dirDelta(mv)
		nr, nc := w.robotR+dr, w.robotC+dc
		if nr < 0 || nr >= w.h || nc < 0 || nc >= w.w {
			continue
		}
		if w.grid[nr][nc] == '#' {
			continue
		}
		if w.grid[nr][nc] == '.' {
			// Просто двигаем робота
			w.grid[w.robotR][w.robotC] = '.'
			w.grid[nr][nc] = '@'
			w.robotR, w.robotC = nr, nc
		} else if w.grid[nr][nc] == 'O' {
			// Пытаемся толкнуть коробку
			br, bc := nr+dr, nc+dc
			if br < 0 || br >= w.h || bc < 0 || bc >= w.w {
				continue
			}
			if w.grid[br][bc] == '#' || w.grid[br][bc] == 'O' || w.grid[br][bc] == '@' {
				// Сдвинуть нельзя
				continue
			}
			// Можно сдвинуть
			w.grid[br][bc] = 'O'
			w.grid[nr][nc] = '@'
			w.grid[w.robotR][w.robotC] = '.'
			w.robotR, w.robotC = nr, nc
		}
	}
}

// Подсчет суммы GPS координат для части 1
func sumGPSPart1(w *warehouse) int {
	// Найдем расстояния от верхнего и левого края.
	// Размер карты w.h x w.w
	// Для коробки GPS = 100*row + col
	// row,col начинаются с 0
	// Коробки 'O'
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

// Создание увеличенной карты для части 2
func scaleMap(lines []string) []string {
	// Правила замены:
	// # -> ##
	// O -> []
	// . -> ..
	// @ -> @.
	// Высота та же, ширина удваивается
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
			default:
				// На всякий случай, если встретится что-то чужое
				sb.WriteRune(ch)
				sb.WriteRune(ch)
			}
		}
		scaled = append(scaled, sb.String())
	}
	return scaled
}

// Поиск робота на увеличенной карте
func findRobotPart2(grid [][]rune) (int, int) {
	for r := 0; r < len(grid); r++ {
		for c := 0; c < len(grid[r]); c++ {
			if grid[r][c] == '@' {
				return r, c
			}
		}
	}
	panic("no robot in part2")
}

// Проверка, является ли позиция коробкой (2 символа '[]')
func isBoxPart2(grid [][]rune, r, c int) bool {
	// Коробка занимает две клетки по горизонтали: '[' и ']'
	// Проверим что в (r,c) - '[' и в (r,c+1) - ']'
	if c+1 < len(grid[r]) && grid[r][c] == '[' && grid[r][c+1] == ']' {
		return true
	}
	return false
}

// Проверка, стена ли это (два символа '#')
func isWallPart2(grid [][]rune, r, c int) bool {
	if c+1 < len(grid[r]) && grid[r][c] == '#' && grid[r][c+1] == '#' {
		return true
	}
	return false
}

// Проверка пустого места ".."
func isEmptyPart2(grid [][]rune, r, c int) bool {
	// Пустое место ".." или ".@" и т.п. для робота
	// Но робот '@.' - надо аккуратно
	// Пустым считаем ".." или ".@" или "@." или любой вариант без '[' ']' '#' для этой клетки.
	if grid[r][c] == '.' && c+1 < len(grid[r]) && (grid[r][c+1] == '.' || grid[r][c+1] == '@') {
		return true
	}
	if grid[r][c] == '@' && c+1 < len(grid[r]) && (grid[r][c+1] == '.') {
		// Клетка робота - не пустая, но для движения это так же как пустая?
		// Нет, робот - уникален. В момент движения робот смещается.
		// Нам нужно различать клетки с роботом и пустоту.
		return false
	}
	// Случай если это '.@'? Может ли быть '.@'? Нет, не по правилам конвертации.
	return false
}

// Проверка клетки на робот
func isRobotPart2(grid [][]rune, r, c int) bool {
	if grid[r][c] == '@' {
		return true
	}
	return false
}

// Для удобства определим функцию получения типа клетки в part2:
func getCellTypePart2(grid [][]rune, r, c int) cellType {
	if r < 0 || r >= len(grid) || c < 0 || c >= len(grid[0]) {
		return cellWall // за пределами - считаем стеной
	}

	// Определим приоритет:
	// Стена: '#' следом '#'
	if c+1 < len(grid[r]) && grid[r][c] == '#' && grid[r][c+1] == '#' {
		return cellWall
	}
	// Коробка: '[' ']'
	if isBoxPart2(grid, r, c) {
		return cellBox
	}
	// Робот: '@'
	if isRobotPart2(grid, r, c) {
		return cellRobot
	}
	// Пустота: ".." или ".@". Но ".@" не бывает, '@' всегда идет первым символом в паре?
	// См. условие: @->@. Значит робот всегда '@' '.'
	// '.' '.' - пусто
	if c+1 < len(grid[r]) && grid[r][c] == '.' && grid[r][c+1] == '.' {
		return cellEmpty
	}

	// Если встретились нестандартные пары, попробуем распознать:
	// Возможные варианты:
	// '@.' - робот
	// '[]' - коробка
	// '##' - стена
	// '..' - пусто
	// Других быть не должно
	return cellEmpty
}

// Установка клетки в part2
func setCellPart2(grid [][]rune, r, c int, ct cellType) {
	switch ct {
	case cellWall:
		// Не трогаем стены в процессе, они не меняются
	case cellEmpty:
		grid[r][c] = '.'
		grid[r][c+1] = '.'
	case cellBox:
		grid[r][c] = '['
		grid[r][c+1] = ']'
	case cellRobot:
		grid[r][c] = '@'
		grid[r][c+1] = '.'
	}
}

// Очистить клетку от содержимого, сделать пустой
func clearCellPart2(grid [][]rune, r, c int) {
	grid[r][c] = '.'
	grid[r][c+1] = '.'
}

// Двинуть робота в part2
func simulatePart2(grid [][]rune, moves string) (int, int) {
	rR, cR := findRobotPart2(grid)
	for _, mv := range moves {
		dr, dc := dirDelta(mv)
		nr, nc := rR+dr, cR+dc
		// Проверим границы
		if nr < 0 || nr >= len(grid) || nc < 0 || nc+1 >= len(grid[0]) {
			continue
		}

		// Определим что за клетка
		ct := getCellTypePart2(grid, nr, nc)
		if ct == cellWall {
			// Движение невозможно
			continue
		}
		if ct == cellEmpty {
			// Перемещение робота
			clearCellPart2(grid, rR, cR)
			setCellPart2(grid, nr, nc, cellRobot)
			rR, cR = nr, nc
			continue
		}
		if ct == cellRobot {
			// Теоретически не должно быть второго робота
			continue
		}
		if ct == cellBox {
			// Толкаем коробку
			br, bc := nr+dr, nc+dc
			if br < 0 || br >= len(grid) || bc < 0 || bc+1 >= len(grid[0]) {
				continue
			}
			nextCt := getCellTypePart2(grid, br, bc)
			if nextCt == cellWall || nextCt == cellBox || nextCt == cellRobot {
				// Не можем сдвинуть
				continue
			}
			// Сдвиг коробки
			// Стереть коробку с nr,nc
			clearCellPart2(grid, nr, nc)
			// Поставить коробку на br,bc
			setCellPart2(grid, br, bc, cellBox)
			// Робот переходит на nr,nc
			clearCellPart2(grid, rR, cR)
			setCellPart2(grid, nr, nc, cellRobot)
			rR, cR = nr, nc
		}
	}
	return rR, cR
}

// Подсчет GPS координат для части 2
// Теперь коробки '[”]' шириной 2. Расстояние считается до верхнего и левого края карты
// Край карты - это ряды и колонки с 0.
// GPS(box) = 100*r + c, где r,c - координаты левой части коробки '['
// Это просто индекс строки и столбца, где '[' стоит.
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
		panic("usage: program input.txt")
	}
	filename := os.Args[1]
	mapLines, moveLines := readInput(filename)
	moves := strings.Join(moveLines, "")

	// Сборка карты для части 1
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

	fmt.Println(part1Sum)
	fmt.Println(elapsed1)

	// Часть 2
	scaledLines := scaleMap(mapLines)
	var grid2 [][]rune
	for _, l := range scaledLines {
		grid2 = append(grid2, []rune(l))
	}

	start = time.Now()
	simulatePart2(grid2, moves)
	part2Sum := sumGPSPart2(grid2)
	elapsed2 := time.Since(start)

	fmt.Println(part2Sum)
	fmt.Println(elapsed2)
}
