package main

import (
	"fmt"
	"os"
	"regexp"
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
	lines := strings.Split(string(bs), "\n")
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}

	part1(lines)
	part2(lines)
}

func part1(lines []string) {
	timeStart := time.Now()
	var xmas int

	xmas += countXMAS(lines)
	xmas += countXMAS(diagonal(lines))

	for i := 0; i < 3; i++ {
		lines = rotate90(lines)
		xmas += countXMAS(lines)
		xmas += countXMAS(diagonal(lines))
	}
	fmt.Printf("Part 1: %d\t\tin %v\n", xmas, time.Since(timeStart))
}

var reXMAS = regexp.MustCompile(`XMAS`)

func countXMAS(lines []string) int {
	var xmas int
	for _, line := range lines {
		xmas += len(reXMAS.FindAllString(line, -1))
	}
	return xmas
}

func rotate90(lines []string) []string {
	var newLines []string
	for _, l := range lines {
		for i, c := range l {
			if i >= len(newLines) {
				newLines = append(newLines, "")
			}
			newLines[i] = string(c) + newLines[i]
		}
	}
	return newLines
}

func diagonal(lines []string) []string {
	var newLines []string
	for i, l := range lines {
		for j, c := range l {
			if i+j >= len(newLines) {
				newLines = append(newLines, "")
			}
			newLines[i+j] += string(c)
		}
	}
	return newLines
}

func part2(lines []string) {
	timeStart := time.Now()
	var xmas int
	for y, l := range lines[1 : len(lines)-1] {
		for x, c := range l[1 : len(l)-1] {
			if c != 'A' {
				continue
			}
			if !(lines[y][x] == 'M' && lines[y+2][x+2] == 'S' || lines[y][x] == 'S' && lines[y+2][x+2] == 'M') {
				continue
			}
			if !(lines[y][x+2] == 'M' && lines[y+2][x] == 'S' || lines[y][x+2] == 'S' && lines[y+2][x] == 'M') {
				continue
			}
			xmas++
		}
	}

	fmt.Printf("Part 2: %d\t\tin %v\n", xmas, time.Since(timeStart))
}
