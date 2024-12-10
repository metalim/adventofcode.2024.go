// Напиши код на Go для решения следующей задачи. Входные данные в файле указываемом аргументом в командной строке.
// Выведи время решения каждой части.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func isSafe(report []int) bool {
	if len(report) < 2 {
		return true
	}
	increasing := report[1] > report[0]
	for i := 1; i < len(report); i++ {
		diff := report[i] - report[i-1]
		if diff == 0 || diff < -3 || diff > 3 {
			return false
		}
		if (report[i] > report[i-1]) != increasing {
			return false
		}
	}
	return true
}

func canBeSafeWithOneRemoval(report []int) bool {
	for i := range report {
		newRep := append([]int{}, report[:i]...)
		newRep = append(newRep, report[i+1:]...)
		if isSafe(newRep) {
			return true
		}
	}
	return false
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <inputfile>")
		return
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	var reports [][]int
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		var rep []int
		for _, p := range parts {
			n, _ := strconv.Atoi(p)
			rep = append(rep, n)
		}
		reports = append(reports, rep)
	}

	start1 := time.Now()
	part1count := 0
	for _, r := range reports {
		if isSafe(r) {
			part1count++
		}
	}
	end1 := time.Now()

	start2 := time.Now()
	part2count := 0
	for _, r := range reports {
		if isSafe(r) || canBeSafeWithOneRemoval(r) {
			part2count++
		}
	}
	end2 := time.Now()

	fmt.Printf("Part 1: %d\nTime: %v\n", part1count, end1.Sub(start1))
	fmt.Printf("Part 2: %d\nTime: %v\n", part2count, end2.Sub(start2))
}
