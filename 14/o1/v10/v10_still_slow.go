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
	x  int
	y  int
	dx int
	dy int
}

// Pair представляет пару индексов роботов
type Pair struct {
	i int
	j int
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
				robots[j].x %= Width
				if robots[j].x < 0 {
					robots[j].x += Width
				}
				robots[j].y %= Height
				if robots[j].y < 0 {
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
func calculateAverageSquaredDistance(robots []Robot, pairs []Pair) float64 {
	var totalDistanceSquared int
	for _, pair := range pairs {
		dx := robots[pair.i].x - robots[pair.j].x
		dy := robots[pair.i].y - robots[pair.j].y

		// Вычисляем кратчайшее расстояние с учетом оборачивания
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

		distanceSquared := dx*dx + dy*dy
		totalDistanceSquared += distanceSquared
	}
	numPairs := len(pairs)
	if numPairs == 0 {
		return 0
	}
	return float64(totalDistanceSquared) / float64(numPairs)
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

	// Предварительное создание списка пар роботов
	numRobots := len(robotsForPart2)
	pairs := make([]Pair, 0, numRobots*(numRobots-1)/2)
	for i := 0; i < numRobots; i++ {
		for j := i + 1; j < numRobots; j++ {
			pairs = append(pairs, Pair{i: i, j: j})
		}
	}

	minSumDistanceSquared := int(^uint(0) >> 1) // Max int
	minSecond := 0
	maxSeconds := 200000

	for seconds := 1; seconds <= maxSeconds; seconds++ {
		simulate(robotsForPart2, 1, true)
		sumDistanceSquared := 0

		for _, pair := range pairs {
			dx := robotsForPart2[pair.i].x - robotsForPart2[pair.j].x
			dy := robotsForPart2[pair.i].y - robotsForPart2[pair.j].y

			// Вычисляем кратчайшее расстояние с учетом оборачивания
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

			distanceSquared := dx*dx + dy*dy
			sumDistanceSquared += distanceSquared
		}

		if sumDistanceSquared < minSumDistanceSquared {
			minSumDistanceSquared = sumDistanceSquared
			minSecond = seconds
		}
	}

	durationPart2 := time.Since(startPart2)
	fmt.Printf("Часть Вторая Ответ: %d секунд\n", minSecond)
	fmt.Printf("Время выполнения Части Второй: %v\n", durationPart2)
}
