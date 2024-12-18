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
func findRegions(grid [][]byte) []([]point) {
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
				// BFS/DFS для региона
				var stack []point
				stack = append(stack, point{i, j})
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
	cells := make(map[point]bool)
	for _, p := range region {
		cells[p] = true
	}
	perim := 0
	for _, p := range region {
		// Считаем грани, не смежные с тем же регионом
		for _, d := range dirs {
			nr, nc := p.r+d.r, p.c+d.c
			if nr < 0 || nr >= R || nc < 0 || nc >= C || grid[nr][nc] != ch {
				perim++
			}
		}
	}
	return perim
}

// Для части 2 нам нужны "стороны" - количество прямолинейных секций ограждения.
// Подход:
// 1) Найти граничные сегменты.
// 2) Сформировать циклы границ региона.
// 3) Подсчитать количество сторон: двигаясь по циклу, считать количество изменений направления.
func calcSides(grid [][]byte, region []point) int {
	R, C := len(grid), len(grid[0])
	ch := grid[region[0].r][region[0].c]
	cells := make(map[point]bool)
	for _, p := range region {
		cells[p] = true
	}

	// Найдём все граничные ребра
	type edge struct {
		p1, p2 point
	}
	// Представим граничные ребра как набор вертикальных или горизонтальных сегментов между узлами сетки.
	// Узлы - это "углы" клеток. Для клетки (r,c) углы: (r,c), (r,c+1), (r+1,c), (r+1,c+1) в "сеточных координатах".
	// Если клетка принадлежит региону и сосед снаружи - ребро по границе.
	// Соберём все ребра в виде узлов (r,c) сетки вершин.
	edgesMap := make(map[[2]point]bool)
	addEdge := func(a, b point) {
		if a.r > b.r || (a.r == b.r && a.c > b.c) {
			a, b = b, a
		}
		edgesMap[[2]point{a, b}] = true
	}

	for _, p := range region {
		r, c := p.r, p.c
		// Проверяем каждое из 4 направлений, если за ней нет такого же региона - это граничное ребро
		for i, d := range dirs {
			nr, nc := r+d.r, c+d.c
			if nr < 0 || nr >= R || nc < 0 || nc >= C || grid[nr][nc] != ch {
				// ребро между вершинами в зависимости от направления
				// dirs: up, right, down, left
				// up   : ребро между (r,c) и (r,c+1)
				// right: ребро между (r,c+1) и (r+1,c+1)
				// down : ребро между (r+1,c) и (r+1,c+1)
				// left : ребро между (r,c) и (r+1,c)
				switch i {
				case 0: // up
					addEdge(point{r, c}, point{r, c + 1})
				case 1: // right
					addEdge(point{r, c + 1}, point{r + 1, c + 1})
				case 2: // down
					addEdge(point{r + 1, c}, point{r + 1, c + 1})
				case 3: // left
					addEdge(point{r, c}, point{r + 1, c})
				}
			}
		}
	}

	// Теперь у нас есть набор граничных ребер в виде соединённых точек. Нужно выделить все контуры.
	// Построим граф смежности: для каждой вершины список соседей.
	graph := make(map[point][]point)
	for e := range edgesMap {
		a, b := e[0], e[1]
		graph[a] = append(graph[a], b)
		graph[b] = append(graph[b], a)
	}

	visitedVert := make(map[point]bool)
	var cycles [][]point

	for v := range graph {
		if !visitedVert[v] {
			// Обход для выделения цикла(ов). Граничные ребра образуют один или несколько замкнутых контуров.
			// Можно идти по ребрам, пока не вернемся к стартовой точке.
			// Найдём циклы, двигаясь жадно: у каждой вершины в контуре должно быть 2 соседа (кроме вырожденных случаев).
			var cycle []point
			start := v
			cur := v
			var prev point
			for {
				visitedVert[cur] = true
				cycle = append(cycle, cur)
				// Найдём следующую вершину, не равную prev
				var next point
				for _, nx := range graph[cur] {
					if nx != prev {
						next = nx
						break
					}
				}
				if next == (point{}) {
					// нет продолжения
					break
				}
				prev = cur
				cur = next
				if cur == start {
					break
				}
			}
			// Проверим замкнутость
			if len(cycle) > 2 && cycle[0] == cycle[len(cycle)-1] {
				cycle = cycle[:len(cycle)-1]
			}
			cycles = append(cycles, cycle)
		}
	}

	// Подсчёт сторон для каждого цикла:
	// Пройдем по циклу и отметим направления переходов между вершинами
	// Для подсчёта направлений: если точки (r,c) идут в порядке цикла, то направление шага - (dr,dc).
	// Стороны - это количество отрезков прямого направления. Изменение направления увеличивает число сторон.
	// Изначально sides = 1. Каждый поворот увеличивает.
	calcCycleSides := func(cycle []point) int {
		if len(cycle) < 2 {
			return 0
		}
		// Найдем вектор предыдущего шага
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

	sides := 0
	for _, cyc := range cycles {
		if len(cyc) > 2 {
			sides += calcCycleSides(cyc)
		} else if len(cyc) == 2 {
			// Если вдруг контур из двух точек (линейный), это 1 сторона
			sides += 1
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
