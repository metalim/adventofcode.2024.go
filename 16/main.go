package main

import (
	"flag"
	"fmt"
	"math"
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

type Vec2 struct {
	x, y int
}

func (v Vec2) Add(other Vec2) Vec2 {
	return Vec2{v.x + other.x, v.y + other.y}
}

func (v Vec2) Sub(other Vec2) Vec2 {
	return Vec2{v.x - other.x, v.y - other.y}
}

type Dir int

const (
	Right Dir = iota
	Up
	Left
	Down
)

var Directions = []Vec2{
	Right: {1, 0},
	Up:    {0, -1},
	Left:  {-1, 0},
	Down:  {0, 1},
}

type Parsed struct {
	Map  []string
	W, H int

	Start, End Vec2
}

func parseInput(input string) Parsed {
	lines := strings.Split(input, "\n")
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}

	parsed := Parsed{
		Map: lines,
		W:   len(lines[0]),
		H:   len(lines),
	}

	for y, line := range lines {
		for x, c := range line {
			if c == 'S' {
				parsed.Start = Vec2{x, y}
			}
			if c == 'E' {
				parsed.End = Vec2{x, y}
			}
		}
	}
	return parsed
}

type Deer struct {
	Pos Vec2
	Dir Dir
}

const Wall = '#'

func bfs(parsed Parsed, start Deer) map[Deer]int {
	minScores := map[Deer]int{start: 0}
	next := map[Deer]int{start: 0}
	current := map[Deer]int{}
	for len(next) > 0 {
		current, next = next, current
		clear(next)
		for deer, score := range current {
			for dir, dirVec := range Directions {
				dir := Dir(dir)
				nDeer := Deer{Pos: deer.Pos.Add(dirVec), Dir: dir}
				if nDeer.Pos.x < 0 || nDeer.Pos.x >= parsed.W || nDeer.Pos.y < 0 || nDeer.Pos.y >= parsed.H {
					continue
				}

				if parsed.Map[nDeer.Pos.y][nDeer.Pos.x] == Wall {
					continue
				}

				nextScore := score + 1
				if dir != deer.Dir {
					nextScore += 1000
				}
				if minScore, ok := minScores[nDeer]; !ok || nextScore < minScore {
					minScores[nDeer] = nextScore
					next[nDeer] = nextScore
				}
			}
		}
	}
	return minScores
}

func getMinScore(minScores map[Deer]int, pos Vec2) (int, Dir) {
	minScore := math.MaxInt
	var minDir Dir
	for dir := range Directions {
		dir := Dir(dir)
		if score, ok := minScores[Deer{Pos: pos, Dir: dir}]; ok {
			if score < minScore {
				minScore = score
				minDir = dir
			}
		}
	}
	return minScore, minDir
}
func part1(parsed Parsed) {
	timeStart := time.Now()
	start := Deer{Pos: parsed.Start, Dir: 0}
	minScores := bfs(parsed, start)
	minScore, _ := getMinScore(minScores, parsed.End)
	fmt.Printf("Part 1: %d\t\tin %v\n", minScore, time.Since(timeStart))
}

func part2(parsed Parsed) {
	timeStart := time.Now()
	start := Deer{Pos: parsed.Start, Dir: 0}
	minScores := bfs(parsed, start)
	minScore, minDir := getMinScore(minScores, parsed.End)
	prev := map[Deer]int{{Pos: parsed.End, Dir: minDir}: minScore}
	current := map[Deer]int{}
	paths := map[Vec2]int{parsed.End: minScore}
	for len(prev) > 0 {
		current, prev = prev, current
		clear(prev)
		for deer, score := range current {
			for dir, dirVec := range Directions {
				dir := Dir(dir)
				pp := deer.Pos.Sub(dirVec)
				if pp.x < 0 || pp.x >= parsed.W || pp.y < 0 || pp.y >= parsed.H {
					continue
				}
				prevScore := score - 1
				for pDir := range Directions {
					pDir := Dir(pDir)
					pDeer := Deer{Pos: pp, Dir: pDir}
					pDeerScore := prevScore
					if pDir != dir {
						pDeerScore -= 1000
					}
					if score, ok := minScores[pDeer]; ok && score == pDeerScore {
						prev[pDeer] = pDeerScore
						paths[pDeer.Pos] = pDeerScore
					}
				}
			}
		}
	}
	fmt.Printf("Part 2: %d\t\tin %v\n", len(paths), time.Since(timeStart))
}

/*
########
#.....E#
#.####.#
#S.....#
########

#     #     #     #     #     #     #     #
#  1002^ 2003> 2004> 2005> 2006> 1007^    #
#  1001^    #     #     #     #  1006^    #
#     0>    1>    2>    3>    4>    5>    #
#     #     #     #     #     #     #     #

*/
