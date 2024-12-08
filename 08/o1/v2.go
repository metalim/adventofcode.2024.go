// o1/o1.go:31:6: declared and not used: freq
// o1/o1.go:59:6: declared and not used: freq

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
func part1(positions map[rune][]point, maxX, maxY int) map[point]struct{} {
	antinodes := make(map[point]struct{})
	for _, ants := range positions {
		if len(ants) < 2 {
			continue
		}
		for i := 0; i < len(ants); i++ {
			for j := i + 1; j < len(ants); j++ {
				p1, p2 := ants[i], ants[j]
				dx, dy := p2.x-p1.x, p2.y-p1.y
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
func part2(positions map[rune][]point, maxX, maxY int) map[point]struct{} {
	antinodes := make(map[point]struct{})
	for _, ants := range positions {
		if len(ants) < 2 {
			continue
		}
		for i := 0; i < len(ants); i++ {
			for j := i + 1; j < len(ants); j++ {
				p1, p2 := ants[i], ants[j]
				dx, dy := p2.x-p1.x, p2.y-p1.y
				g := gcd(dx, dy)
				dx /= g
				dy /= g
				for t := 0; ; t++ {
					x := p1.x + t*dx
					y := p1.y + t*dy
					if x < 0 || x > maxX || y < 0 || y > maxY {
						break
					}
					antinodes[point{x, y}] = struct{}{}
				}
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
