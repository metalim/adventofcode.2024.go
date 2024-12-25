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

	parsed := parseInput(string(bs))
	part1(parsed)
}

type Parsed [][]string

func parseInput(input string) Parsed {
	lines := strings.Split(input, "\n")
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}

	// don't parse numbers, that'll be the task itself lol
	// (I did parse in previous years, and found out parsing is the main part of the task)
	var parsed [][]string
	for i := 0; i < len(lines); i += 8 {
		parsed = append(parsed, lines[i:i+8])
	}
	return parsed
}

func part1(parsed Parsed) {
	timeStart := time.Now()
	var keys, locks [][5]int
	for _, grid := range parsed {
		if grid[0] == "#####" {
			var lock [5]int
			for x := 0; x < 5; x++ {
				for y := 1; y <= 5; y++ {
					if grid[y][x] == '#' {
						lock[x] = y
					}
				}
			}
			locks = append(locks, lock)
		} else {
			var key [5]int
			for x := 0; x < 5; x++ {
				for y := 5; y >= 1; y-- {
					if grid[y][x] == '#' {
						key[x] = 6 - y
					}
				}
			}
			keys = append(keys, key)
		}
	}

	var fitCount int
	for _, key := range keys {
		for _, lock := range locks {
			fit := true
			for x := 0; x < 5; x++ {
				if key[x]+lock[x] > 5 {
					fit = false
					break
				}
			}
			if fit {
				fitCount++
			}
		}
	}
	fmt.Printf("Part 1: %d\t\tin %v\n", fitCount, time.Since(timeStart))
}
