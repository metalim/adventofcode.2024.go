package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

const StepsPart1 = 25
const StepsPart2 = 75

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

type Input []int

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	catch(err)
	return i
}

func parseInput(input string) Input {
	parts := strings.Split(strings.TrimSpace(input), " ")
	nums := make([]int, len(parts))
	for i, part := range parts {
		nums[i] = toInt(part)
	}
	return Input(nums)
}

func blink(stones Input, steps int) Input {
	next := make(Input, len(stones))
	for step := 0; step < steps; step++ {
		for i, stone := range stones {
			if stone == 0 {
				next[i] = 1
				continue
			}
			s := strconv.Itoa(stone)
			if len(s)%2 == 0 {
				next[i] = toInt(s[:len(s)/2])
				next = append(next, toInt(s[len(s)/2:]))
				continue
			}
			next[i] = stone * 2024
		}
		stones, next = next, stones
		next = slices.Grow(next, len(stones))[:len(stones)]
	}
	return stones
}

func part1(stones Input) {
	timeStart := time.Now()
	stones = blink(stones, StepsPart1)

	fmt.Printf("Part 1: %d\t\tin %v\n", len(stones), time.Since(timeStart))
}

func blinkStoneOnce(stone int) []int {
	if stone == 0 {
		return []int{1}
	}
	s := strconv.Itoa(stone)
	if len(s)%2 == 0 {
		return []int{toInt(s[:len(s)/2]), toInt(s[len(s)/2:])}
	}
	return []int{stone * 2024}
}

var mem = make(map[[2]int]int)

func blinkStoneCount(stone int, steps int) int {
	if v, ok := mem[[2]int{stone, steps}]; ok {
		return v
	}

	next := blinkStoneOnce(stone)
	mem[[2]int{stone, 1}] = len(next)
	if steps == 1 {
		return len(next)
	}

	var count int
	for _, v := range next {
		count += blinkStoneCount(v, steps-1)
	}
	mem[[2]int{stone, steps}] = count
	return count
}

func blinkStonesCount(stones Input, steps int) int {
	var count int
	for _, stone := range stones {
		count += blinkStoneCount(stone, steps)
	}
	return count
}

func part2(stones Input) {
	timeStart := time.Now()
	count := blinkStonesCount(stones, StepsPart2)
	fmt.Printf("Mem size: %d\n", len(mem))

	fmt.Printf("Part 2: %d\t\tin %v\n", count, time.Since(timeStart))
}
