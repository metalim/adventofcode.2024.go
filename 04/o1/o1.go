// Напиши код на Go для решения следующей задачи. Входные данные в файле указываемом аргументом в командной строке.
// Выведи время решения каждой части.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <input_file>")
		return
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var grid []string
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if len(line) > 0 {
			grid = append(grid, line)
		}
	}
	if err := sc.Err(); err != nil {
		panic(err)
	}

	startPart1 := time.Now()
	part1Count := countXMAS(grid)
	elapsedPart1 := time.Since(startPart1)

	startPart2 := time.Now()
	part2Count := countX_MAS(grid)
	elapsedPart2 := time.Since(startPart2)

	fmt.Printf("Part 1: %d (Time: %v)\n", part1Count, elapsedPart1)
	fmt.Printf("Part 2: %d (Time: %v)\n", part2Count, elapsedPart2)
}

// Часть 1: Подсчет всех вхождений "XMAS" во всех 8 направлениях
func countXMAS(grid []string) int {
	word := "XMAS"
	count := 0
	rows := len(grid)
	if rows == 0 {
		return 0
	}
	cols := len(grid[0])
	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}}

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			for _, d := range dirs {
				if checkWord(grid, r, c, d[0], d[1], word) {
					count++
				}
			}
		}
	}
	return count
}

func checkWord(grid []string, sr, sc, dr, dc int, word string) bool {
	for i := 0; i < len(word); i++ {
		r := sr + dr*i
		c := sc + dc*i
		if r < 0 || r >= len(grid) || c < 0 || c >= len(grid[0]) {
			return false
		}
		if grid[r][c] != word[i] {
			return false
		}
	}
	return true
}

// Часть 2: Подсчет всех вхождений X-MAS.
// Шаблон (3x3):
//
//	M . S
//	. A .
//	M . S
//
// Каждая диагональ "MAS" или "SAM".
func countX_MAS(grid []string) int {
	count := 0
	rows := len(grid)
	if rows < 3 {
		return 0
	}
	cols := len(grid[0])
	if cols < 3 {
		return 0
	}

	for r := 1; r < rows-1; r++ {
		for c := 1; c < cols-1; c++ {
			if grid[r][c] != 'A' {
				continue
			}
			// Проверяем диагонали:
			// Левая верхняя и правая нижняя: (r-1,c-1), (r+1,c+1)
			// Правая верхняя и левая нижняя: (r-1,c+1), (r+1,c-1)

			tl := grid[r-1][c-1]
			br := grid[r+1][c+1]
			tr := grid[r-1][c+1]
			bl := grid[r+1][c-1]

			// Проверка, что одна диагональ формирует MAS или SAM и другая тоже.
			diag1 := string([]byte{tl, 'A', br})
			diag2 := string([]byte{tr, 'A', bl})

			if isMAS(diag1) && isMAS(diag2) {
				count++
			}
		}
	}

	return count
}

func isMAS(s string) bool {
	return s == "MAS" || s == "SAM"
}
