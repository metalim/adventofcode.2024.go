package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
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

	reMul := regexp.MustCompile(`mul\((\d+),(\d+)\)`)
	var total int
	for _, line := range lines {
		matches := reMul.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			a, _ := strconv.Atoi(match[1])
			b, _ := strconv.Atoi(match[2])
			total += a * b
		}
	}
	fmt.Printf("Part 1: %d\t\tin %v\n", total, time.Since(timeStart))
}

func part2(lines []string) {
	timeStart := time.Now()
	re := regexp.MustCompile(`mul\((\d+),(\d+)\)|do\(\)|don't\(\)`)
	var total int
	enabled := true
	for _, line := range lines {
		matches := re.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if match[0] == "do()" {
				enabled = true
				continue
			}
			if match[0] == "don't()" {
				enabled = false
				continue
			}
			if !enabled {
				continue
			}
			a, _ := strconv.Atoi(match[1])
			b, _ := strconv.Atoi(match[2])
			total += a * b
		}
	}

	fmt.Printf("Part 2: %d\t\tin %v\n", total, time.Since(timeStart))
}
