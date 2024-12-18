package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"time"
)

// Константы для размеров сетки
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
// Если wrap=true, позиции оборачиваются по границам сетки.
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

// calculateAverageDistance вычисляет среднее евклидово расстояние между всеми парами роботов
func calculateAverageDistance(robots []Robot) float64 {
	totalDistance := 0.0
	pairs := 0
	for i := 0; i < len(robots); i++ {
		for j := i + 1; j < len(robots); j++ {
			dx := robots[i].x - robots[j].x
			dy := robots[i].y - robots[j].y
			// Учитываем оборачивание для кратчайшего расстояния
			dx = min(abs(dx), Width-abs(dx))
			dy = min(abs(dy), Height-abs(dy))
			distance := math.Sqrt(float64(dx*dx + dy*dy))
			totalDistance += distance
			pairs++
		}
	}
	if pairs == 0 {
		return 0
	}
	return totalDistance / float64(pairs)
}

// abs возвращает абсолютное значение целого числа
func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// min возвращает меньшее из двух целых чисел
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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

	minAverageDistance := math.MaxFloat64
	minSecond := 0
	maxSeconds := 200000

	for seconds := 1; seconds <= maxSeconds; seconds++ {
		robotsForPart2 = simulate(robotsForPart2, 1, true)
		avgDist := calculateAverageDistance(robotsForPart2)

		if avgDist < minAverageDistance {
			minAverageDistance = avgDist
			minSecond = seconds
		}
	}

	durationPart2 := time.Since(startPart2)
	fmt.Printf("Часть Вторая Ответ: %d секунд\n", minSecond)
	fmt.Printf("Время выполнения Части Второй: %v\n", durationPart2)
}
