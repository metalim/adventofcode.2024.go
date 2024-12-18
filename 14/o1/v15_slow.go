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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: go run main.go <input_file>")
		return
	}
	filename := os.Args[1]

	// Предварительная компиляция регулярного выражения
	re := regexp.MustCompile(`p=(-?\d+),(-?\d+) v=(-?\d+),(-?\d+)`)

	// Чтение входных данных
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Ошибка при открытии файла: %v\n", err)
		return
	}
	defer file.Close()

	// Инициализация срезов для координат и скоростей
	var xs_init, ys_init, dxs, dys []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if matches == nil || len(matches) != 5 {
			fmt.Printf("Неверный формат строки: %s\n", line)
			return
		}
		x, err := strconv.Atoi(matches[1])
		if err != nil {
			fmt.Printf("Ошибка при парсинге x: %v\n", err)
			return
		}
		y, err := strconv.Atoi(matches[2])
		if err != nil {
			fmt.Printf("Ошибка при парсинге y: %v\n", err)
			return
		}
		dx, err := strconv.Atoi(matches[3])
		if err != nil {
			fmt.Printf("Ошибка при парсинге dx: %v\n", err)
			return
		}
		dy, err := strconv.Atoi(matches[4])
		if err != nil {
			fmt.Printf("Ошибка при парсинге dy: %v\n", err)
			return
		}
		xs_init = append(xs_init, x)
		ys_init = append(ys_init, y)
		dxs = append(dxs, dx)
		dys = append(dys, dy)
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Ошибка при чтении файла: %v\n", err)
		return
	}

	n := len(xs_init)
	if n == 0 {
		fmt.Println("Нет роботов для симуляции.")
		return
	}

	// Часть Первая
	startPart1 := time.Now()

	// Копируем начальные позиции для части первой
	xs_p1 := make([]int, n)
	ys_p1 := make([]int, n)
	copy(xs_p1, xs_init)
	copy(ys_p1, ys_init)

	// Симуляция 100 секунд
	for i := 0; i < 100; i++ {
		for j := 0; j < n; j++ {
			xs_p1[j] += dxs[j]
			ys_p1[j] += dys[j]

			// Оборачивание по ширине
			if xs_p1[j] >= Width {
				xs_p1[j] -= Width
			} else if xs_p1[j] < 0 {
				xs_p1[j] += Width
			}

			// Оборачивание по высоте
			if ys_p1[j] >= Height {
				ys_p1[j] -= Height
			} else if ys_p1[j] < 0 {
				ys_p1[j] += Height
			}
		}
	}

	// Вычисление фактора безопасности
	midX := Width / 2
	midY := Height / 2
	q1, q2, q3, q4 := 0, 0, 0, 0
	for j := 0; j < n; j++ {
		if xs_p1[j] == midX || ys_p1[j] == midY {
			continue
		}
		if xs_p1[j] < midX && ys_p1[j] < midY {
			q1++
		} else if xs_p1[j] >= midX && ys_p1[j] < midY {
			q2++
		} else if xs_p1[j] < midX && ys_p1[j] >= midY {
			q3++
		} else {
			q4++
		}
	}
	safetyFactor := q1 * q2 * q3 * q4
	durationPart1 := time.Since(startPart1)

	fmt.Printf("Часть Первая Ответ: %d\n", safetyFactor)
	fmt.Printf("Время выполнения Части Первой: %v\n", durationPart1)

	// Часть Вторая
	startPart2 := time.Now()

	// Копируем начальные позиции для части второй
	xs_p2 := make([]int, n)
	ys_p2 := make([]int, n)
	copy(xs_p2, xs_init)
	copy(ys_p2, ys_init)

	minSecond := 0
	minSumDistanceSquared := int(^uint(0) >> 1) // Максимальное значение int
	maxSeconds := 200000

	for seconds := 1; seconds <= maxSeconds; seconds++ {
		// Обновление позиций на 1 секунду
		for j := 0; j < n; j++ {
			xs_p2[j] += dxs[j]
			ys_p2[j] += dys[j]

			// Оборачивание по ширине
			if xs_p2[j] >= Width {
				xs_p2[j] -= Width
			} else if xs_p2[j] < 0 {
				xs_p2[j] += Width
			}

			// Оборачивание по высоте
			if ys_p2[j] >= Height {
				ys_p2[j] -= Height
			} else if ys_p2[j] < 0 {
				ys_p2[j] += Height
			}
		}

		// Вычисление суммы квадратов расстояний
		sumDistanceSquared := 0
		for i := 0; i < n; i++ {
			x1, y1 := xs_p2[i], ys_p2[i]
			for j := i + 1; j < n; j++ {
				dx := x1 - xs_p2[j]
				dy := y1 - ys_p2[j]

				// Кратчайшее расстояние с учётом оборачивания
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

		// Проверка на новый минимум
		if sumDistanceSquared < minSumDistanceSquared {
			minSumDistanceSquared = sumDistanceSquared
			minSecond = seconds
			avgDistSquared := float64(sumDistanceSquared) / float64(n*(n-1)/2)
			fmt.Printf("Новый минимум: %d секунд, Среднее квадратичное расстояние: %.2f\n", minSecond, avgDistSquared)
		}
	}

	durationPart2 := time.Since(startPart2)
	fmt.Printf("Часть Вторая Ответ: %d секунд\n", minSecond)
	fmt.Printf("Время выполнения Части Второй: %v\n", durationPart2)
}
