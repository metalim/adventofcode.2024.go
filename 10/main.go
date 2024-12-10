package main

import (
	"flag"
	"fmt"
	"maps"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

const PRINT_MAP = false

func catch(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("Usage: go run main.go input.txt")
		os.Exit(1)
	}

	bs, err := os.ReadFile(flag.Arg(0))
	catch(err)

	input := parseInput(string(bs))
	part1(input)
	part2(input)
}

type Input []string

func parseInput(input string) Input {
	lines := strings.Split(input, "\n")
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	return Input(lines)
}

type point [2]int // [y, x]
func (p point) Add(q point) point {
	return point{p[0] + q[0], p[1] + q[1]}
}

var directions = []point{{-1, 0}, {0, -1}, {1, 0}, {0, 1}}

func part1(input Input) {
	timeStart := time.Now()

	// usage: trails[np][peak]=struct{}{}
	trails := map[point]map[point]struct{}{}
	peaks := map[point]struct{}{}
	for y, line := range input {
		for x, c := range line {
			if c == '9' {
				p := point{y, x}
				peaks[p] = struct{}{}
				trails[p] = map[point]struct{}{}
				trails[p][p] = struct{}{}
			}
		}
	}
	next := maps.Clone(peaks)
	cur := map[point]struct{}{}
	H, W := len(input), len(input[0])
	for v := byte('8'); v >= byte('0'); v-- {
		printMap(input, next, true, "Searching for %c, from %d points", v, len(next))
		cur, next = next, cur
		clear(next) // reuse, lol
		for p := range cur {
			for _, d := range directions {
				np := p.Add(d)
				if np[0] < 0 || np[0] >= H || np[1] < 0 || np[1] >= W {
					continue
				}
				if input[np[0]][np[1]] != v {
					continue
				}
				next[np] = struct{}{}
				if _, ok := trails[np]; !ok {
					trails[np] = map[point]struct{}{}
				}
				for peak := range trails[p] {
					trails[np][peak] = struct{}{}
				}
			}
		}
	}
	printMap(input, next, false, "Final map:")

	trailheads := next
	var sum int
	for p := range trailheads {
		sum += len(trails[p])
	}
	fmt.Printf("Part 1: %d\t\tin %v\n", sum, time.Since(timeStart))
}

var cPoint = color.New(color.FgYellow).Add(color.Bold)
var cNeighbor = color.New(color.FgGreen)
var cFiller = color.New(color.FgBlack)
var cHead = color.New(color.FgRed)

func printMap(input Input, points map[point]struct{}, printNeighbors bool, format string, a ...interface{}) {
	if !PRINT_MAP {
		return
	}
	fmt.Printf("\n"+format+"\n", a...)
	for y, line := range input {
		for x, c := range line {
			p := point{y, x}
			if _, ok := points[p]; ok {
				cPoint.Printf("%c", c)
			} else {
				var found bool
				if printNeighbors {
					for _, d := range directions {
						if _, ok := points[p.Add(d)]; ok {
							cNeighbor.Printf("%c", c)
							found = true
							break
						}
					}
				}
				if !found {
					if c == '0' {
						cHead.Printf("%c", c)
					} else {
						cFiller.Printf("%c", c)
					}
				}
			}
		}
		fmt.Println()
	}
}

func part2(input Input) {
	timeStart := time.Now()
	for _, line := range input {
		_ = line
	}

	fmt.Printf("Part 2: \t\tin %v\n", time.Since(timeStart))
}
