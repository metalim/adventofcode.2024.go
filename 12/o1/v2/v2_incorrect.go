package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

type point struct {
	r, c int
}

var dirs = []point{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

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

func calcPerimeter(grid [][]byte, region []point) int {
	R, C := len(grid), len(grid[0])
	ch := grid[region[0].r][region[0].c]
	cells := make(map[point]bool, len(region))
	for _, p := range region {
		cells[p] = true
	}
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

// Оптимизированный подсчёт сторон для второй части
// Идея:
// 1) Собираем граничные рёбра.
// 2) Эти рёбра образуют один или несколько замкнутых контуров.
// 3) Для каждого контура считаем количество сторон (учитывая, что смежные коллинеарные сегменты объединяются в один).
func calcSides(grid [][]byte, region []point) int {
	R, C := len(grid), len(grid[0])
	ch := grid[region[0].r][region[0].c]
	cells := make(map[point]bool, len(region))
	for _, p := range region {
		cells[p] = true
	}

	// Собираем граничные ребра
	// Каждый ребро задаётся парой вершин (в "узловых" координатах)
	// Вершины - это углы клеток: для клетки (r,c) углы это (r,c), (r,c+1), (r+1,c), (r+1,c+1)
	type edge struct {
		a, b point
	}
	graph := make(map[point][]point)

	addEdge := func(a, b point) {
		graph[a] = append(graph[a], b)
		graph[b] = append(graph[b], a)
	}

	for _, p := range region {
		r, c := p.r, p.c
		for i, d := range dirs {
			nr, nc := r+d.r, c+d.c
			if nr < 0 || nr >= R || nc < 0 || nc >= C || grid[nr][nc] != ch {
				// Граничный ребро
				var v1, v2 point
				switch i {
				case 0: // up
					v1 = point{r, c}
					v2 = point{r, c + 1}
				case 1: // right
					v1 = point{r, c + 1}
					v2 = point{r + 1, c + 1}
				case 2: // down
					v1 = point{r + 1, c}
					v2 = point{r + 1, c + 1}
				case 3: // left
					v1 = point{r, c}
					v2 = point{r + 1, c}
				}
				addEdge(v1, v2)
			}
		}
	}

	visitedVert := make(map[point]bool)
	totalSides := 0

	// Функция подсчета сторон в цикле:
	// Пройдем по циклу вершин и подсчитаем количество поворотов.
	// Количество сторон = количество линейных сегментов = 1 + кол-во смен направления.
	calcCycleSides := func(cycle []point) int {
		if len(cycle) < 2 {
			return 0
		}
		// Найдем векторы направлений и посчитаем повороты
		// Замкнём цикл, добавим начальную вершину в конец для удобства
		cycle = append(cycle, cycle[0])
		sides := 1
		prevDr := cycle[1].r - cycle[0].r
		prevDc := cycle[1].c - cycle[0].c
		for i := 1; i < len(cycle)-1; i++ {
			dr := cycle[i+1].r - cycle[i].r
			dc := cycle[i+1].c - cycle[i].c
			if dr != prevDr || dc != prevDc {
				sides++
			}
			prevDr, prevDc = dr, dc
		}
		return sides
	}

	// Находим все циклы. У каждого узла степень 2 или 0 в корректном контуре.
	// Для каждого не посещенного узла пройдем по циклу
	for start := range graph {
		if visitedVert[start] {
			continue
		}

		// Пробуем выстроить цикл начиная с start
		cycle := []point{}
		cur := start
		var prev point
		for {
			visitedVert[cur] = true
			cycle = append(cycle, cur)
			neighbors := graph[cur]
			// Ищем следующий узел, не равный prev
			var next point
			for _, nx := range neighbors {
				if nx != prev {
					next = nx
					break
				}
			}
			if next == (point{}) {
				// Нет продолжения - не должно случиться в корректной задаче
				break
			}
			prev = cur
			cur = next
			if cur == start {
				break
			}
		}

		// Если получили цикл
		if len(cycle) > 2 && cycle[0] == cycle[len(cycle)-1] {
			cycle = cycle[:len(cycle)-1]
		}
		sides := calcCycleSides(cycle)
		totalSides += sides
	}

	return totalSides
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
