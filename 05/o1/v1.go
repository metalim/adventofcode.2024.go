package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Проверяет корректность порядка страниц для данного обновления
func isCorrectOrder(update []int, rules map[int]map[int]bool) bool {
	pos := make(map[int]int)
	for i, p := range update {
		pos[p] = i
	}
	for x, deps := range rules {
		if _, ok := pos[x]; !ok {
			continue
		}
		for y := range deps {
			if _, ok := pos[y]; !ok {
				continue
			}
			if pos[x] > pos[y] {
				return false
			}
		}
	}
	return true
}

// Сортировка страниц обновления в правильном порядке
func sortUpdate(update []int, rules map[int]map[int]bool) []int {
	// Строим граф зависимостей
	inDegree := make(map[int]int)
	graph := make(map[int][]int)
	pages := make(map[int]bool)
	for _, p := range update {
		pages[p] = true
	}
	for x, deps := range rules {
		if !pages[x] {
			continue
		}
		for y := range deps {
			if !pages[y] {
				continue
			}
			graph[x] = append(graph[x], y)
		}
	}
	// Подсчет входных степеней
	for x := range graph {
		for _, y := range graph[x] {
			inDegree[y]++
		}
	}
	// Топологическая сортировка
	var q []int
	for p := range pages {
		if inDegree[p] == 0 {
			q = append(q, p)
		}
	}
	var result []int
	for len(q) > 0 {
		x := q[0]
		q = q[1:]
		result = append(result, x)
		for _, y := range graph[x] {
			inDegree[y]--
			if inDegree[y] == 0 {
				q = append(q, y)
			}
		}
	}
	return result
}

// Получение средней страницы (для нечетного количества элементов)
func middlePage(update []int) int {
	return update[len(update)/2]
}

func main() {
	start := time.Now()

	if len(os.Args) < 2 {
		fmt.Println("Укажите путь к файлу")
		return
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	// Чтение правил
	rules := make(map[int]map[int]bool)
	var line string
	for sc.Scan() {
		line = strings.TrimSpace(sc.Text())
		if line == "" {
			break
		}
		parts := strings.Split(line, "|")
		if len(parts) != 2 {
			continue
		}
		x, _ := strconv.Atoi(parts[0])
		y, _ := strconv.Atoi(parts[1])
		if rules[x] == nil {
			rules[x] = make(map[int]bool)
		}
		rules[x][y] = true
	}

	part1Start := time.Now()

	var correctSum, incorrectSum int
	var incorrectUpdates [][]int

	// Чтение обновлений
	for sc.Scan() {
		line = strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		strPages := strings.Split(line, ",")
		var update []int
		for _, sp := range strPages {
			p, _ := strconv.Atoi(strings.TrimSpace(sp))
			update = append(update, p)
		}

		if isCorrectOrder(update, rules) {
			correctSum += middlePage(update)
		} else {
			incorrectUpdates = append(incorrectUpdates, update)
		}
	}

	fmt.Println("Part1:", correctSum)
	fmt.Println("Time Part1:", time.Since(part1Start))

	part2Start := time.Now()

	// Обработка некорректных обновлений
	for _, up := range incorrectUpdates {
		sortedUp := sortUpdate(up, rules)
		incorrectSum += middlePage(sortedUp)
	}

	fmt.Println("Part2:", incorrectSum)
	fmt.Println("Time Part2:", time.Since(part2Start))

	fmt.Println("Total Time:", time.Since(start))
}
