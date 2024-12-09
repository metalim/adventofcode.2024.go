// Напиши код на Go для решения следующей задачи. Входные данные в файле указываемом аргументом в командной строке.
// Выведи время решения каждой части.

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

type point struct{ x, y int }

func gcd(a, b int) int {
	if b == 0 {
		if a < 0 {
			return -a
		}
		return a
	}
	return gcd(b, a%b)
}

// Возвращает множество точек-антинодов для части 1
// Условие: Для каждой пары антенн одной частоты считаем две точки:
// P1 - V и P1 + 2*V, где V = P2 - P1
func part1(positions map[rune][]point, maxX, maxY int) map[point]struct{} {
	antinodes := make(map[point]struct{})
	for freq, ants := range positions {
		if len(ants) < 2 {
			continue
		}
		for i := 0; i < len(ants); i++ {
			for j := i + 1; j < len(ants); j++ {
				p1, p2 := ants[i], ants[j]
				dx, dy := p2.x-p1.x, p2.y-p1.y
				// Антиноды: P1 - (dx, dy) и P1 + 2*(dx, dy)
				candidates := []point{
					{x: p1.x - dx, y: p1.y - dy},
					{x: p1.x + 2*dx, y: p1.y + 2*dy},
				}
				for _, c := range candidates {
					if c.x >= 0 && c.y >= 0 && c.y <= maxY && c.x <= maxX {
						antinodes[c] = struct{}{}
					}
				}
			}
		}
	}
	return antinodes
}

// Возвращает множество точек-антинодов для части 2
// Условие: Любая точка на линии, проходящей через хотя бы две антенны одной частоты
func part2(positions map[rune][]point, maxX, maxY int) map[point]struct{} {
	antinodes := make(map[point]struct{})
	for freq, ants := range positions {
		if len(ants) < 2 {
			continue
		}
		// Для каждой пары антенн получаем направляющий вектор (dx_min, dy_min)
		// и генерируем все точки на линии в пределах карты
		for i := 0; i < len(ants); i++ {
			for j := i + 1; j < len(ants); j++ {
				p1, p2 := ants[i], ants[j]
				dx, dy := p2.x-p1.x, p2.y-p1.y
				g := gcd(dx, dy)
				dx /= g
				dy /= g
				// Генерируем точки по линии в обе стороны от p1
				// В одну сторону
				for t := 0; ; t++ {
					x := p1.x + t*dx
					y := p1.y + t*dy
					if x < 0 || x > maxX || y < 0 || y > maxY {
						break
					}
					antinodes[point{x, y}] = struct{}{}
				}
				// В обратную сторону
				for t := -1; ; t-- {
					x := p1.x + t*dx
					y := p1.y + t*dy
					if x < 0 || x > maxX || y < 0 || y > maxY {
						break
					}
					antinodes[point{x, y}] = struct{}{}
				}
			}
		}
	}
	return antinodes
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <input_file>")
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	var grid []string
	for sc.Scan() {
		grid = append(grid, sc.Text())
	}
	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}

	maxY := len(grid) - 1
	maxX := 0
	if len(grid) > 0 {
		maxX = len(grid[0]) - 1
	}

	positions := make(map[rune][]point)
	for y, line := range grid {
		for x, ch := range line {
			if ch != '.' {
				positions[ch] = append(positions[ch], point{x, y})
			}
		}
	}

	start1 := time.Now()
	a1 := part1(positions, maxX, maxY)
	part1Time := time.Since(start1)

	start2 := time.Now()
	a2 := part2(positions, maxX, maxY)
	part2Time := time.Since(start2)

	fmt.Println(len(a1))
	fmt.Println(len(a2))
	fmt.Println("Part1 time:", part1Time)
	fmt.Println("Part2 time:", part2Time)
}
