/*

Напиши код на Go для решения следующей задачи.
Входные данные в файле указываемом аргументом в командной строке.
Выведи ответ и время решения после решения каждой части.

https://adventofcode.com/2024/day/12

--- Day 12: Garden Groups ---
Why not search for the Chief Historian near the gardener and his massive farm? There's plenty of food, so The Historians grab something to eat while they search.

You're about to settle near a complex arrangement of garden plots when some Elves ask if you can lend a hand. They'd like to set up fences around each region of garden plots, but they can't figure out how much fence they need to order or how much it will cost. They hand you a map (your puzzle input) of the garden plots.

Each garden plot grows only a single type of plant and is indicated by a single letter on your map. When multiple garden plots are growing the same type of plant and are touching (horizontally or vertically), they form a region. For example:

AAAA
BBCD
BBCC
EEEC
This 4x4 arrangement includes garden plots growing five different types of plants (labeled A, B, C, D, and E), each grouped into their own region.

In order to accurately calculate the cost of the fence around a single region, you need to know that region's area and perimeter.

The area of a region is simply the number of garden plots the region contains. The above map's type A, B, and C plants are each in a region of area 4. The type E plants are in a region of area 3; the type D plants are in a region of area 1.

Each garden plot is a square and so has four sides. The perimeter of a region is the number of sides of garden plots in the region that do not touch another garden plot in the same region. The type A and C plants are each in a region with perimeter 10. The type B and E plants are each in a region with perimeter 8. The lone D plot forms its own region with perimeter 4.

Visually indicating the sides of plots in each region that contribute to the perimeter using - and |, the above map's regions' perimeters are measured as follows:

+-+-+-+-+
|A A A A|
+-+-+-+-+     +-+
              |D|
+-+-+   +-+   +-+
|B B|   |C|
+   +   + +-+
|B B|   |C C|
+-+-+   +-+ +
          |C|
+-+-+-+   +-+
|E E E|
+-+-+-+
Plants of the same type can appear in multiple separate regions, and regions can even appear within other regions. For example:

OOOOO
OXOXO
OOOOO
OXOXO
OOOOO
The above map contains five regions, one containing all of the O garden plots, and the other four each containing a single X plot.

The four X regions each have area 1 and perimeter 4. The region containing 21 type O plants is more complicated; in addition to its outer edge contributing a perimeter of 20, its boundary with each X region contributes an additional 4 to its perimeter, for a total perimeter of 36.

Due to "modern" business practices, the price of fence required for a region is found by multiplying that region's area by its perimeter. The total price of fencing all regions on a map is found by adding together the price of fence for every region on the map.

In the first example, region A has price 4 * 10 = 40, region B has price 4 * 8 = 32, region C has price 4 * 10 = 40, region D has price 1 * 4 = 4, and region E has price 3 * 8 = 24. So, the total price for the first example is 140.

In the second example, the region with all of the O plants has price 21 * 36 = 756, and each of the four smaller X regions has price 1 * 4 = 4, for a total price of 772 (756 + 4 + 4 + 4 + 4).

Here's a larger example:

RRRRIICCFF
RRRRIICCCF
VVRRRCCFFF
VVRCCCJFFF
VVVVCJJCFE
VVIVCCJJEE
VVIIICJJEE
MIIIIIJJEE
MIIISIJEEE
MMMISSJEEE
It contains:

A region of R plants with price 12 * 18 = 216.
A region of I plants with price 4 * 8 = 32.
A region of C plants with price 14 * 28 = 392.
A region of F plants with price 10 * 18 = 180.
A region of V plants with price 13 * 20 = 260.
A region of J plants with price 11 * 20 = 220.
A region of C plants with price 1 * 4 = 4.
A region of E plants with price 13 * 18 = 234.
A region of I plants with price 14 * 22 = 308.
A region of M plants with price 5 * 12 = 60.
A region of S plants with price 3 * 8 = 24.
So, it has a total price of 1930.

What is the total price of fencing all regions on your map?

--- Part Two ---
Fortunately, the Elves are trying to order so much fence that they qualify for a bulk discount!

Under the bulk discount, instead of using the perimeter to calculate the price, you need to use the number of sides each region has. Each straight section of fence counts as a side, regardless of how long it is.

Consider this example again:

AAAA
BBCD
BBCC
EEEC
The region containing type A plants has 4 sides, as does each of the regions containing plants of type B, D, and E. However, the more complex region containing the plants of type C has 8 sides!

Using the new method of calculating the per-region price by multiplying the region's area by its number of sides, regions A through E have prices 16, 16, 32, 4, and 12, respectively, for a total price of 80.

The second example above (full of type X and O plants) would have a total price of 436.

Here's a map that includes an E-shaped region full of type E plants:

EEEEE
EXXXX
EEEEE
EXXXX
EEEEE
The E-shaped region has an area of 17 and 12 sides for a price of 204. Including the two regions full of type X plants, this map has a total price of 236.

This map has a total price of 368:

AAAAAA
AAABBA
AAABBA
ABBAAA
ABBAAA
AAAAAA
It includes two regions full of type B plants (each with 4 sides) and a single region full of type A plants (with 4 sides on the outside and 8 more sides on the inside, a total of 12 sides). Be especially careful when counting the fence around regions like the one full of type A plants; in particular, each section of fence has an in-side and an out-side, so the fence does not connect across the middle of the region (where the two B regions touch diagonally). (The Elves would have used the Möbius Fencing Company instead, but their contract terms were too one-sided.)

The larger example from before now has the following updated prices:

A region of R plants with price 12 * 10 = 120.
A region of I plants with price 4 * 4 = 16.
A region of C plants with price 14 * 22 = 308.
A region of F plants with price 10 * 12 = 120.
A region of V plants with price 13 * 10 = 130.
A region of J plants with price 11 * 12 = 132.
A region of C plants with price 1 * 4 = 4.
A region of E plants with price 13 * 8 = 104.
A region of I plants with price 14 * 16 = 224.
A region of M plants with price 5 * 6 = 30.
A region of S plants with price 3 * 6 = 18.
Adding these together produces its new total price of 1206.

What is the new total price of fencing all regions on your map?

*/

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
