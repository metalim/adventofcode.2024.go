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

var opsPart1 = []string{"+", "*"}
var opsPart2 = []string{"+", "*", "||"}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <inputfile>")
		return
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var lines []string
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	// Part 1
	startPart1 := time.Now()
	sumPart1 := 0
	for _, line := range lines {
		testVal, nums := parseLine(line)
		if canMakeValue(nums, testVal, opsPart1) {
			sumPart1 += testVal
		}
	}
	fmt.Println("Part 1:", sumPart1, "Time:", time.Since(startPart1))

	// Part 2
	startPart2 := time.Now()
	sumPart2 := 0
	for _, line := range lines {
		testVal, nums := parseLine(line)
		if canMakeValue(nums, testVal, opsPart2) {
			sumPart2 += testVal
		}
	}
	fmt.Println("Part 2:", sumPart2, "Time:", time.Since(startPart2))
}

func parseLine(line string) (int, []int) {
	parts := strings.Split(line, ":")
	testVal, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
	right := strings.Fields(strings.TrimSpace(parts[1]))
	var nums []int
	for _, r := range right {
		n, _ := strconv.Atoi(r)
		nums = append(nums, n)
	}
	return testVal, nums
}

func canMakeValue(nums []int, target int, ops []string) bool {
	if len(nums) == 1 {
		return nums[0] == target
	}
	return tryOps(nums, target, ops, 0, nums[0])
}

func tryOps(nums []int, target int, ops []string, idx int, current int) bool {
	if idx == len(nums)-1 {
		return current == target
	}
	nextVal := nums[idx+1]
	for _, op := range ops {
		switch op {
		case "+":
			if tryOps(nums, target, ops, idx+1, current+nextVal) {
				return true
			}
		case "*":
			if tryOps(nums, target, ops, idx+1, current*nextVal) {
				return true
			}
		case "||":
			concatVal := concatNumbers(current, nextVal)
			// При конкатенации нужно пересобрать цепочку вычислений,
			// т.к. конкатенация происходит "на месте" текущего и следующего элемента.
			// Поэтому мы создадим новый срез чисел, в котором текущий элемент заменим на конкатенированный,
			// а следующий пропустим.
			newNums := make([]int, 0, len(nums)-1)
			newNums = append(newNums, nums[:idx]...)
			newNums = append(newNums, concatVal)
			newNums = append(newNums, nums[idx+2:]...)
			if canMakeValue(newNums, target, ops) {
				return true
			}
		}
	}
	return false
}

func concatNumbers(a, b int) int {
	as := strconv.Itoa(a)
	bs := strconv.Itoa(b)
	res, _ := strconv.Atoi(as + bs)
	return res
}
