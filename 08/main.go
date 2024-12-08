package main

import (
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
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go input.txt")
		os.Exit(1)
	}

	bs, err := os.ReadFile(os.Args[1])
	catch(err)

	input := parseInput(string(bs))
	part1(input)
	part2(input)
}

func parseInput(input string) []string {
	lines := strings.Split(input, "\n")
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	return lines
}

type Vec2 [2]int

func part1(lines []string) {
	timeStart := time.Now()
	nodes := map[rune][]Vec2{}
	for y, line := range lines {
		for x, c := range line {
			if c != '.' {
				nodes[c] = append(nodes[c], Vec2{y, x})
			}
		}
	}
	H := len(lines)
	W := len(lines[0])
	antinodes := map[Vec2]rune{}
	for k, ns := range nodes {
		for i, pos1 := range ns {
			for j, pos2 := range ns {
				if i == j {
					continue
				}
				y := 2*pos1[0] - pos2[0]
				x := 2*pos1[1] - pos2[1]
				if y < 0 || y >= H || x < 0 || x >= W {
					continue
				}
				antinodes[Vec2{y, x}] = k
			}
		}
	}

	fmt.Printf("Part 1: %d\t\tin %v\n", len(antinodes), time.Since(timeStart))
}

func part2(lines []string) {
	timeStart := time.Now()
	nodes := map[rune][]Vec2{}
	for y, line := range lines {
		for x, c := range line {
			if c != '.' {
				nodes[c] = append(nodes[c], Vec2{y, x})
			}
		}
	}
	H := len(lines)
	W := len(lines[0])
	antinodes := map[Vec2]rune{}
	for k, ns := range nodes {
		for i, pos1 := range ns {
			for j, pos2 := range ns {
				if i == j {
					continue
				}
				dy := pos1[0] - pos2[0]
				dx := pos1[1] - pos2[1]
				y := pos1[0]
				x := pos1[1]
				for 0 <= y && y < H && 0 <= x && x < W {
					antinodes[Vec2{y, x}] = k
					y += dy
					x += dx
				}
			}
		}
	}

	fmt.Printf("Part 2: %d\t\tin %v\n", len(antinodes), time.Since(timeStart))
}
