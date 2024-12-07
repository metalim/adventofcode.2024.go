// Напиши код на Go для решения следующей задачи. Входные данные в файле указываемом аргументом в командной строке.
// Выведи время решения каждой части.

package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type posDir struct {
	r, c, d int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: program input.txt")
		return
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	var grid []string
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		grid = append(grid, sc.Text())
	}
	if err := sc.Err(); err != nil {
		fmt.Println("Error:", err)
		return
	}

	rows := len(grid)
	cols := len(grid[0])
	dirs := [4][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

	startR, startC, startD := -1, -1, -1
FindStart:
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			switch grid[i][j] {
			case '^':
				startR, startC, startD = i, j, 0
				break FindStart
			case '>':
				startR, startC, startD = i, j, 1
				break FindStart
			case 'v':
				startR, startC, startD = i, j, 2
				break FindStart
			case '<':
				startR, startC, startD = i, j, 3
				break FindStart
			}
		}
	}

	original := make([][]byte, rows)
	for i := range grid {
		original[i] = []byte(grid[i])
	}

	// Часть 1
	part1Start := time.Now()
	{
		r, c, d := startR, startC, startD
		visited := make(map[[2]int]bool)
		visited[[2]int{r, c}] = true

		for {
			nr, nc := r+dirs[d][0], c+dirs[d][1]
			if nr < 0 || nr >= rows || nc < 0 || nc >= cols {
				break
			}
			if original[nr][nc] == '#' {
				d = (d + 1) % 4
			} else {
				r, c = nr, nc
				visited[[2]int{r, c}] = true
			}
		}

		fmt.Println(len(visited)) // Результат части 1
	}
	part1Time := time.Since(part1Start)
	fmt.Println("Part1 time:", part1Time)

	// Часть 2
	part2Start := time.Now()
	simWithObstacle := func(or, oc int) bool {
		original[or][oc] = '#'
		r, c, d := startR, startC, startD
		visitedStates := make(map[posDir]bool)
		visitedStates[posDir{r, c, d}] = true

		for {
			nr, nc := r+dirs[d][0], c+dirs[d][1]
			if nr < 0 || nr >= rows || nc < 0 || nc >= cols {
				original[or][oc] = '.'
				return false
			}
			if original[nr][nc] == '#' {
				d = (d + 1) % 4
			} else {
				r, c = nr, nc
				st := posDir{r, c, d}
				if visitedStates[st] {
					original[or][oc] = '.'
					return true
				}
				visitedStates[st] = true
			}
		}
	}

	count := 0
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if (i == startR && j == startC) || original[i][j] != '.' {
				continue
			}
			if simWithObstacle(i, j) {
				count++
			}
		}
	}

	fmt.Println(count) // Результат части 2
	part2Time := time.Since(part2Start)
	fmt.Println("Part2 time:", part2Time)
}
