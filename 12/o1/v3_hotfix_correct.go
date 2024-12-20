package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

// Структура для координат
type point struct {
	r, c int
}

// Направления для обхода
var dirs = []point{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

// Читаем карту
func readMap(filename string) ([][]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var grid [][]byte
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if line != "" {
			grid = append(grid, []byte(line))
		}
	}
	return grid, sc.Err()
}

// Поиск регионов
func findRegions(grid [][]byte) [][]point {
	R, C := len(grid), len(grid[0])
	visited := make([][]bool, R)
	for i := range visited {
		visited[i] = make([]bool, C)
	}
	var regions [][]point

	for i := 0; i < R; i++ {
		for j := 0; j < C; j++ {
			if !visited[i][j] {
				ch := grid[i][j]
				stack := []point{{i, j}}
				visited[i][j] = true
				var region []point
				for len(stack) > 0 {
					cur := stack[len(stack)-1]
					stack = stack[:len(stack)-1]
					region = append(region, cur)
					for _, d := range dirs {
						nr, nc := cur.r+d.r, cur.c+d.c
						if nr >= 0 && nr < R && nc >= 0 && nc < C {
							if !visited[nr][nc] && grid[nr][nc] == ch {
								visited[nr][nc] = true
								stack = append(stack, point{nr, nc})
							}
						}
					}
				}
				regions = append(regions, region)
			}
		}
	}
	return regions
}

// Подсчет периметра для части 1
func calcPerimeter(grid [][]byte, region []point) int {
	R, C := len(grid), len(grid[0])
	ch := grid[region[0].r][region[0].c]
	perim := 0
	for _, p := range region {
		for _, d := range dirs {
			nr, nc := p.r+d.r, p.c+d.c
			if nr < 0 || nr >= R || nc < 0 || nc >= C || grid[nr][nc] != ch {
				perim++
			}
		}
	}
	return perim
}

// Подсчет сторон для части 2
func calcSides(grid [][]byte, region []point) int {
	R, C := len(grid), len(grid[0])
	// ch := grid[region[0].r][region[0].c] // manually commented out unused variable

	// Создаём карту для быстрого доступа к клеткам региона
	cells := make([][]bool, R)
	for i := 0; i < R; i++ {
		cells[i] = make([]bool, C)
	}
	for _, p := range region {
		cells[p.r][p.c] = true
	}

	sides := 0

	// Считаем горизонтальные стороны
	for r := 0; r < R; r++ {
		inSegment := false
		for c := 0; c < C; c++ {
			if cells[r][c] {
				// Проверяем верхнюю грань
				if r == 0 || !cells[r-1][c] {
					if !inSegment {
						sides++
						inSegment = true
					}
				} else {
					inSegment = false
				}
			} else {
				inSegment = false
			}
		}
	}

	// Считаем нижние горизонтальные стороны
	for r := 0; r < R; r++ {
		inSegment := false
		for c := 0; c < C; c++ {
			if cells[r][c] {
				// Проверяем нижнюю грань
				if r == R-1 || !cells[r+1][c] {
					if !inSegment {
						sides++
						inSegment = true
					}
				} else {
					inSegment = false
				}
			} else {
				inSegment = false
			}
		}
	}

	// Считаем левые вертикальные стороны
	for c := 0; c < C; c++ {
		inSegment := false
		for r := 0; r < R; r++ {
			if cells[r][c] {
				// Проверяем левую грань
				if c == 0 || !cells[r][c-1] {
					if !inSegment {
						sides++
						inSegment = true
					}
				} else {
					inSegment = false
				}
			} else {
				inSegment = false
			}
		}
	}

	// Считаем правые вертикальные стороны
	for c := 0; c < C; c++ {
		inSegment := false
		for r := 0; r < R; r++ {
			if cells[r][c] {
				// Проверяем правую грань
				if c == C-1 || !cells[r][c+1] {
					if !inSegment {
						sides++
						inSegment = true
					}
				} else {
					inSegment = false
				}
			} else {
				inSegment = false
			}
		}
	}

	return sides
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <input_file>")
	}
	filename := os.Args[1]

	grid, err := readMap(filename)
	if err != nil {
		log.Fatal(err)
	}

	// Проверка, что карта не пустая
	if len(grid) == 0 || len(grid[0]) == 0 {
		log.Fatal("Empty grid")
	}

	startPart1 := time.Now()
	regions := findRegions(grid)
	part1 := 0
	for _, r := range regions {
		area := len(r)
		perim := calcPerimeter(grid, r)
		part1 += area * perim
	}
	durationPart1 := time.Since(startPart1)
	fmt.Println("Part 1:", part1)
	fmt.Println("Time Part 1:", durationPart1)

	startPart2 := time.Now()
	part2 := 0
	for _, r := range regions {
		area := len(r)
		sides := calcSides(grid, r)
		part2 += area * sides
	}
	durationPart2 := time.Since(startPart2)
	fmt.Println("Part 2:", part2)
	fmt.Println("Time Part 2:", durationPart2)
}
