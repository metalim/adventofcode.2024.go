package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
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

type Point [2]int // y, x

func (p Point) Add(d Point) Point {
	return Point{p[0] + d[0], p[1] + d[1]}
}

var directions = []Point{{-1, 0}, {0, -1}, {1, 0}, {0, 1}}

func bfs(input Input, start Point, symbol byte) map[Point]struct{} {
	visited := map[Point]struct{}{start: {}}
	queue := map[Point]struct{}{start: {}}
	found := map[Point]struct{}{start: {}}
	for len(queue) > 0 {
		next := map[Point]struct{}{}
		for p := range queue {
			for _, d := range directions {
				np := p.Add(d)
				if np[0] < 0 || np[0] >= len(input) || np[1] < 0 || np[1] >= len(input[0]) {
					continue
				}
				if _, ok := visited[np]; ok {
					continue
				}
				if input[np[0]][np[1]] != symbol {
					continue
				}
				visited[np] = struct{}{}
				found[np] = struct{}{}
				next[np] = struct{}{}
			}
		}
		queue = next
	}
	return found
}

type Plot struct {
	symbol byte
	points map[Point]struct{}
}

func findPlots(input Input) map[Point]*Plot {
	plots := map[Point]*Plot{}
	allFound := map[Point]struct{}{}
	for y, line := range input {
		for x, r := range line {
			p := Point{y, x}
			if _, ok := allFound[p]; ok {
				continue
			}
			found := bfs(input, p, byte(r))
			plots[p] = &Plot{symbol: byte(r), points: found}
			for p := range found {
				allFound[p] = struct{}{}
			}
		}
	}
	return plots
}

func part1(input Input) {
	timeStart := time.Now()

	var cost int
	plots := findPlots(input)

	for _, plot := range plots {
		area := len(plot.points)
		var perimeter int
		for p := range plot.points {
			for _, d := range directions {
				np := p.Add(d)
				if _, ok := plot.points[np]; !ok {
					perimeter++
				}
			}
		}
		cost += area * perimeter
	}

	fmt.Printf("Part 1: %d\t\tin %v\n", cost, time.Since(timeStart))
}

func part2(input Input) {
	timeStart := time.Now()

	plots := findPlots(input)

	var cost int
	for _, plot := range plots {
		area := len(plot.points)
		// line of walls is just one side
		var sides int
		for p := range plot.points {
			for i, d := range directions {
				np := p.Add(d)
				if _, ok := plot.points[np]; ok {
					continue
				}
				sideDir := directions[(i+1)%4]
				sideP := p.Add(sideDir)
				if _, ok := plot.points[sideP]; ok {
					sideNP := sideP.Add(d)
					if _, ok := plot.points[sideNP]; !ok {
						continue
					}
				}
				sides++
			}
		}

		cost += area * sides
	}

	fmt.Printf("Part 2: %d\t\tin %v\n", cost, time.Since(timeStart))
}
