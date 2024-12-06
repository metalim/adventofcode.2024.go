package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const PRINT_GRID = false

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
	if len(os.Args) != 2 {
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

	var found bool
	for y, line := range lines {
		grid = append(grid, []rune(line))
		x := strings.Index(line, "^")
		if x >= 0 {
			if found {
				panic("Multiple guards found")
			}
			found = true
			guard = Guard{
				y:   y,
				x:   x,
				dir: 0,
			}
		}
	}
	if !found {
		panic("No guard found")
	}
	return grid, guard
}

func walkOut(grid [][]rune, guard Guard) {
	var H = len(grid)
	var W = len(grid[0])
	for {
		grid[guard.y][guard.x] = 'X'
		y := guard.y + Directions[guard.dir].y
		x := guard.x + Directions[guard.dir].x
		if y < 0 || H <= y || x < 0 || W <= x {
			return
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
	if PRINT_GRID {
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
	walkOut(grid, guard)
	var loopCount int
	for y, line := range grid {
		for x, cell := range line {
			if cell != 'X' {
				continue
			}
			grid[y][x] = '#'
			if hasLoop(grid, guard) {
				loopCount++
				grid[y][x] = 'O'
			} else {
				grid[y][x] = ' '
			}
		}
	}
	printGrid(grid)
	fmt.Printf("Part 2: \t\t%d\tin %v\n", loopCount, time.Since(timeStart))
}

func hasLoop(grid [][]rune, guard Guard) bool {
	var H = len(grid)
	var W = len(grid[0])
	visited := map[Guard]struct{}{}
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
