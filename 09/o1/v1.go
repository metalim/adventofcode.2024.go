// Напиши код на Go для решения следующей задачи. Входные данные в файле указываемом аргументом в командной строке.
// Выведи время решения каждой части.

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <inputfile>")
	}

	startOverall := time.Now()

	// Чтение входных данных
	inputFile := os.Args[1]
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()

	// Парсинг входной строки
	var blocks []int
	fileID := 0
	isFile := true
	for i := 0; i < len(line); i++ {
		n, _ := strconv.Atoi(string(line[i]))
		if isFile {
			for j := 0; j < n; j++ {
				blocks = append(blocks, fileID)
			}
			fileID++
		} else {
			for j := 0; j < n; j++ {
				blocks = append(blocks, -1)
			}
		}
		isFile = !isFile
	}
	original := make([]int, len(blocks))
	copy(original, blocks)

	// Часть 1: Посимвольное перемещение
	startPart1 := time.Now()
	blocksPart1 := make([]int, len(original))
	copy(blocksPart1, original)

	for {
		leftFree := -1
		rightFile := -1
		// Находим левый свободный блок и правый файловый блок
		for i, b := range blocksPart1 {
			if b == -1 {
				leftFree = i
				break
			}
		}
		if leftFree == -1 {
			break
		}
		for i := len(blocksPart1) - 1; i >= 0; i-- {
			if blocksPart1[i] >= 0 {
				rightFile = i
				break
			}
		}
		if rightFile == -1 || rightFile < leftFree {
			break
		}
		// Перемещение
		blocksPart1[leftFree] = blocksPart1[rightFile]
		blocksPart1[rightFile] = -1
	}
	// Вычисляем чек-сумму
	checksum1 := 0
	for i, b := range blocksPart1 {
		if b >= 0 {
			checksum1 += i * b
		}
	}
	timePart1 := time.Since(startPart1)

	// Часть 2: Перемещение целых файлов
	startPart2 := time.Now()
	blocksPart2 := make([]int, len(original))
	copy(blocksPart2, original)

	// Определяем файлы (start, end)
	type filePos struct{ start, end int }
	files := []filePos{}
	curID := 0
	var startPos = -1
	for i, b := range blocksPart2 {
		if b == curID {
			if startPos == -1 {
				startPos = i
			}
		} else {
			if startPos != -1 {
				files = append(files, filePos{startPos, i - 1})
				startPos = -1
				curID++
			}
			if b == curID {
				startPos = i
			}
		}
	}
	if startPos != -1 {
		files = append(files, filePos{startPos, len(blocksPart2) - 1})
	}

	// Двигаемся по файлам в порядке убывания ID
	for id := len(files) - 1; id >= 0; id-- {
		f := files[id]
		length := f.end - f.start + 1
		// Ищем свободный сегмент слева от f.start, достаточно длинный
		target := -1
		for i := 0; i < f.start; {
			if blocksPart2[i] == -1 {
				segStart := i
				for i < len(blocksPart2) && blocksPart2[i] == -1 {
					i++
				}
				segLen := i - segStart
				if segLen >= length {
					target = segStart
					break
				}
			} else {
				i++
			}
		}
		// Переносим, если нашли место
		if target != -1 {
			// Освобождаем старые блоки
			fileBlocks := make([]int, length)
			for k := f.start; k <= f.end; k++ {
				fileBlocks[k-f.start] = blocksPart2[k]
				blocksPart2[k] = -1
			}
			// Копируем файл в новое место
			for k := 0; k < length; k++ {
				blocksPart2[target+k] = fileBlocks[k]
			}
		}
	}

	// Чек-сумма для второй части
	checksum2 := 0
	for i, b := range blocksPart2 {
		if b >= 0 {
			checksum2 += i * b
		}
	}
	timePart2 := time.Since(startPart2)

	// Выводим результаты
	fmt.Printf("Part 1: %d (time: %v)\n", checksum1, timePart1)
	fmt.Printf("Part 2: %d (time: %v)\n", checksum2, timePart2)
	fmt.Printf("Total time: %v\n", time.Since(startOverall))
}
