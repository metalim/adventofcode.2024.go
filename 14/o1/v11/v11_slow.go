package main

import (
	"bufio"
	"fmt"
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
	x, y, dx, dy int
}

// parseLine парсит строку ввода и возвращает робота
func parseLine(line string, re *regexp.Regexp) (Robot, error) {
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
func readInput(filename string, re *regexp.Regexp) ([]Robot, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	robots := make([]Robot, 0, 100)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		robot, err := parseLine(line, re)
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
func simulate(robots []Robot, n int, wrap bool) {
	for i := 0; i < n; i++ {
		for j := 0; j < len(robots); j++ {
			robots[j].x += robots[j].dx
			robots[j].y += robots[j].dy

			if wrap {
				// Оборачивание по ширине
				if robots[j].x >= Width {
					robots[j].x -= Width
				} else if robots[j].x < 0 {
					robots[j].x += Width
				}

				// Оборачивание по высоте
				if robots[j].y >= Height {
					robots[j].y -= Height
				} else if robots[j].y < 0 {
					robots[j].y += Height
				}
			}
		}
	}
}

// calculateSafetyFactor вычисляет фактор безопасности после 100 секунд
func calculateSafetyFactor(robots []Robot) int {
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

// calculateAverageSquaredDistance вычисляет среднее квадратичное расстояние между всеми парами роботов
func calculateAverageSquaredDistance(robots []Robot) float64 {
	totalDistanceSquared := 0
	pairs := 0
	n := len(robots)
	for i := 0; i < n; i++ {
		x1, y1 := robots[i].x, robots[i].y
		for j := i + 1; j < n; j++ {
			x2, y2 := robots[j].x, robots[j].y

			dx := x1 - x2
			dy := y1 - y2

			// Кратчайшее расстояние с учетом оборачивания
			if dx > Width/2 {
				dx -= Width
			} else if dx < -Width/2 {
				dx += Width
			}
			if dy > Height/2 {
				dy -= Height
			} else if dy < -Height/2 {
				dy += Height
			}

			// Абсолютные значения
			if dx < 0 {
				dx = -dx
			}
			if dy < 0 {
				dy = -dy
			}

			totalDistanceSquared += dx*dx + dy*dy
			pairs++
		}
	}
	if pairs == 0 {
		return 0
	}
	return float64(totalDistanceSquared) / float64(pairs)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: go run main.go <input_file>")
		return
	}
	filename := os.Args[1]

	// Предварительная компиляция регулярного выражения
	re := regexp.MustCompile(`p=(-?\d+),(-?\d+) v=(-?\d+),(-?\d+)`)

	// Чтение входных данных
	robotsInitial, err := readInput(filename, re)
	if err != nil {
		fmt.Printf("Ошибка при чтении входных данных: %v\n", err)
		return
	}

	// Часть Первая
	startPart1 := time.Now()
	robotsAfter100 := make([]Robot, len(robotsInitial))
	copy(robotsAfter100, robotsInitial)
	simulate(robotsAfter100, 100, true)
	safetyFactor := calculateSafetyFactor(robotsAfter100)
	durationPart1 := time.Since(startPart1)
	fmt.Printf("Часть Первая Ответ: %d\n", safetyFactor)
	fmt.Printf("Время выполнения Части Первой: %v\n", durationPart1)

	// Часть Вторая
	startPart2 := time.Now()
	robotsForPart2 := make([]Robot, len(robotsInitial))
	copy(robotsForPart2, robotsInitial)

	minSumDistanceSquared := int(^uint(0) >> 1) // Max int
	minSecond := 0

	maxSeconds := 200000

	for seconds := 1; seconds <= maxSeconds; seconds++ {
		simulate(robotsForPart2, 1, true)
		sumDistanceSquared := 0
		n := len(robotsForPart2)
		for i := 0; i < n; i++ {
			x1, y1 := robotsForPart2[i].x, robotsForPart2[i].y
			for j := i + 1; j < n; j++ {
				x2, y2 := robotsForPart2[j].x, robotsForPart2[j].y

				dx := x1 - x2
				dy := y1 - y2

				// Кратчайшее расстояние с учетом оборачивания
				if dx > Width/2 {
					dx -= Width
				} else if dx < -Width/2 {
					dx += Width
				}
				if dy > Height/2 {
					dy -= Height
				} else if dy < -Height/2 {
					dy += Height
				}

				// Абсолютные значения
				if dx < 0 {
					dx = -dx
				}
				if dy < 0 {
					dy = -dy
				}

				sumDistanceSquared += dx*dx + dy*dy
			}
		}

		if sumDistanceSquared < minSumDistanceSquared {
			minSumDistanceSquared = sumDistanceSquared
			minSecond = seconds
			fmt.Printf("Новый минимум: %d секунд, Среднее квадратичное расстояние: %.2f\n", minSecond, float64(minSumDistanceSquared)/float64(n*(n-1)/2))
		}
	}

	durationPart2 := time.Since(startPart2)
	fmt.Printf("Часть Вторая Ответ: %d секунд\n", minSecond)
	fmt.Printf("Время выполнения Части Второй: %v\n", durationPart2)
}
