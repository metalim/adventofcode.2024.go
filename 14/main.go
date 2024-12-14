package main

import (
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

const Part1Moves = 100
const Part2Moves = 10000

var W, H = 101, 103

func catch(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.IntVar(&W, "w", 101, "width")
	flag.IntVar(&H, "h", 103, "height")
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("Usage: go run . input.txt")
		os.Exit(1)
	}

	bs, err := os.ReadFile(flag.Arg(0))
	catch(err)

	input := parseInput(string(bs))
	part1(input)
	input = parseInput(string(bs))
	part2(input)
}

type Robot struct {
	P Point
	V Point
}

func (r *Robot) Move(t int) {
	r.P.X = mod(r.P.X+r.V.X*t, W)
	r.P.Y = mod(r.P.Y+r.V.Y*t, H)
}

type Point struct {
	X int
	Y int
}

func (p Point) Add(v Point) Point {
	return Point{X: p.X + v.X, Y: p.Y + v.Y}
}

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	catch(err)
	return i
}

type Input []*Robot

var reRobot = regexp.MustCompile(`p=(-?\d+),(-?\d+) v=(-?\d+),(-?\d+)`)

func parseInput(input string) Input {
	lines := strings.Split(input, "\n")
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	robots := make(Input, 0, len(lines))
	for _, line := range lines {
		parts := reRobot.FindStringSubmatch(line)
		r := &Robot{
			P: Point{X: toInt(parts[1]), Y: toInt(parts[2])},
			V: Point{X: toInt(parts[3]), Y: toInt(parts[4])},
		}
		robots = append(robots, r)
	}
	return robots
}

func mod(a, b int) int {
	return (a%b + b) % b
}

func safetyFactor(robots Input) int {
	var safety [4]int
	for _, r := range robots {
		switch {
		case r.P.X < W/2 && r.P.Y < H/2:
			safety[0]++
		case r.P.X > W/2 && r.P.Y < H/2:
			safety[1]++
		case r.P.X < W/2 && r.P.Y > H/2:
			safety[2]++
		case r.P.X > W/2 && r.P.Y > H/2:
			safety[3]++
		default:
			// ignore robots that are not in any quadrant
		}
	}
	return safety[0] * safety[1] * safety[2] * safety[3]
}

func part1(robots Input) {
	timeStart := time.Now()

	fmt.Printf("Moving %d robots %d times in a %dx%d grid\n", len(robots), Part1Moves, W, H)
	for _, r := range robots {
		r.Move(Part1Moves)
	}
	safetyFactor := safetyFactor(robots)
	fmt.Printf("Part 1: %d\t\tin %v\n", safetyFactor, time.Since(timeStart))
}

func bit(v bool) int {
	if v {
		return 1
	}
	return 0
}

var quarters = []rune(" ▘▝▀▖▌▞▛▗▚▐▜▄▙▟█")

var tree = color.New(color.FgGreen)
var box = color.New(color.FgRed)
var snow = color.New(color.FgWhite)

func fprintGridCompact(w io.Writer, robots Input) {
	grid := map[Point]bool{}
	for _, r := range robots {
		grid[r.P] = true
	}
	for y := 0; y < H; y += 2 {
		for x := 0; x < W; x += 2 {
			var bits int
			var nBits int
			for dy := 0; dy < 2; dy++ {
				for dx := 0; dx < 2; dx++ {
					p := Point{x + dx, y + dy}
					if _, ok := grid[p]; ok {
						bits |= 1 << (dy*2 + dx)
						nBits++
					}
				}
			}
			char := quarters[bits]
			o := snow
			switch nBits {
			case 2:
				o = box
			case 3, 4:
				o = tree
			}
			o.Fprintf(w, "%c", char)
		}
		w.Write([]byte("\n"))
	}
}

func part2(robots Input) {
	timeStart := time.Now()
	movingRobots := make(Input, len(robots))
	for i, r := range robots {
		movingRobots[i] = &Robot{P: r.P, V: r.V}
	}
	minMetric := math.MaxInt
	var minStep int
	for i := 1; i <= Part2Moves; i++ {
		var avg Point
		for _, r := range movingRobots {
			r.Move(1)
			avg = avg.Add(r.P)
		}
		avg = Point{X: avg.X / len(movingRobots), Y: avg.Y / len(movingRobots)}
		var asd int // average squared deviation
		for _, r := range movingRobots {
			asd += (r.P.X-avg.X)*(r.P.X-avg.X) + (r.P.Y-avg.Y)*(r.P.Y-avg.Y)
		}
		asd = asd / len(movingRobots)
		if asd < minMetric {
			minMetric = asd
			minStep = i
			fmt.Printf("New min metric: %d at step %d\n", minMetric, minStep)
		}
		if minMetric == 0 {
			break
		}
	}

	for _, r := range robots {
		r.Move(minStep)
	}
	fmt.Printf("Part 2: %d\t\tin %v\n", minStep, time.Since(timeStart))
	fprintGridCompact(os.Stdout, robots)
}
