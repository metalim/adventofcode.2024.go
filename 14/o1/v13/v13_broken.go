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

// Функция для парсинга строки и добавления данных в срезы
func parseLine(line string, re *regexp.Regexp, xs, ys, dxs, dys *[]int) error {
	matches := re.FindStringSubmatch(line)
	if matches == nil || len(matches) != 5 {
		return fmt.Errorf("неверный формат строки: %s", line)
	}
	x, err := strconv.Atoi(matches[1])
	if err != nil {
		return err
	}
	y, err := strconv.Atoi(matches[2])
	if err != nil {
		return err
	}
	dx, err := strconv.Atoi(matches[3])
	if err != nil {
		return err
	}
	dy, err := strconv.Atoi(matches[4])
	if err != nil {
		return err
	}
	*xs = append(*xs, x)
	*ys = append(*ys, y)
	*dxs = append(*dxs, dx)
	*dys = append(*dys, dy)
	return nil
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
		err := parseLine(line, re, &xs, &ys, &dxs, &dys)
		if err != nil {
			fmt.Printf("Ошибка при парсинге строки: %v\n", err)
			return
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Ошибка при чтении файла: %v\n", err)
		return
	}

	n := len(xs)

	// Часть Первая
	startPart1 := time.Now()

	for i := 0; i < n; i++ {
		xs[i] += dxs[i] * 100
		ys[i] += dys[i] * 100

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
	// Так как позиции уже изменены для Часть Первой, необходимо считать их заново
	// Для оптимизации, можно считать входные данные дважды или сохранить начальные позиции
	// Здесь мы предположим, что позиции уже были сохранены до изменения

	// Для простоты, переоткрываем файл и читаем заново
	file.Seek(0, 0) // Сброс позиции чтения файла
	xs = nil
	ys = nil
	dxs = nil
	dys = nil
	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		err := parseLine(line, re, &xs, &ys, &dxs, &dys)
		if err != nil {
			fmt.Printf("Ошибка при парсинге строки: %v\n", err)
			return
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Ошибка при чтении файла: %v\n", err)
		return
	}

	// Инициализация минимальной суммы квадратов расстояний
	minSumDistanceSquared := 1<<63 - 1 // Максимальное значение int64
	minSecond := 0

	// Параметры для раннего завершения симуляции
	patience := 1000     // Количество секунд без улучшений для остановки
	noImprovement := 0   // Счётчик секунд без улучшений
	maxSeconds := 200000 // Максимальное количество секунд симуляции

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
			noImprovement = 0
			// Вывод нового минимума
			fmt.Printf("Новый минимум: %d секунд, Среднее квадратичное расстояние: %.2f\n", minSecond, float64(minSumDistanceSquared)/float64(n*(n-1)/2))
		} else {
			noImprovement++
			if noImprovement >= patience {
				// Если нет улучшений в течение `patience` секунд, завершаем симуляцию
				break
			}
		}
	}

	durationPart2 := time.Since(startPart2)
	fmt.Printf("Часть Вторая Ответ: %d секунд\n", minSecond)
	fmt.Printf("Время выполнения Части Второй: %v\n", durationPart2)
}
