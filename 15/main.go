package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

const GPSY = 100
const Space = ' '

var Print = false
var Print1 = false
var Print2 = false
var Freq = 20
var Delay = 50 * time.Millisecond

func catch(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.BoolVar(&Print, "print", false, "print the grid")
	flag.BoolVar(&Print1, "print1", false, "print the grid for part 1")
	flag.BoolVar(&Print2, "print2", false, "print the grid for part 2")
	flag.IntVar(&Freq, "freq", 20, "frequency of the print")
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("Usage: go run . input.txt")
		os.Exit(1)
	}

	Delay = time.Second / time.Duration(Freq)
	bs, err := os.ReadFile(flag.Arg(0))
	catch(err)

	input := parseInput(string(bs))
	part1(input)
	part2(input)
}

type Input struct {
	Room         []string
	Instructions string
}

func parseInput(input string) Input {
	parts := strings.Split(input, "\n\n")
	return Input{
		Room:         strings.Split(parts[0], "\n"),
		Instructions: strings.ReplaceAll(parts[1], "\n", ""),
	}
}

type Point struct {
	X int
	Y int
}

func (p Point) Add(q Point) Point {
	return Point{X: p.X + q.X, Y: p.Y + q.Y}
}

var directions = map[rune]Point{
	'v': {0, 1},
	'^': {0, -1},
	'>': {1, 0},
	'<': {-1, 0},
}

func canMove(p Point, dir Point, room map[Point]rune) bool {
	np := p.Add(dir)
	switch room[np] {
	case '#':
		return false
	case 'O':
		return canMove(np, dir, room)
	case '[':
		if dir.X == 0 {
			return canMove(np, dir, room) && canMove(np.Add(Point{X: 1, Y: 0}), dir, room)
		}
		return canMove(np, dir, room)
	case ']':
		if dir.X == 0 {
			return canMove(np, dir, room) && canMove(np.Add(Point{X: -1, Y: 0}), dir, room)
		}
		return canMove(np, dir, room)
	default:
		return true
	}
}

func move(p Point, dir Point, room map[Point]rune) (Point, bool) {
	np := p.Add(dir)
	switch room[np] {
	case '#':
		return p, false
	case 'O':
		_, ok := move(np, dir, room)
		if !ok {
			return p, false
		}
	case '[':
		if dir.X == 0 {
			move(np.Add(Point{X: 1, Y: 0}), dir, room)
		}
		move(np, dir, room)
	case ']':
		if dir.X == 0 {
			move(np.Add(Point{X: -1, Y: 0}), dir, room)
		}
		move(np, dir, room)
	}
	room[p], room[np] = Space, room[p]
	return np, true
}

func part1(input Input) {
	if Print1 {
		saved := Print
		Print = true
		defer func() {
			Print = saved
		}()
	}
	timeStart := time.Now()
	room := map[Point]rune{}
	robot := Point{X: 0, Y: 0}
	H := len(input.Room)
	W := len(input.Room[0])
	for y, line := range input.Room {
		for x, c := range line {
			room[Point{X: x, Y: y}] = c
			if c == '@' {
				robot = Point{X: x, Y: y}
			}
		}
	}
	initPrint()
	printGrid(room, W, H, 0, input.Instructions)
	for i, instruction := range input.Instructions {
		robot, _ = move(robot, directions[instruction], room)
		printGrid(room, W, H, i, input.Instructions)
	}
	var sum int
	for pos, c := range room {
		if c == 'O' {
			sum += pos.X + pos.Y*GPSY
		}
	}

	fmt.Printf("Part 1: %d\t\tin %v\n", sum, time.Since(timeStart))
}

func part2(input Input) {
	if Print2 {
		saved := Print
		Print = true
		defer func() {
			Print = saved
		}()
	}

	timeStart := time.Now()
	room := map[Point]rune{}
	robot := Point{X: 0, Y: 0}
	H := len(input.Room)
	W := len(input.Room[0]) * 2
	for y, line := range input.Room {
		for x, c := range line {
			switch c {
			case 'O':
				room[Point{X: x * 2, Y: y}] = '['
				room[Point{X: x*2 + 1, Y: y}] = ']'
			case '@':
				room[Point{X: x * 2, Y: y}] = c
				room[Point{X: x*2 + 1, Y: y}] = Space
				robot = Point{X: x * 2, Y: y}
			default:
				room[Point{X: x * 2, Y: y}] = c
				room[Point{X: x*2 + 1, Y: y}] = c
			}
		}
	}
	initPrint()
	printGrid(room, W, H, 0, input.Instructions)
	for i, instruction := range input.Instructions {
		if canMove(robot, directions[instruction], room) {
			robot, _ = move(robot, directions[instruction], room)
		}

		printGrid(room, W, H, i, input.Instructions)
	}
	var sum int
	for pos, c := range room {
		if c == '[' {
			sum += pos.X + pos.Y*GPSY
		}
	}

	fmt.Printf("Part 2: %d\t\tin %v\n", sum, time.Since(timeStart))
}
