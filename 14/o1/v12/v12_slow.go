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
	var xs, ys, dxs, dys []int
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
		xs = append(xs, x)
		ys = append(ys, y)
		dxs = append(dxs, dx)
		dys = append(dys, dy)
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Ошибка при чтении файла: %v\n", err)
		return
	}

	n := len(xs)

	// Часть Первая
	startPart1 := time.Now()

	// Обновление позиций после 100 секунд с оборачиванием
	for i := 0; i < n; i++ {
		xs[i] += dxs[i] * 100
		ys[i] += dys[i] * 100

		// Оборачивание по ширине
		if xs[i] >= Width {
			xs[i] %= Width
		} else if xs[i] < 0 {
			xs[i] = (xs[i]%Width + Width) % Width
		}

		// Оборачивание по высоте
		if ys[i] >= Height {
			ys[i] %= Height
		} else if ys[i] < 0 {
			ys[i] = (ys[i]%Height + Height) % Height
		}
	}

	// Вычисление фактора безопасности
	midX := Width / 2
	midY := Height / 2

	q1, q2, q3, q4 := 0, 0, 0, 0
	for i := 0; i < n; i++ {
		if xs[i] == midX || ys[i] == midY {
			continue
		}
		if xs[i] < midX && ys[i] < midY {
			q1++
		} else if xs[i] >= midX && ys[i] < midY {
			q2++
		} else if xs[i] < midX && ys[i] >= midY {
			q3++
		} else if xs[i] >= midX && ys[i] >= midY {
			q4++
		}
	}
	safetyFactor := q1 * q2 * q3 * q4
	durationPart1 := time.Since(startPart1)

	fmt.Printf("Часть Первая Ответ: %d\n", safetyFactor)
	fmt.Printf("Время выполнения Части Первой: %v\n", durationPart1)

	// Часть Вторая
	startPart2 := time.Now()

	// Восстановление начальных позиций для Часть Второй
	// Предполагаем, что начальные позиции уже находятся в xs и ys
	// Поэтому нужно сохранить их перед модификацией
	// Но так как мы уже изменили xs и ys для Часть Первой, необходимо считать заново
	// Для эффективности, лучше считать входные данные дважды
	// Однако, чтобы избежать этого, можно сохранить начальные позиции

	// Восстановление начальных позиций
	// Создаём копии начальных позиций
	initialXs := make([]int, n)
	initialYs := make([]int, n)
	copy(initialXs, xs)
	copy(initialYs, ys)

	// Сброс позиций на начальные значения
	for i := 0; i < n; i++ {
		xs[i] = initialXs[i] - dxs[i]*100
		ys[i] = initialYs[i] - dys[i]*100

		// Оборачивание по ширине
		if xs[i] >= Width {
			xs[i] %= Width
		} else if xs[i] < 0 {
			xs[i] = (xs[i]%Width + Width) % Width
		}

		// Оборачивание по высоте
		if ys[i] >= Height {
			ys[i] %= Height
		} else if ys[i] < 0 {
			ys[i] = (ys[i]%Height + Height) % Height
		}
	}

	minSecond := 0
	minSumDistanceSquared := 1<<63 - 1 // Максимальное значение int64

	maxSeconds := 200000

	for seconds := 1; seconds <= maxSeconds; seconds++ {
		// Обновление позиций на 1 секунду с оборачиванием
		for i := 0; i < n; i++ {
			xs[i] += dxs[i]
			ys[i] += dys[i]

			// Оборачивание по ширине
			if xs[i] >= Width {
				xs[i] -= Width
			} else if xs[i] < 0 {
				xs[i] += Width
			}

			// Оборачивание по высоте
			if ys[i] >= Height {
				ys[i] -= Height
			} else if ys[i] < 0 {
				ys[i] += Height
			}
		}

		// Вычисление суммы квадратов расстояний
		sumDistanceSquared := 0
		for i := 0; i < n; i++ {
			x1, y1 := xs[i], ys[i]
			for j := i + 1; j < n; j++ {
				dx := x1 - xs[j]
				dy := y1 - ys[j]

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
			// Выводим только новый минимум
			fmt.Printf("Новый минимум: %d секунд, Среднее квадратичное расстояние: %.2f\n", minSecond, float64(minSumDistanceSquared)/float64(n*(n-1)/2))
		}
	}

	durationPart2 := time.Since(startPart2)
	fmt.Printf("Часть Вторая Ответ: %d секунд\n", minSecond)
	fmt.Printf("Время выполнения Части Второй: %v\n", durationPart2)
}
