package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
)

const (
	LengthInput  = 1024
	LengthSample = 12
)

var LengthPart1 = LengthInput
var PrintGrid = false

func catch(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.IntVar(&LengthPart1, "length", LengthInput, "Length of the input for part 1")
	flag.BoolVar(&PrintGrid, "print", false, "Print the grid")
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("Usage: go run . input.txt")
		os.Exit(1)
	}

	bs, err := os.ReadFile(flag.Arg(0))
	catch(err)

	parsed := parseInput(string(bs))
	part1(parsed)
	part2(parsed)
}

type Point struct {
	X, Y int
}

func (p Point) Add(d Point) Point {
	return Point{p.X + d.X, p.Y + d.Y}
}

type Parsed struct {
	Points []Point
	BR     Point
}

func parseInput(input string) Parsed {
	lines := strings.Split(input, "\n")
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	parsed := Parsed{Points: make([]Point, len(lines))}
	for i, line := range lines {
		parsed.Points[i] = Point{}
		_, err := fmt.Sscanf(line, "%d,%d", &parsed.Points[i].X, &parsed.Points[i].Y)
		catch(err)
		if parsed.BR.X < parsed.Points[i].X {
			parsed.BR.X = parsed.Points[i].X
		}
		if parsed.BR.Y < parsed.Points[i].Y {
			parsed.BR.Y = parsed.Points[i].Y
		}
	}
	return parsed
}

var Dirs = []Point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

func bfs(grid Grid, start, end Point) int {
	visited := make(map[Point]bool)
	visited[start] = true
	steps := 0
	next := []Point{start}
	current := []Point{}

	for ; len(next) > 0; steps++ {
		current, next = next, current[:0]
		for _, p := range current {
			if p == end {
				return steps
			}
			for _, d := range Dirs {
				np := p.Add(d)
				if visited[np] {
					continue
				}
				if np.X < 0 || np.X > grid.BR.X || np.Y < 0 || np.Y > grid.BR.Y {
					continue
				}
				if _, ok := grid.Grid[np]; ok {
					continue
				}
				visited[np] = true
				next = append(next, np)
			}
		}
	}
	return -1
}

type Grid struct {
	Grid map[Point]struct{}
	BR   Point
}

var red = color.New(color.FgHiRed)

func (g Grid) Print(ps ...Point) {
	if !PrintGrid {
		return
	}
	fmt.Println("Grid:")
	for y := 0; y <= g.BR.Y; y++ {
	NextX:
		for x := 0; x <= g.BR.X; x++ {
			for _, p := range ps {
				if p.X == x && p.Y == y {
					red.Print("#")
					continue NextX
				}
			}

			if _, ok := g.Grid[Point{x, y}]; ok {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func NewGrid(parsed Parsed, length int) Grid {
	grid := Grid{Grid: make(map[Point]struct{}), BR: parsed.BR}
	for i, p := range parsed.Points {
		if i == length {
			break
		}
		grid.Grid[p] = struct{}{}
	}
	return grid
}

var Start = Point{0, 0}

func part1(parsed Parsed) {
	timeStart := time.Now()
	grid := NewGrid(parsed, LengthPart1)
	grid.Print()
	steps := bfs(grid, Start, grid.BR)
	fmt.Printf("Part 1: %d\t\tin %v\n", steps, time.Since(timeStart))
}

func part2(parsed Parsed) {
	part2_cut(parsed)
	part2_binary_search(parsed)
}

func part2_binary_search(parsed Parsed) {
	timeStart := time.Now()
	step := sort.Search(len(parsed.Points), func(i int) bool {
		grid := NewGrid(parsed, i+1)
		steps := bfs(grid, Start, grid.BR)
		return steps == -1
	})
	if step == len(parsed.Points) {
		fmt.Printf("Part 2: No solution found\t\tin %v\n", time.Since(timeStart))
		return
	}
	grid := NewGrid(parsed, step)
	grid.Print(parsed.Points[step])
	fmt.Printf("Part 2: %d, @(%d,%d)\t\tin %v\n", step, parsed.Points[step].X, parsed.Points[step].Y, time.Since(timeStart))

}

var neighborsWithDiagonal = []Point{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}

type CutSet struct {
	Points map[Point]struct{}
	TR, BL bool
}

func findJoin(parsed Parsed) (int, Point) {
	sets := map[Point]*CutSet{}

	for steps, p := range parsed.Points {
		var pSet *CutSet
		for _, dp := range neighborsWithDiagonal {
			np := p.Add(dp)
			if npSet, ok := sets[np]; ok {
				if pSet == nil {
					pSet = npSet
				} else if pSet == npSet {
					continue
				} else {
					pSet.TR = pSet.TR || npSet.TR
					pSet.BL = pSet.BL || npSet.BL
					if pSet.TR && pSet.BL {
						return steps, p
					}
					// merge sets
					for o := range npSet.Points {
						pSet.Points[o] = struct{}{}
						sets[o] = pSet
					}
				}
			}
		}
		if pSet == nil {
			pSet = &CutSet{Points: make(map[Point]struct{})}
		}
		pSet.Points[p] = struct{}{}
		if p.X == 0 || p.Y == parsed.BR.Y {
			pSet.BL = true
		}
		if p.X == parsed.BR.X || p.Y == 0 {
			pSet.TR = true
		}
		if pSet.BL && pSet.TR {
			return steps, p
		}
		sets[p] = pSet
	}
	return -1, Point{-1, -1}
}

func part2_cut(parsed Parsed) {
	timeStart := time.Now()
	steps, p := findJoin(parsed)
	fmt.Printf("Part 2: %d, @(%d,%d) \t\tin %v\n", steps, p.X, p.Y, time.Since(timeStart))
}
