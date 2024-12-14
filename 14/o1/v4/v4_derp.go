/*
➜ go run ./o1/v3 input.txt
Часть Первая Ответ: 221655456
Время выполнения Части Первой: 860.791µs
Часть Вторая Ответ: 0 секунд
Время выполнения Части Второй: 10.75µs

ебать ты идиот. Я же тебе сказал: НЕ ВСЕ РОБОТЫ КУЧКУЮТСЯ. Т.е. часть роботов складывается в рисунок ёлки, а другая часть всё так же продолжает шататься по всему полю. Используй ДРУГУЮ эвристику.
*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

// Constants for grid size (used only for Part One)
const (
	Width  = 101
	Height = 103
)

// Robot представляет робота с позицией и скоростью
type Robot struct {
	x  int
	y  int
	dx int
	dy int
}

// parseLine парсит строку ввода и возвращает робота
func parseLine(line string) (Robot, error) {
	// Пример строки: p=0,4 v=3,-3
	re := regexp.MustCompile(`p=(-?\d+),(-?\d+) v=(-?\d+),(-?\d+)`)
	matches := re.FindStringSubmatch(line)
	if matches == nil || len(matches) != 5 {
		return Robot{}, fmt.Errorf("неверный формат строки: %s", line)
	}
	x, err := strconv.Atoi(matches[1])
	if err != nil {
		return Robot{}, err
	}
	y, err := strconv.Atoi(matches[2])
	if err != nil {
		return Robot{}, err
	}
	dx, err := strconv.Atoi(matches[3])
	if err != nil {
		return Robot{}, err
	}
	dy, err := strconv.Atoi(matches[4])
	if err != nil {
		return Robot{}, err
	}
	return Robot{x: x, y: y, dx: dx, dy: dy}, nil
}

// readInput читает входной файл и возвращает срез роботов
func readInput(filename string) ([]Robot, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	robots := []Robot{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		robot, err := parseLine(line)
		if err != nil {
			return nil, err
		}
		robots = append(robots, robot)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return robots, nil
}

// simulate симулирует движение роботов на n секунд.
// Если wrap=true, позиции оборачиваются по границам (используется для Part One).
// Если wrap=false, позиции продолжают двигаться в бесконечном пространстве (используется для Part Two).
func simulate(robots []Robot, n int, wrap bool) []Robot {
	simulated := make([]Robot, len(robots))
	copy(simulated, robots)
	for i := 0; i < n; i++ {
		for j := 0; j < len(simulated); j++ {
			simulated[j].x += simulated[j].dx
			simulated[j].y += simulated[j].dy
			if wrap {
				simulated[j].x = simulated[j].x % Width
				if simulated[j].x < 0 {
					simulated[j].x += Width
				}
				simulated[j].y = simulated[j].y % Height
				if simulated[j].y < 0 {
					simulated[j].y += Height
				}
			}
		}
	}
	return simulated
}

// calculateSafetyFactor вычисляет фактор безопасности после 100 секунд
func calculateSafetyFactor(robots []Robot) int {
	// Определение квадрантов
	midX := Width / 2
	midY := Height / 2

	q1, q2, q3, q4 := 0, 0, 0, 0
	for _, robot := range robots {
		// Игнорировать роботов, находящихся точно на середине
		if robot.x == midX || robot.y == midY {
			continue
		}
		if robot.x < midX && robot.y < midY {
			q1++
		} else if robot.x >= midX && robot.y < midY {
			q2++
		} else if robot.x < midX && robot.y >= midY {
			q3++
		} else if robot.x >= midX && robot.y >= midY {
			q4++
		}
	}
	return q1 * q2 * q3 * q4
}

// countClosePairs считает количество пар роботов, находящихся на расстоянии <= threshold
func countClosePairs(robots []Robot, threshold int) int {
	count := 0
	for i := 0; i < len(robots); i++ {
		for j := i + 1; j < len(robots); j++ {
			dx := robots[i].x - robots[j].x
			dy := robots[i].y - robots[j].y
			// Используем манхэттенское расстояние или евклидово, в зависимости от предпочтения
			distanceSquared := dx*dx + dy*dy
			if distanceSquared <= threshold*threshold {
				count++
			}
		}
	}
	return count
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: go run main.go <input_file>")
		return
	}
	filename := os.Args[1]
	robotsInitial, err := readInput(filename)
	if err != nil {
		fmt.Printf("Ошибка при чтении входных данных: %v\n", err)
		return
	}

	// Часть Первая
	startPart1 := time.Now()
	robotsAfter100 := simulate(robotsInitial, 100, true)
	safetyFactor := calculateSafetyFactor(robotsAfter100)
	durationPart1 := time.Since(startPart1)
	fmt.Printf("Часть Первая Ответ: %d\n", safetyFactor)
	fmt.Printf("Время выполнения Части Первой: %v\n", durationPart1)

	// Часть Вторая
	startPart2 := time.Now()
	robotsForPart2 := make([]Robot, len(robotsInitial))
	copy(robotsForPart2, robotsInitial)

	threshold := 10 // Пороговое расстояние для близких пар, можно настроить
	maxSeconds := 100000
	maxStable := 100 // Количество последовательных секунд без улучшения для остановки

	maxClosePairs := 0
	minSecond := 0
	stableCount := 0

	for seconds := 1; seconds <= maxSeconds; seconds++ {
		robotsForPart2 = simulate(robotsForPart2, 1, false)
		currentClosePairs := countClosePairs(robotsForPart2, threshold)

		if currentClosePairs > maxClosePairs {
			maxClosePairs = currentClosePairs
			minSecond = seconds
			stableCount = 0
		} else {
			stableCount++
			if stableCount >= maxStable {
				break
			}
		}
	}

	durationPart2 := time.Since(startPart2)
	fmt.Printf("Часть Вторая Ответ: %d секунд\n", minSecond)
	fmt.Printf("Время выполнения Части Второй: %v\n", durationPart2)
}
