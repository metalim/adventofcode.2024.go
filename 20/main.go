package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	SaveAtLeast1 = 100
	CheatTime1   = 2

	SaveAtLeast2 = 100
	CheatTime2   = 20
)

func catch(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
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

func (p Point) Add(o Point) Point {
	return Point{X: p.X + o.X, Y: p.Y + o.Y}
}

var Dirs = []Point{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

type Parsed struct {
	Map   []string
	W, H  int
	Start Point
	End   Point
}

func parseInput(input string) *Parsed {
	lines := strings.Split(input, "\n")
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	parsed := &Parsed{Map: lines, W: len(lines[0]), H: len(lines)}
	for y, line := range lines {
		for x, c := range line {
			if c == 'S' {
				parsed.Start = Point{X: x, Y: y}
			}
			if c == 'E' {
				parsed.End = Point{X: x, Y: y}
			}
		}
	}
	return parsed
}

const Wall = '#'

type Distances map[Point]int

func (d Distances) Print(parsed *Parsed, label string) {
	return
	fmt.Printf("\n%s:\n", label)
	for y := 0; y < parsed.H; y++ {
		for x := 0; x < parsed.W; x++ {
			p := Point{X: x, Y: y}
			if steps, ok := d[p]; ok {
				fmt.Printf("%3d", steps)
			} else {
				fmt.Printf("  %c", parsed.Map[y][x])
			}
		}
		fmt.Println()
	}
}

func bfs(parsed *Parsed, start, end Point) (steps int, visited Distances) {
	visited = Distances{start: 0} // save steps to reach point
	next := []Point{start}
	var cur []Point

	var step int
	for len(next) > 0 {
		step++ // start with 1
		cur, next = next, cur[:0]
		for _, p := range cur {
			if p == end {
				return step - 1, visited
			}
			for _, dir := range Dirs {
				nb := p.Add(dir)
				if nb.X < 0 || nb.X >= parsed.W || nb.Y < 0 || nb.Y >= parsed.H {
					continue
				}
				if parsed.Map[nb.Y][nb.X] == Wall {
					continue
				}

				if _, ok := visited[nb]; ok {
					continue
				}
				visited[nb] = step
				next = append(next, nb)
			}
		}
	}
	panic("path not found")
}

// brain! work!
// there were DFS tasks last days, so brain is **buzzled** with mem + DFS
// (typo is nice)
// Do we really need DFS? It just keeps popping in my head
// Nice move, Eric! 2 relaxing days and alignment with the task
// I think same trick was used in previous years, not sure
//
// 50 minutes in, and still no vision
// Cursor keeps suggesting BULLSHIT, that also is distracting
// Fuck you, Cursor! You hear me?
//
// Ok, vision is following:
// 1. BFS forward from the start
// 2. BFS backward from the end
// 3. cheat and count
// easy!

func findWaysToCheat(parsed *Parsed, start, end Point, saveAtLeast, maxCheatTime int) (ways, stepsWithoutCheating int) {
	// 1. simple bfs first, fill forward map
	stepsForward, forward := bfs(parsed, start, end)
	forward.Print(parsed, "forward")

	// 2. bfs backwards
	stepsBackward, backward := bfs(parsed, end, start)
	if stepsForward != stepsBackward {
		panic("wut?")
	}
	backward.Print(parsed, "backward")

	// 3. iterate over all visited points and check if we can still cheat from there
	// 2 cheat steps = 1 wall, because we need to land on empty space again
	for p, cheatStep := range forward {
		// p = (1,3)
		if cheatStep > stepsForward-saveAtLeast-2 {
			continue
		}
		// now we need to skip up to maxCheatTime
		for dy := -maxCheatTime; dy <= maxCheatTime; dy++ {
			cheatTimeLeft := maxCheatTime - abs(dy)
			for dx := -cheatTimeLeft; dx <= cheatTimeLeft; dx++ {
				np := Point{X: p.X + dx, Y: p.Y + dy}
				if np.X < 0 || np.X >= parsed.W || np.Y < 0 || np.Y >= parsed.H {
					continue
				}
				if parsed.Map[np.Y][np.X] == Wall {
					continue
				}
				cheatTime := abs(dx) + abs(dy)
				if backStep, ok := backward[np]; ok {
					if cheatStep+cheatTime+backStep <= stepsForward-saveAtLeast {
						ways++
					}
				}
			}
		}
	}

	return ways, stepsForward
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func part1(parsed *Parsed) {
	timeStart := time.Now()
	ways, stepsWithoutCheating := findWaysToCheat(parsed, parsed.Start, parsed.End, SaveAtLeast1, CheatTime1)

	fmt.Printf("Part 1: %d to save %d steps out of %d, with cheat time %d\t\tin %v\n", ways, SaveAtLeast1, stepsWithoutCheating, CheatTime1, time.Since(timeStart))
}

func part2(parsed *Parsed) {
	timeStart := time.Now()
	ways, stepsWithoutCheating := findWaysToCheat(parsed, parsed.Start, parsed.End, SaveAtLeast2, CheatTime2)

	fmt.Printf("Part 2: %d to save %d steps out of %d, with cheat time %d\t\tin %v\n", ways, SaveAtLeast2, stepsWithoutCheating, CheatTime2, time.Since(timeStart))

}
