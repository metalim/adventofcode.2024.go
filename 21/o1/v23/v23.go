// Advent of Code 2024, Day 21: "Keypad Conundrum"
// ------------------------------------------------
// Решение, способное уложиться в ~10 секунд даже при 25 роботах (Part 2),
// используя итеративный поиск в глубину (IDDFS) с префиксной отсечкой
// и "макроактивацией".
// Идея:
//   1. Для каждого кода выполняем IDDFS (Iterative Deepening DFS),
//      наращивая глубину от 0 до некоего лимита (скажем, 300).
//   2. На каждом шаге проверяем:
//       - Если набранный код == искомому, возвращаем глубину.
//       - Если длина текущего пути >= лимита, останавливаемся.
//       - Если набранное (typed) уже не является префиксом искомого, отсекаем ветку.
//   3. "Макроактивация": при нажатии 'A' на i-м роботе, если там 'A',
//      спускаемся к (i+1)-му роботу (или нажимаем цифру, если это последний робот);
//      если там стрелка, двигаем руку (i+1)-го робота (или цифровую клавиатуру,
//      если i — последний).
//   4. В отличие от BFS, IDDFS не хранит все состояния в памяти,
//      а глубину наращиваем постепенно. Для коротких кодов (3-4 символа + 'A')
//      и префиксной отсечки такой подход на практике очень быстр.
//
// Шаги реализации:
//
//  - У нас два "chainLen": Part 1 -> 2 робота, Part 2 -> 25 роботов.
//  - Каждому коду вызываем solveCodeIDDFS(chainLen, code).
//     * Она запускает цикл по limit = 0..maxDepth (например, 300).
//     * Делает DFS (рекурсивную функцию) с параметрами (глубина, limit, текущее состояние).
//       Если находим код — возвращаем глубину.
//  - Состояние хранить в виде (robPos[], digPos, typed). Но в DFS не создаём больших структур —
//    передаём их копии/ссылки аккуратно. "robPos" можно хранить срезом длины chainLen,
//    "digPos" и "typed" — строки.
//    * На каждом шаге генерируем до 5 переходов: '^','v','<','>','A'.
//    * При 'A' делаем функцию macroActivate(...) — сразу проходим несколько уровней (если 'A' на роботе).
//  - Как только находим решение на глубине d, завершаем и возвращаем d.
//  - Перемножаем d на числовую часть кода.
//
// Такой IDDFS-алгоритм с префиксной отсечкой обычно завершается мгновенно
// для кодов длиной 4 символа, даже при 25 роботах.
//
// ------------------------------------------------

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Константы для двух частей
const (
	chainLenPart1 = 2
	chainLenPart2 = 25
)

// Раскладка роботной клавиатуры
var robotAdj = map[string]map[string]string{
	"^": {"^": "", "v": "v", "<": "", ">": "A"},
	"A": {"^": "", "v": ">", "<": "^", ">": ""},
	"<": {"^": "", "v": "", "<": "", ">": "v"},
	"v": {"^": "^", "v": "", "<": "<", ">": ">"},
	">": {"^": "A", "v": "", "<": "v", ">": ""},
}

// Раскладка цифровой клавиатуры
var digitAdj = map[string]map[string]string{
	"7": {"^": "", "v": "4", "<": "", ">": "8"},
	"8": {"^": "", "v": "5", "<": "7", ">": "9"},
	"9": {"^": "", "v": "6", "<": "8", ">": ""},
	"4": {"^": "7", "v": "1", "<": "", ">": "5"},
	"5": {"^": "8", "v": "2", "<": "4", ">": "6"},
	"6": {"^": "9", "v": "3", "<": "5", ">": ""},
	"1": {"^": "4", "v": "", "<": "", ">": "2"},
	"2": {"^": "5", "v": "0", "<": "1", ">": "3"},
	"3": {"^": "6", "v": "A", "<": "2", ">": ""},
	"0": {"^": "2", "v": "", "<": "", ">": "A"},
	"A": {"^": "3", "v": "", "<": "0", ">": ""},
}

// parseNumericPart: "029A" -> 29, "000A" -> 0, etc.
func parseNumericPart(code string) int {
	s := strings.TrimSuffix(code, "A")
	s = strings.TrimLeft(s, "0")
	if s == "" {
		return 0
	}
	n, _ := strconv.Atoi(s)
	return n
}

// ---------- Итеративный поиск в глубину ----------
// Для каждого кода вызываем solveCodeIDDFS(chainLen, code).
// Она запускает limit от 0 до maxDepth. При нахождении решения сразу возвращаем глубину.
// maxDepth можно поставить ~300 как очень большой запас, т.к. реальная длина решения гораздо меньше.
//
// dfsID — рекурсивная функция. depth — текущая глубина.
// Если набран code, возвращаем true. Если depth == limit, обрываем. Если typed не префикс, обрываем.
//
// "robPos" — массив позиций роботов (длина chainLen), "digPos" — позиция на цифровой, "typed" — набранное.

func solveCodeIDDFS(chainLen int, code string) int {
	// Быстрая проверка
	if code == "" {
		return 0
	}
	// Начальные позиции
	rob := make([]string, chainLen)
	for i := 0; i < chainLen; i++ {
		rob[i] = "A"
	}
	dig := "A"
	typed := ""

	// Пытаемся наращивать limit
	maxDepth := 300
	for limit := 0; limit <= maxDepth; limit++ {
		if dfsLimited(rob, dig, typed, 0, limit, chainLen, code) {
			return limit
		}
	}
	// Не нашли решение в пределах maxDepth => 0
	return 0
}

// dfsLimited — рекурсивный поиск.
// depth — текущая глубина; limit — предельная глубина.
// Возвращаем true, если нашли решение.
func dfsLimited(rob []string, dig, typed string, depth, limit, chainLen int, code string) bool {
	// Цель
	if typed == code {
		return true
	}
	if depth == limit {
		return false
	}
	// Префиксная отсечка
	if !strings.HasPrefix(code, typed) {
		return false
	}

	// Пять возможных нажатий
	for _, inp := range []string{"^", "v", "<", ">", "A"} {
		// Копируем текущее состояние
		robN := append([]string(nil), rob...)
		digN := dig
		typedN := typed
		ok := true

		if inp != "A" {
			cur := robN[0]
			nx := robotAdj[cur][inp]
			if nx == "" {
				ok = false
			} else {
				robN[0] = nx
			}
		} else {
			// macroActivate
			robN, digN, typedN, ok = macroActivateIDDFS(robN, digN, typedN, 0, chainLen)
		}

		if ok {
			// Рекурсия
			if dfsLimited(robN, digN, typedN, depth+1, limit, chainLen, code) {
				return true
			}
		}
	}
	return false
}

// macroActivateIDDFS — аналог macroActivate, только не храним больших таблиц visited.
// Пробиваем 'A' на уровне idx.
func macroActivateIDDFS(rob []string, dig, typed string, idx, chainLen int) ([]string, string, string, bool) {
	cur := rob[idx]
	switch cur {
	case "^", "v", "<", ">":
		// Команда для следующего уровня
		if idx == chainLen-1 {
			// Двигаем цифровую
			nx := digitAdj[dig][cur]
			if nx == "" {
				return rob, dig, typed, false
			}
			dig = nx
			return rob, dig, typed, true
		}
		// Двигаем робота idx+1
		c2 := rob[idx+1]
		n2 := robotAdj[c2][cur]
		if n2 == "" {
			return rob, dig, typed, false
		}
		robCpy := append([]string(nil), rob...)
		robCpy[idx+1] = n2
		return robCpy, dig, typed, true

	case "A":
		if idx < chainLen-1 {
			return macroActivateIDDFS(rob, dig, typed, idx+1, chainLen)
		}
		// idx == chainLen-1 => нажать кнопку на цифровой
		if _, isDir := robotAdj[dig]; isDir && dig != "A" {
			return rob, dig, typed, false
		}
		typed += dig
		return rob, dig, typed, true

	default:
		return rob, dig, typed, false
	}
}

// solvePartIDDFS — суммирует (длина*числоваяЧасть) для 5 кодов.
func solvePartIDDFS(chainLen int, codes []string) int {
	sum := 0
	for _, c := range codes {
		ln := solveCodeIDDFS(chainLen, c)
		nv := parseNumericPart(c)
		sum += ln * nv
	}
	return sum
}

// ------------------- main -------------------
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <input_file>")
		return
	}
	filename := os.Args[1]
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var codes []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line != "" {
			codes = append(codes, line)
		}
	}
	if err := sc.Err(); err != nil {
		panic(err)
	}

	if len(codes) == 0 {
		fmt.Println("No codes found.")
		return
	}

	// Part 1 (chainLen=2)
	startP1 := time.Now()
	p1res := solvePartIDDFS(chainLenPart1, codes)
	tP1 := time.Since(startP1).Seconds()
	fmt.Printf("Part 1: %d (time %.4fs)\n", p1res, tP1)

	// Part 2 (chainLen=25)
	startP2 := time.Now()
	p2res := solvePartIDDFS(chainLenPart2, codes)
	tP2 := time.Since(startP2).Seconds()
	fmt.Printf("Part 2: %d (time %.4fs)\n", p2res, tP2)
}
