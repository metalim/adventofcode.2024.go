package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"regexp"
	"slices"
	"strconv"
	"time"
)

const (
	Acost = 3
	Bcost = 1

	Part1MaxPresses = 100
	Part2Add        = 1e13
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
	part1(slices.Clone(input))
	part2(input)
}

type Input []Machine

type Machine struct {
	A     Point
	B     Point
	Prize Point
}

type Point struct {
	X int
	Y int
}

func (p Point) AddInt(v int) Point {
	return Point{p.X + v, p.Y + v}
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	catch(err)
	return i
}

var reMachine = regexp.MustCompile(`Button A: X\+(\d+), Y\+(\d+)\nButton B: X\+(\d+), Y\+(\d+)\nPrize: X=(\d+), Y=(\d+)`)

func parseInput(input string) Input {
	ms := reMachine.FindAllStringSubmatch(input, -1)
	machines := make(Input, 0, len(ms))
	for _, m := range ms {
		buttonA := Point{atoi(m[1]), atoi(m[2])}
		buttonB := Point{atoi(m[3]), atoi(m[4])}
		prize := Point{atoi(m[5]), atoi(m[6])}
		machines = append(machines, Machine{buttonA, buttonB, prize})
	}
	return machines
}

func part1(machines Input) {
	timeStart := time.Now()
	var totalMinCost int
	for _, m := range machines {
		minCost := math.MaxInt
		for a := 0; a < Part1MaxPresses; a++ {
			for b := 0; b < Part1MaxPresses; b++ {
				if m.A.X*a+m.B.X*b == m.Prize.X && m.A.Y*a+m.B.Y*b == m.Prize.Y {
					cost := Acost*a + Bcost*b
					if cost < minCost {
						minCost = cost
					}
				}
			}
		}
		if minCost == math.MaxInt {
			continue
		}
		totalMinCost += minCost
	}

	fmt.Printf("Part 1: %d\t\tin %v\n", totalMinCost, time.Since(timeStart))
}

func solve(p, a, b Point) int {
	/*
		solve for i and j
		i*a.X + j*b.X = p.X
		i*a.Y + j*b.Y = p.Y

		i = (p.X - j*b.X) / a.X
		(p.X - j*b.X) * a.Y / a.X + j*b.Y = p.Y
		p.X*a.Y - j*a.Y*b.X + j*a.X*b.Y = p.Y*a.X
		j*(a.X*b.Y - a.Y*b.X) = p.Y*a.X - p.X*a.Y
		j = (p.Y*a.X - p.X*a.Y) / (a.X*b.Y - a.Y*b.X)

		same way for i:
		i = (p.Y*b.X - p.X*b.Y) / (a.X*b.Y - a.Y*b.X)

		or multiply by -1 numerator and denominator to get same denominator (determinant):
		i = (p.X*b.Y - p.Y*b.X) / (a.X*b.Y - a.Y*b.X)
	*/

	det := a.X*b.Y - a.Y*b.X
	if det == 0 {
		return 0
	}
	detA := p.X*b.Y - p.Y*b.X
	detB := p.Y*a.X - p.X*a.Y
	if detA%det != 0 || detB%det != 0 {
		return 0
	}
	i := detA / det
	j := detB / det

	if i < 0 || j < 0 {
		return 0
	}
	return Acost*i + Bcost*j
}

func part2(machines Input) {
	timeStart := time.Now()
	var totalCost int
	for _, m := range machines {
		cost := solve(m.Prize.AddInt(Part2Add), m.A, m.B)
		totalCost += cost
	}

	fmt.Printf("Part 2: %d\t\tin %v\n", totalCost, time.Since(timeStart))
}
