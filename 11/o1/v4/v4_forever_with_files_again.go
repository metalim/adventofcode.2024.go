/*
```
➜ go run ./o1/v3/ sample2.txt
3
92.167µs
3
19µs
```

ответ неверный
*/

package main

import (
	"bufio"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"
)

/*
Другой подход: полностью потоковая обработка с минимальным использованием памяти.

Идея:
- На каждом шаге читать один файл построчно (каждый камень — отдельная строка).
- Для каждого камня применить правило и сразу записать результаты в новый файл, не храня всё в памяти.
- После окончания шага заменить входной файл на выходной.
- Так мы избежим взрыва памяти, используя только файлы и простые операции.

Оптимизации:
- Не храним все камни в памяти, работаем построчно.
- Используем big.Int для умножения, так как иных данных у нас нет.
- Правило разбиения: просто делим строку пополам, trimLeft для удаления ведущих нулей.
- Да, будет долго, но не будет переполнения памяти.

*/

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <inputfile>")
	}
	inputFile := os.Args[1]

	// Считаем камни после 25 шагов
	start := time.Now()
	tmpAfter25 := evolveFile(inputFile, 25)
	count25 := countLines(tmpAfter25)
	fmt.Println(count25)
	fmt.Println(time.Since(start))

	// Считаем камни после 75 шагов (ещё 50 шагов)
	start2 := time.Now()
	tmpAfter75 := evolveFile(tmpAfter25, 50)
	count75 := countLines(tmpAfter75)
	fmt.Println(count75)
	fmt.Println(time.Since(start2))
}

func evolveFile(inputFile string, times int) string {
	in := inputFile
	for i := 0; i < times; i++ {
		out := in + ".out"
		processFile(in, out)
		os.Remove(in)
		os.Rename(out, in)
	}
	return in
}

func processFile(inFile, outFile string) {
	fi, err := os.Open(inFile)
	if err != nil {
		log.Fatal(err)
	}
	defer fi.Close()

	fo, err := os.Create(outFile)
	if err != nil {
		log.Fatal(err)
	}
	defer fo.Close()

	scanner := bufio.NewScanner(fi)
	writer := bufio.NewWriter(fo)

	m2024 := big.NewInt(2024)

	for scanner.Scan() {
		num := scanner.Text()
		if num == "0" {
			// Правило 1
			writer.WriteString("1\n")
		} else {
			L := len(num)
			if L%2 == 0 {
				// Правило 2
				half := L / 2
				leftStr := strings.TrimLeft(num[:half], "0")
				if leftStr == "" {
					leftStr = "0"
				}
				rightStr := strings.TrimLeft(num[half:], "0")
				if rightStr == "" {
					rightStr = "0"
				}
				writer.WriteString(leftStr + "\n")
				writer.WriteString(rightStr + "\n")
			} else {
				// Правило 3
				bi := new(big.Int)
				bi.SetString(num, 10)
				bi.Mul(bi, m2024)
				writer.WriteString(bi.String() + "\n")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	writer.Flush()
}

func countLines(filename string) int {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	count := 0
	for scanner.Scan() {
		count++
	}
	return count
}
