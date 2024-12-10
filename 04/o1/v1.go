/*
Напиши код на Go для решения следующей задачи. Входные данные в файле указываемом аргументом в командной строке.
Выведи время решения каждой части.

--- Day 4: Ceres Search ---
"Looks like the Chief's not here. Next!" One of The Historians pulls out a device and pushes the only button on it. After a brief flash, you recognize the interior of the Ceres monitoring station!

As the search for the Chief continues, a small Elf who lives on the station tugs on your shirt; she'd like to know if you could help her with her word search (your puzzle input). She only has to find one word: XMAS.

This word search allows words to be horizontal, vertical, diagonal, written backwards, or even overlapping other words. It's a little unusual, though, as you don't merely need to find one instance of XMAS - you need to find all of them. Here are a few ways XMAS might appear, where irrelevant characters have been replaced with .:


..X...
.SAMX.
.A..A.
XMAS.S
.X....
The actual word search will be full of letters instead. For example:

MMMSXXMASM
MSAMXMSMSA
AMXSXMAAMM
MSAMASMSMX
XMASAMXAMM
XXAMMXXAMA
SMSMSASXSS
SAXAMASAAA
MAMMMXMMMM
MXMXAXMASX
In this word search, XMAS occurs a total of 18 times; here's the same word search again, but where letters not involved in any XMAS have been replaced with .:

....XXMAS.
.SAMXMS...
...S..A...
..A.A.MS.X
XMASAMX.MM
X.....XA.A
S.S.S.S.SS
.A.A.A.A.A
..M.M.M.MM
.X.X.XMASX
Take a look at the little Elf's word search. How many times does XMAS appear?


--- Part Two ---
The Elf looks quizzically at you. Did you misunderstand the assignment?

Looking for the instructions, you flip over the word search to find that this isn't actually an XMAS puzzle; it's an X-MAS puzzle in which you're supposed to find two MAS in the shape of an X. One way to achieve that is like this:

M.S
.A.
M.S
Irrelevant characters have again been replaced with . in the above diagram. Within the X, each MAS can be written forwards or backwards.

Here's the same example from before, but this time all of the X-MASes have been kept instead:

.M.S......
..A..MSMS.
.M.S.MAA..
..A.ASMSM.
.M.S.M....
..........
S.S.S.S.S.
.A.A.A.A..
M.M.M.M.M.
..........
In this example, an X-MAS appears 9 times.

Flip the word search from the instructions back over to the word search side and try again. How many times does an X-MAS appear?

*/

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
