package main

import (
	"bytes"
	"flag"
	"fmt"
	"maps"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var Print = false

func catch(err error) {
	if err != nil {
		panic(err)
	}
}

type Guard struct {
	p   Vec2
	dir int
}
type Vec2 struct {
	y, x int
}

func (v Vec2) Add(o Vec2) Vec2 {
	return Vec2{v.y + o.y, v.x + o.x}
}

type Grid struct {
	Map map[Vec2]rune
	BR  Vec2
}

func (g *Grid) Clone() *Grid {
	return &Grid{
		Map: maps.Clone(g.Map),
		BR:  g.BR,
	}
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

	bs, err := os.ReadFile(flag.Arg(0))
	catch(err)

	grid, guard := parseInput(string(bs))
	part1(grid, guard)
	part2(grid, guard)
}

func parseInput(input string) (grid *Grid, guard Guard) {
	lines := strings.Split(input, "\n")
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}

	grid = &Grid{
		Map: make(map[Vec2]rune),
		BR:  Vec2{len(lines), len(lines[0])},
	}
	var guardFound bool
	for y, line := range lines {
		for x, cell := range line {
			if cell == '#' {
				grid.Map[Vec2{y, x}] = cell
			}
			if cell == '^' {
				if guardFound {
					panic("Multiple guards found")
				}
				guardFound = true
				guard = Guard{
					p:   Vec2{y, x},
					dir: 0,
				}
			}
		}
	}
	if !guardFound {
		panic("No guard found")
	}
	return grid, guard
}

func walkOut(grid *Grid, guard Guard) map[Vec2]struct{} {
	path := map[Vec2]struct{}{}
	for {
		path[guard.p] = struct{}{}
		grid.Map[guard.p] = 'X'
		np := guard.p.Add(Directions[guard.dir])
		if np.y < 0 || grid.BR.y <= np.y || np.x < 0 || grid.BR.x <= np.x {
			return path
		}
		if grid.Map[np] == '#' {
			guard.dir = (guard.dir + 3) % 4 // Turn right
		} else {
			guard.p = np
		}
	}
}

func printGrid(grid *Grid) {
	if !Print {
		return
	}
	fmt.Println(grid.BR)
	var buf bytes.Buffer
	for y := 0; y < grid.BR.y; y++ {
		for x := 0; x < grid.BR.x; x++ {
			r, ok := grid.Map[Vec2{y, x}]
			if !ok {
				r = ' '
			}
			buf.WriteRune(r)
		}
		buf.WriteRune('\n')
	}
	os.Stdout.Write(buf.Bytes())
}

func part1(grid *Grid, guard Guard) {
	timeStart := time.Now()
	path := walkOut(grid, guard)
	printGrid(grid)
	fmt.Printf("Part 1: \t\t%d\tin %v\n", len(path), time.Since(timeStart))
}

var Workers = runtime.NumCPU()

func part2(grid *Grid, guard Guard) {
	timeStart := time.Now()
	path := walkOut(grid, guard)
	fmt.Printf("Using %d workers\n", Workers)
	var loopCount atomic.Int64
	wg := sync.WaitGroup{}
	wg.Add(Workers)
	ch := make(chan Vec2, Workers)
	for i := 0; i < Workers; i++ {
		dupe := grid.Clone()

		go func(grid *Grid) {
			defer wg.Done()
			visited := map[Guard]struct{}{}
			for p := range ch {
				grid.Map[p] = '#'
				if hasLoop(grid, guard, visited) {
					loopCount.Add(1)
					grid.Map[p] = 'O'
				} else {
					grid.Map[p] = ' '
				}
			}
		}(dupe)
	}
	for p := range path {
		ch <- p
	}
	close(ch)
	wg.Wait()
	printGrid(grid)
	fmt.Printf("Part 2: \t\t%d\tin %v\n", loopCount.Load(), time.Since(timeStart))
}

func hasLoop(grid *Grid, guard Guard, visited map[Guard]struct{}) bool {
	clear(visited)
	for {
		if _, ok := visited[guard]; ok {
			return true
		}
		visited[guard] = struct{}{}
		np := guard.p.Add(Directions[guard.dir])
		if np.y < 0 || grid.BR.y <= np.y || np.x < 0 || grid.BR.x <= np.x {
			return false
		}
		if grid.Map[np] == '#' {
			guard.dir = (guard.dir + 3) % 4 // Turn right
		} else {
			guard.p = np
		}
	}
}
