/*
учти, что конкатенация тоже обрабатывается слева направо и не имеет приоритета перед другими операциями.
Результат второй части для примера выдаёт 4097, а должен выдавать 11387
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

// Для решения с учётом новой логики:
// 1. Генерируем все возможные варианты склейки чисел (конкатенация "||").
//    Это даёт нам различные последовательности чисел без оператора "||".
// 2. Для каждой последовательности проверяем все комбинации "+" и "*".
//    Операции вычисляем строго слева направо.
// 3. Если хоть в одном варианте результат совпал с искомым, уравнение считается выполнимым.

var opsAddMul = []string{"+", "*"} // Для части 1
var opsAll = []string{"+", "*"}    // Для части 2 после формирования групп без "||" мы используем только + и *
// Конкатенация будет обработана отдельно.

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

	// Часть 1: только + и *
	startPart1 := time.Now()
	sumPart1 := 0
	for _, line := range lines {
		testVal, nums := parseLine(line)
		if canMakeValuePart1(nums, testVal) {
			sumPart1 += testVal
		}
	}
	fmt.Println("Part 1:", sumPart1, "Time:", time.Since(startPart1))

	// Часть 2: теперь добавляем возможность конкатенации
	startPart2 := time.Now()
	sumPart2 := 0
	for _, line := range lines {
		testVal, nums := parseLine(line)
		if canMakeValuePart2(nums, testVal) {
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

// canMakeValuePart1 проверяет, можно ли получить testVal, вставляя только "+" или "*" между числами.
func canMakeValuePart1(nums []int, testVal int) bool {
	return tryAllOps(nums, testVal, opsAddMul)
}

// canMakeValuePart2 сначала генерирует все варианты группирования чисел с помощью "||", а затем для каждого варианта
// проверяет, можно ли получить testVal, вставляя "+" или "*".
func canMakeValuePart2(nums []int, testVal int) bool {
	groupings := generateConcatenations(nums, 0, nums[0], []int{})
	for _, g := range groupings {
		if tryAllOps(g, testVal, opsAddMul) {
			return true
		}
	}
	return false
}

// generateConcatenations генерирует все варианты числовых последовательностей, которые могут получиться,
// если между числами вставлять или не вставлять "||".
// idx - текущий индекс,
// currentNum - текущее число, которое мы набираем (может быть результатом конкатенации),
// currList - уже сформированная часть результата до currentNum.
func generateConcatenations(nums []int, idx int, currentNum int, currList []int) [][]int {
	if idx == len(nums)-1 {
		// Дошли до конца, добавляем текущий накопленный элемент
		result := make([]int, len(currList)+1)
		copy(result, currList)
		result[len(currList)] = currentNum
		return [][]int{result}
	}
	// Вариант 1: Не конкатенируем, добавляем currentNum в список и идём дальше
	noConcat := generateConcatenations(nums, idx+1, nums[idx+1], append(currList, currentNum))

	// Вариант 2: Конкатенируем текущий currentNum с nums[idx+1]
	concatVal, _ := strconv.Atoi(strconv.Itoa(currentNum) + strconv.Itoa(nums[idx+1]))
	concat := generateConcatenations(nums, idx+1, concatVal, currList)

	return append(noConcat, concat...)
}

// tryAllOps перебирает все возможные комбинации "+" и "*" между числами.
func tryAllOps(nums []int, target int, operators []string) bool {
	// Если только одно число
	if len(nums) == 1 {
		return nums[0] == target
	}
	return opsRec(nums, target, 1, nums[0], operators)
}

// opsRec рекурсивно перебирает операторы между числами и вычисляет слева направо
func opsRec(nums []int, target int, idx int, current int, ops []string) bool {
	if idx == len(nums) {
		return current == target
	}
	nextVal := nums[idx]
	for _, op := range ops {
		var val int
		if op == "+" {
			val = current + nextVal
		} else {
			val = current * nextVal
		}
		if opsRec(nums, target, idx+1, val, ops) {
			return true
		}
	}
	return false
}
