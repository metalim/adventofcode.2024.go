/*
Напиши код на Go для решения следующей задачи. Входные данные в файле указываемом аргументом в командной строке.
Выведи время решения каждой части.

--- Day 9: Disk Fragmenter ---
Another push of the button leaves you in the familiar hallways of some friendly amphipods! Good thing you each somehow got your own personal mini submarine. The Historians jet away in search of the Chief, mostly by driving directly into walls.

While The Historians quickly figure out how to pilot these things, you notice an amphipod in the corner struggling with his computer. He's trying to make more contiguous free space by compacting all of the files, but his program isn't working; you offer to help.

He shows you the disk map (your puzzle input) he's already generated. For example:

2333133121414131402
The disk map uses a dense format to represent the layout of files and free space on the disk. The digits alternate between indicating the length of a file and the length of free space.

So, a disk map like 12345 would represent a one-block file, two blocks of free space, a three-block file, four blocks of free space, and then a five-block file. A disk map like 90909 would represent three nine-block files in a row (with no free space between them).

Each file on disk also has an ID number based on the order of the files as they appear before they are rearranged, starting with ID 0. So, the disk map 12345 has three files: a one-block file with ID 0, a three-block file with ID 1, and a five-block file with ID 2. Using one character for each block where digits are the file ID and . is free space, the disk map 12345 represents these individual blocks:

0..111....22222
The first example above, 2333133121414131402, represents these individual blocks:

00...111...2...333.44.5555.6666.777.888899
The amphipod would like to move file blocks one at a time from the end of the disk to the leftmost free space block (until there are no gaps remaining between file blocks). For the disk map 12345, the process looks like this:

0..111....22222
02.111....2222.
022111....222..
0221112...22...
02211122..2....
022111222......
The first example requires a few more steps:

00...111...2...333.44.5555.6666.777.888899
009..111...2...333.44.5555.6666.777.88889.
0099.111...2...333.44.5555.6666.777.8888..
00998111...2...333.44.5555.6666.777.888...
009981118..2...333.44.5555.6666.777.88....
0099811188.2...333.44.5555.6666.777.8.....
009981118882...333.44.5555.6666.777.......
0099811188827..333.44.5555.6666.77........
00998111888277.333.44.5555.6666.7.........
009981118882777333.44.5555.6666...........
009981118882777333644.5555.666............
00998111888277733364465555.66.............
0099811188827773336446555566..............
The final step of this file-compacting process is to update the filesystem checksum. To calculate the checksum, add up the result of multiplying each of these blocks' position with the file ID number it contains. The leftmost block is in position 0. If a block contains free space, skip it instead.

Continuing the first example, the first few blocks' position multiplied by its file ID number are 0 * 0 = 0, 1 * 0 = 0, 2 * 9 = 18, 3 * 9 = 27, 4 * 8 = 32, and so on. In this example, the checksum is the sum of these, 1928.

Compact the amphipod's hard drive using the process he requested. What is the resulting filesystem checksum? (Be careful copy/pasting the input for this puzzle; it is a single, very long line.)

--- Part Two ---
Upon completion, two things immediately become clear. First, the disk definitely has a lot more contiguous free space, just like the amphipod hoped. Second, the computer is running much more slowly! Maybe introducing all of that file system fragmentation was a bad idea?

The eager amphipod already has a new plan: rather than move individual blocks, he'd like to try compacting the files on his disk by moving whole files instead.

This time, attempt to move whole files to the leftmost span of free space blocks that could fit the file. Attempt to move each file exactly once in order of decreasing file ID number starting with the file with the highest file ID number. If there is no span of free space to the left of a file that is large enough to fit the file, the file does not move.

The first example from above now proceeds differently:

00...111...2...333.44.5555.6666.777.888899
0099.111...2...333.44.5555.6666.777.8888..
0099.1117772...333.44.5555.6666.....8888..
0099.111777244.333....5555.6666.....8888..
00992111777.44.333....5555.6666.....8888..
The process of updating the filesystem checksum is the same; now, this example's checksum would be 2858.

Start over, now compacting the amphipod's hard drive using this new method instead. What is the resulting filesystem checksum?

*/

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
