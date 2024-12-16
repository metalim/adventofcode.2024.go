package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

var Print = false

func catch(err error) {
	if err != nil {
		panic(err)
	}
}

type Guard struct {
	y, x int
	dir  int
}
type Vec2 struct {
	y, x int
}

var Directions = []Vec2{
	{-1, 0}, // North
	{0, -1}, // West
	{1, 0},  // South
	{0, 1},  // East
}

func main() {
	flag.BoolVar(&Print, "print", false, "Print the grid")
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("Usage: go run main.go input.txt")
		os.Exit(1)
	}

	bs, err := os.ReadFile(os.Args[1])
	catch(err)

	grid, guard := parseInput(string(bs))
	part1(grid, guard)
	part2(grid, guard)
}

func parseInput(input string) (grid [][]rune, guard Guard) {
	lines := strings.Split(input, "\n")
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}

	var guardFound bool
	for y, line := range lines {
		grid = append(grid, []rune(line))
		x := strings.Index(line, "^")
		if x >= 0 {
			if guardFound {
				panic("Multiple guards found")
			}
			guardFound = true
			guard = Guard{
				y:   y,
				x:   x,
				dir: 0,
			}
		}
	}
	if !guardFound {
		panic("No guard found")
	}
	return grid, guard
}

func walkOut(grid [][]rune, guard Guard) map[Vec2]struct{} {
	var H = len(grid)
	var W = len(grid[0])
	path := map[Vec2]struct{}{}
	for {
		path[Vec2{guard.y, guard.x}] = struct{}{}
		grid[guard.y][guard.x] = 'X'
		y := guard.y + Directions[guard.dir].y
		x := guard.x + Directions[guard.dir].x
		if y < 0 || H <= y || x < 0 || W <= x {
			return path
		}
		if grid[y][x] == '#' {
			guard.dir = (guard.dir + 3) % 4 // Turn right
		} else {
			guard.y = y
			guard.x = x
		}
	}
}

func printGrid(grid [][]rune) {
	if Print {
		for _, line := range grid {
			fmt.Println(string(line))
		}
	}
}
func part1(grid [][]rune, guard Guard) {
	timeStart := time.Now()
	walkOut(grid, guard)
	var count int
	for _, line := range grid {
		for _, cell := range line {
			if cell == 'X' {
				count++
			}
		}
	}
	printGrid(grid)
	fmt.Printf("Part 1: \t\t%d\tin %v\n", count, time.Since(timeStart))
}

func part2(grid [][]rune, guard Guard) {
	timeStart := time.Now()
	path := walkOut(grid, guard)
	var loopCount int
	for p := range path {
		y, x := p.y, p.x
		grid[y][x] = '#'
		if hasLoop(grid, guard) {
			loopCount++
			grid[y][x] = 'O'
		} else {
			grid[y][x] = ' '
		}
	}
	printGrid(grid)
	fmt.Printf("Part 2: \t\t%d\tin %v\n", loopCount, time.Since(timeStart))
}

var visited = map[Guard]struct{}{} // to not allocate on each call

func hasLoop(grid [][]rune, guard Guard) bool {
	clear(visited)
	var H = len(grid)
	var W = len(grid[0])
	for {
		if _, ok := visited[guard]; ok {
			return true
		}
		visited[guard] = struct{}{}
		y := guard.y + Directions[guard.dir].y
		x := guard.x + Directions[guard.dir].x
		if y < 0 || H <= y || x < 0 || W <= x {
			return false
		}
		if grid[y][x] == '#' {
			guard.dir = (guard.dir + 3) % 4 // Turn right
		} else {
			guard.y = y
			guard.x = x
		}
	}
}
