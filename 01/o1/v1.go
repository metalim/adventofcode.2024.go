// Напиши код на Go для решения следующей задачи. Входные данные в файле указываемом аргументом в командной строке.
// Выведи время решения каждой части.

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <inputfile>")
		return
	}
	startTotal := time.Now()

	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var leftList, rightList []int

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			continue
		}
		l, err := strconv.Atoi(parts[0])
		if err != nil {
			panic(err)
		}
		r, err := strconv.Atoi(parts[1])
		if err != nil {
			panic(err)
		}
		leftList = append(leftList, l)
		rightList = append(rightList, r)
	}

	// Part One
	part1Start := time.Now()
	sort.Ints(leftList)
	sort.Ints(rightList)
	var distSum int
	for i := range leftList {
		diff := leftList[i] - rightList[i]
		if diff < 0 {
			diff = -diff
		}
		distSum += diff
	}
	part1Time := time.Since(part1Start)

	// Part Two
	part2Start := time.Now()
	countRight := make(map[int]int)
	for _, v := range rightList {
		countRight[v]++
	}
	var simScore int
	for _, v := range leftList {
		simScore += v * countRight[v]
	}
	part2Time := time.Since(part2Start)

	fmt.Println("Part 1:", distSum)
	fmt.Println("Time Part 1:", part1Time)
	fmt.Println("Part 2:", simScore)
	fmt.Println("Time Part 2:", part2Time)
	fmt.Println("Total Time:", time.Since(startTotal))
}
