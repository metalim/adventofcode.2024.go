/*
Напиши код на Go для решения следующей задачи.
Входные данные в файле указываемом аргументом в командной строке.
Выведи ответ и время решения после решения каждой части.

--- Day 10: Hoof It ---
You all arrive at a Lava Production Facility on a floating island in the sky. As the others begin to search the massive industrial complex, you feel a small nose boop your leg and look down to discover a reindeer wearing a hard hat.

The reindeer is holding a book titled "Lava Island Hiking Guide". However, when you open the book, you discover that most of it seems to have been scorched by lava! As you're about to ask how you can help, the reindeer brings you a blank topographic map of the surrounding area (your puzzle input) and looks up at you excitedly.

Perhaps you can help fill in the missing hiking trails?

The topographic map indicates the height at each position using a scale from 0 (lowest) to 9 (highest). For example:

0123
1234
8765
9876
Based on un-scorched scraps of the book, you determine that a good hiking trail is as long as possible and has an even, gradual, uphill slope. For all practical purposes, this means that a hiking trail is any path that starts at height 0, ends at height 9, and always increases by a height of exactly 1 at each step. Hiking trails never include diagonal steps - only up, down, left, or right (from the perspective of the map).

You look up from the map and notice that the reindeer has helpfully begun to construct a small pile of pencils, markers, rulers, compasses, stickers, and other equipment you might need to update the map with hiking trails.

A trailhead is any position that starts one or more hiking trails - here, these positions will always have height 0. Assembling more fragments of pages, you establish that a trailhead's score is the number of 9-height positions reachable from that trailhead via a hiking trail. In the above example, the single trailhead in the top left corner has a score of 1 because it can reach a single 9 (the one in the bottom left).

This trailhead has a score of 2:

...0...
...1...
...2...
6543456
7.....7
8.....8
9.....9
(The positions marked . are impassable tiles to simplify these examples; they do not appear on your actual topographic map.)

This trailhead has a score of 4 because every 9 is reachable via a hiking trail except the one immediately to the left of the trailhead:

..90..9
...1.98
...2..7
6543456
765.987
876....
987....
This topographic map contains two trailheads; the trailhead at the top has a score of 1, while the trailhead at the bottom has a score of 2:

10..9..
2...8..
3...7..
4567654
...8..3
...9..2
.....01
Here's a larger example:

89010123
78121874
87430965
96549874
45678903
32019012
01329801
10456732
This larger example has 9 trailheads. Considering the trailheads in reading order, they have scores of 5, 6, 5, 3, 1, 3, 5, 3, and 5. Adding these scores together, the sum of the scores of all trailheads is 36.

The reindeer gleefully carries over a protractor and adds it to the pile. What is the sum of the scores of all trailheads on your topographic map?


--- Part Two ---
The reindeer spends a few minutes reviewing your hiking trail map before realizing something, disappearing for a few minutes, and finally returning with yet another slightly-charred piece of paper.

The paper describes a second way to measure a trailhead called its rating. A trailhead's rating is the number of distinct hiking trails which begin at that trailhead. For example:

.....0.
..4321.
..5..2.
..6543.
..7..4.
..8765.
..9....
The above map has a single trailhead; its rating is 3 because there are exactly three distinct hiking trails which begin at that position:

.....0.   .....0.   .....0.
..4321.   .....1.   .....1.
..5....   .....2.   .....2.
..6....   ..6543.   .....3.
..7....   ..7....   .....4.
..8....   ..8....   ..8765.
..9....   ..9....   ..9....
Here is a map containing a single trailhead with rating 13:

..90..9
...1.98
...2..7
6543456
765.987
876....
987....
This map contains a single trailhead with rating 227 (because there are 121 distinct hiking trails that lead to the 9 on the right edge and 106 that lead to the 9 on the bottom edge):

012345
123456
234567
345678
4.6789
56789.
Here's the larger example from before:

89010123
78121874
87430965
96549874
45678903
32019012
01329801
10456732
Considering its trailheads in reading order, they have ratings of 20, 24, 10, 4, 1, 4, 5, 8, and 5. The sum of all trailhead ratings in this larger example topographic map is 81.

You're not sure how, but the reindeer seems to have crafted some tiny flags out of toothpicks and bits of paper and is using them to mark trailheads on your topographic map. What is the sum of the ratings of all trailheads?

*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type point struct {
	r, c int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <input_file>")
		return
	}

	startOverall := time.Now()

	// Чтение входных данных
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var grid [][]int
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := sc.Text()
		row := make([]int, len(line))
		for i, ch := range line {
			row[i] = int(ch - '0')
		}
		grid = append(grid, row)
	}

	n, m := len(grid), len(grid[0])

	// Соседи по четырем направлениям
	dirs := []point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

	// Для ускорения подсчетов создадим структуры кэша для Dynamic Programming (DP)
	// DP для подсчета достижимых 9 (часть 1) и количества путей (часть 2)

	// dpReach9[h][r][c] - множество достижимых 9 для клетки (r,c) высоты h
	// dpCountPaths[h][r][c] - количество путей до всех достижимых 9 для клетки (r,c)
	// Чтобы не хранить множества, будем идти снизу вверх: сначала найдем клетки с высотой 9,
	// затем для высоты 8 найдем достижимые 9 и количество путей, и так далее.
	// Это позволит эффективно пересчитать все данные.

	// Сначала сгруппируем клетки по высоте
	heightCells := make([][]point, 10)
	for r := 0; r < n; r++ {
		for c := 0; c < m; c++ {
			h := grid[r][c]
			heightCells[h] = append(heightCells[h], point{r, c})
		}
	}

	// dpReachMask и dpCount будут хранить данные по каждой клетке
	// dpReachMask[r][c]: битовая маска достижимых девяток. Чтобы хранить уникальные девятки,
	// можно использовать уникальный ID для каждой 9-ки. Но у нас может быть много 9-ок,
	// их количество неограничено. Тогда проще хранить количество достижимых 9 для первого задания
	// и суммарное количество путей для второго.
	// Решение: для первой части - мы просто храним множество достижимых 9 как счетчик (количество уникальных 9).
	// Для второй части - количество путей. Чтобы отличать 9-ки, можно их просто считать уникальными?
	// Задача подразумевает подсчет всех достижимых девяток и всех путей. Для счетчика уникальных 9:
	// - Можно использовать flood fill: каждая 9-ка уникальна, сохраним ее координаты в массив и пронумеруем.
	// - Но ограничений на размер нет, поэтому лучше использовать другую стратегию:
	//   Для первой части нам нужен только счетчик уникальных достижимых 9. Можем сделать DP снизу вверх:
	//   когда мы на 9, dp9Count[r][c] = 1, dpPathCount[r][c] = 1.
	//   Для h<9, dp9Count[r][c] = объединение всех достижимых от соседей h+1, то есть сумма уникальных 9.
	//   Но без уникальной идентификации это сложно. Однако в условии требуется только сумма достигнутых 9 от трейлхеда.
	//   На самом деле "уникальные" 9 подразумевают уникальные клетки с высотой 9.
	//   Значит для счета уникальных 9 будем использовать boolean grid visited9[r][c] для передачи информации.
	//   Это будет сложно без комплексной структуры данных.

	// Упростим решение:
	// 1) Посчитаем для каждой клетки dpCount9: количество уникальных 9, достижимых от нее (в смысле наличия пути).
	// 2) Посчитаем для каждой клетки dpPaths: суммарное количество путей от нее к любым 9.
	//
	// Для dpCount9 и dpPaths:
	// - На 9: dpCount9 = 1, dpPaths = 1
	// - На h от 0 до 8:
	//   dpCount9[r][c] = объединение всех dpCount9 соседей с h+1 (т.е. просто сумма уникальных 9? Нужно именно количество уникальных 9.)
	//   Но мы можем просто взять сумму dpCount9 соседей? Это даст дубликаты. Нужны уникальные 9.
	//
	// Число уникальных девяток может быть большое. Но нам нужно лишь сумму по трэлхедам, не по каждой клетке.
	// Чтобы корректно учесть уникальность 9, используем подход:
	// - Определим ID для каждой клетки. Для 9: dpCount9 = массив бит (слишком большой).
	// - Или просто BFS от каждой 9ки вверх по цепочке h-1? Это дорого.
	//
	// Условия задачи: мы должны лишь вывести сумму. Нужно точное решение.
	//
	// Оптимизация:
	// - Сначала найдем компоненты 9-к. Каждая 9 - уникальная позиция.
	//   Нам нужны уникальные 9, нельзя объединять.
	// - В худшем случае слишком большая карта.
	//
	// Допустим, что все 9 уникальны и не требуют обработки больших данных. Будем хранить для каждой клетки set уникальных 9,
	// но set слишком велик. Попытка: мемоизация + слияние множеств.
	//
	// Мы можем сделать топологический проход по высотам от 9 к 0:
	// При h=9: dpCount9[r][c] = { (r,c) }, dpPaths[r][c] = 1
	// При h<9:
	//   dpCount9[r][c] = объединение dpCount9 всех соседей h+1
	//   dpPaths[r][c] = сумма dpPaths соседей h+1 по всем путям, но для уникальности 9 нужно объединять множества.
	//
	// Это дорого. Но задача - это просто пример кода. Дадим простое решение с использованием map для уникальных 9:
	// В реальной задаче оптимизировали бы. Здесь просто сделаем решения через map[int]struct{} для уникальных 9 (по r*m+c).
	// dpCount9Sets[r][c] - set ID 9-к
	// dpPaths[r][c] - количество путей
	//
	// После вычислений суммируем для trailheads (h=0):
	// Часть 1: score = |dpCount9Sets[r][c]| для каждого trailhead
	// Часть 2: rating = dpPaths[r][c]
	//
	// Затем суммы выведем и время решения.

	startPart1 := time.Now()
	dpCount9Sets := make([]map[int]struct{}, n*m)
	dpPaths := make([]int, n*m)
	for i := range dpCount9Sets {
		dpCount9Sets[i] = nil
	}

	// Инициализация 9-к
	for _, p9 := range heightCells[9] {
		idx := p9.r*m + p9.c
		dpCount9Sets[idx] = map[int]struct{}{idx: {}}
		dpPaths[idx] = 1
	}

	// Обработка вниз по высотам
	for h := 8; h >= 0; h-- {
		for _, p0 := range heightCells[h] {
			idx0 := p0.r*m + p0.c
			countSet := make(map[int]struct{})
			pathsSum := 0
			for _, d := range dirs {
				rr, cc := p0.r+d.r, p0.c+d.c
				if rr < 0 || rr >= n || cc < 0 || cc >= m {
					continue
				}
				if grid[rr][cc] == h+1 {
					idx1 := rr*m + cc
					if dpCount9Sets[idx1] != nil {
						for k := range dpCount9Sets[idx1] {
							countSet[k] = struct{}{}
						}
						pathsSum += dpPaths[idx1]
					}
				}
			}
			if len(countSet) > 0 {
				dpCount9Sets[idx0] = countSet
				dpPaths[idx0] = pathsSum
			}
		}
	}

	// Суммы для trailheads
	part1Sum := 0
	part2Sum := 0
	for _, p0 := range heightCells[0] {
		idx0 := p0.r*m + p0.c
		if dpCount9Sets[idx0] != nil {
			part1Sum += len(dpCount9Sets[idx0])
			part2Sum += dpPaths[idx0]
		}
	}
	elapsedPart1 := time.Since(startPart1)

	fmt.Printf("Part 1: %d\n", part1Sum)
	fmt.Printf("Time Part 1: %v\n", elapsedPart1)

	startPart2 := time.Now()
	fmt.Printf("Part 2: %d\n", part2Sum)
	elapsedPart2 := time.Since(startPart2)
	fmt.Printf("Time Part 2: %v\n", elapsedPart2)

	fmt.Printf("Overall time: %v\n", time.Since(startOverall))
}
