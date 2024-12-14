/*
➜ go run ./o1/v2 input.txt
Часть Первая Ответ: 221655456
Время выполнения Части Первой: 735.666µs
Часть Вторая Ответ: 0 секунд
Время выполнения Части Второй: 13µs

ты еблан? симуляция была нормальная в первой версии, но предположение о нахождении минимального bounding box было неверным. Используй ДРУГУЮ метрику для определения ЁЛКИ. Не нарушай условий задачи
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

// calculateVerticalSpread вычисляет вертикальное расстояние между верхним и нижним роботом
func calculateVerticalSpread(robots []Robot) int {
	if len(robots) == 0 {
		return 0
	}
	minY, maxY := robots[0].y, robots[0].y
	for _, robot := range robots[1:] {
		if robot.y < minY {
			minY = robot.y
		}
		if robot.y > maxY {
			maxY = robot.y
		}
	}
	return maxY - minY
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

	prevSpread := calculateVerticalSpread(robotsForPart2)
	seconds := 0
	minSpread := prevSpread
	minSecond := seconds

	for {
		robotsForPart2 = simulate(robotsForPart2, 1, false)
		seconds++
		currentSpread := calculateVerticalSpread(robotsForPart2)

		if currentSpread < minSpread {
			minSpread = currentSpread
			minSecond = seconds
		} else if currentSpread > minSpread {
			// Предполагаем, что минимальная точка была на предыдущей секунде
			break
		}

		// Для предотвращения бесконечного цикла устанавливаем верхний предел
		if seconds > 200000 {
			fmt.Println("Не удалось найти момент отображения пасхального яйца в пределах 200,000 секунд.")
			break
		}
	}

	durationPart2 := time.Since(startPart2)
	fmt.Printf("Часть Вторая Ответ: %d секунд\n", minSecond)
	fmt.Printf("Время выполнения Части Второй: %v\n", durationPart2)
}
