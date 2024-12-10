/*
ты не понял. Конкатенация использует результат вычислений слева, а не написанное число.
Т.е. все операции производятся последовательно слева направо, без приоритетов
*/

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

	startPart1 := time.Now()
	sumPart1 := 0
	for _, line := range lines {
		testVal, nums := parseLine(line)
		if canMakeValue(nums, testVal, opsPart1) {
			sumPart1 += testVal
		}
	}
	fmt.Println("Part 1:", sumPart1, "Time:", time.Since(startPart1))

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
	return opsRec(nums, target, ops, 1, nums[0])
}

// opsRec: перебираем все варианты расстановки операторов слева направо
// current - текущее вычисленное значение
// idx - индекс следующего числа
func opsRec(nums []int, target int, operators []string, idx int, current int) bool {
	if idx == len(nums) {
		return current == target
	}
	nextVal := nums[idx]
	for _, op := range operators {
		var val int
		if op == "+" {
			val = current + nextVal
		} else if op == "*" {
			val = current * nextVal
		} else {
			// Конкатенация
			valStr := strconv.Itoa(current) + strconv.Itoa(nextVal)
			v, err := strconv.Atoi(valStr)
			if err != nil {
				continue
			}
			val = v
		}

		if opsRec(nums, target, operators, idx+1, val) {
			return true
		}
	}
	return false
}
