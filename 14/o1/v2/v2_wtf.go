/*
➜ go run ./o1/v1 input.txt
Part One Answer: 221655456
Time taken for Part One: 738.541µs
Part Two Answer: 1 seconds
Time taken for Part Two: 25.625µs

нет, не все роботы кучкуются в одну кучу. используй другую метрику
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

// boundingArea вычисляет площадь ограничивающего прямоугольника всех роботов
func boundingArea(robots []Robot) int {
	if len(robots) == 0 {
		return 0
	}
	minX, maxX := robots[0].x, robots[0].x
	minY, maxY := robots[0].y, robots[0].y
	for _, robot := range robots[1:] {
		if robot.x < minX {
			minX = robot.x
		}
		if robot.x > maxX {
			maxX = robot.x
		}
		if robot.y < minY {
			minY = robot.y
		}
		if robot.y > maxY {
			maxY = robot.y
		}
	}
	width := maxX - minX + 1
	height := maxY - minY + 1
	return width * height
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

	prevArea := boundingArea(robotsForPart2)
	seconds := 0
	minArea := prevArea
	minSecond := seconds

	for {
		robotsForPart2 = simulate(robotsForPart2, 1, false)
		seconds++
		currentArea := boundingArea(robotsForPart2)

		// Если площадь уменьшилась, обновляем минимальную площадь и время
		if currentArea < minArea {
			minArea = currentArea
			minSecond = seconds
		} else {
			// Если площадь начала увеличиваться, предполагаем, что минимальная достигнута
			// Дополнительно проверяем несколько шагов назад для уверенности
			// Можно настроить количество проверок
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
