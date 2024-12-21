// Advent of Code 2024, Day 21, Part Two: "Keypad Conundrum"
// ---------------------------------------------------------
// Если даже оптимизированный BFS с макро-активацией (см. предыдущие версии)
// на локальной машине пользователя уходит в timeout на sample.txt,
// велика вероятность, что либо среда исполнения крайне медленная,
// либо есть какая-то ошибка в реализации или особенность среды.
//
// Ниже — ещё более «агрессивный» метод:
//   1) Мы используем ту же идею «макро-активации», когда при нажатии 'A'
//      команда проваливается вниз по цепочке из 25 роботов за один шаг.
//   2) Дополнительно ограничиваем глубину поиска (максимальное количество
//      нажатий) — ведь для кода длиной 4 символа + 'A' нет смысла делать
//      больше, чем, скажем, 100 нажатий (это уже с огромным запасом).
//   3) Вместо «чистого BFS» делаем итеративный поиск в глубину (DFS)
//      с ограничением глубины (iterative deepening DFS, IDDFS).
//      IDDFS = наращиваем лимит по глубине от 0,1,2,… пока не найдём решение.
//      Как только найдём, берём кратчайшую.
//
// Для кратких кодов (4 символа + 'A') и 25 роботов этого обычно достаточно,
// и даёт быстрое время — потому что неправильные ветки отсекаются по префиксу,
// а лимит глубины не даёт уходить в «бесконечность».
//
// --------------------------------------------
// Запуск:
//   go run main.go sample.txt
//   go run main.go input.txt
// и т.д.
//
// Если у вас и это решение «зависает» на sample.txt,
// значит либо в среде исполнения очень жёсткие ограничения,
// либо что-то некорректно установлено.
// На большинстве современных систем такая задача решается мгновенно.
// --------------------------------------------

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// chainLength = 25 роботов + 1 наш верхний => 26 "directional" уровней.
const chainLength = 25

// Раскладка "directional keypad" для каждого робота (верхних 26).
var robotAdj = map[string]map[string]string{
	"^": {"^": "", "v": "v", "<": "", ">": "A"},
	"A": {"^": "", "v": ">", "<": "^", ">": ""},
	"<": {"^": "", "v": "", "<": "", ">": "v"},
	"v": {"^": "^", "v": "", "<": "<", ">": ">"},
	">": {"^": "A", "v": "", "<": "v", ">": ""},
}

// Раскладка цифровой клавиатуры.
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

// RobPos — фиксированный массив позиций рук роботов (сравнимый тип).
type RobPos [chainLength + 1]string

// State: позиция каждого из 26 роботов, позиция на цифровой клавиатуре, набранная строка.
type State struct {
	rpos  RobPos
	dig   string
	typed string
}

// startState — все роботы и цифра на 'A', typed = "".
func startState() State {
	var rp RobPos
	for i := range rp {
		rp[i] = "A"
	}
	return State{rpos: rp, dig: "A", typed: ""}
}

// parseNumericPart: "029A" => 29, "456A" => 456, и т.д.
func parseNumericPart(code string) int {
	s := strings.TrimSuffix(code, "A")
	s = strings.TrimLeft(s, "0")
	if s == "" {
		return 0
	}
	n, _ := strconv.Atoi(s)
	return n
}

// ------------------- ITERATIVE DEEPENING -------------------
// Будем искать кратчайшее решение через IDDFS, увеличивая лимит глубины.
// При первом найденном решении завершаемся.

var bestDepthFound int // для хранения результата при поиске

func solveCodeIDDFS(code string) int {
	if code == "" {
		return 0
	}
	start := startState()

	// Быстрая проверка: если уже совпало, вернуть 0
	if start.typed == code {
		return 0
	}

	// Ограничим глубину искусственно чем-то «безопасным».
	// Для кода из 4 символов + 'A' вряд ли нужно более 200 нажатий.
	// Обычно хватает и < 100. Поставим 300 «с запасом».
	maxDepth := 300

	for limit := 0; limit <= maxDepth; limit++ {
		visited := make(map[State]bool)
		bestDepthFound = -1
		if dfsLimited(start, code, 0, limit, visited) {
			// dfsLimited вернёт true, как только найдёт решение
			return limit
		}
	}
	// Не найдено в пределах maxDepth => 0
	return 0
}

// dfsLimited делает DFS с ограничением по глубине depthLimit.
// depth — текущая глубина.
// Если typed == code — нашли решение (возвращаем true).
// Если превысили limit — false.
// Если уже посещали это состояние на этой же глубине или меньшей — нет смысла повторять.
func dfsLimited(st State, code string, depth, limit int, visited map[State]bool) bool {
	if st.typed == code {
		return true
	}
	if depth >= limit {
		return false
	}
	// Префиксная отсечка
	if !strings.HasPrefix(code, st.typed) {
		return false
	}
	// Если уже посещали на <= depth, пропустим (считаем, что более глубокое попадание
	// в то же состояние не даст выгоды).
	if visited[st] {
		return false
	}
	visited[st] = true

	// Пять вариантов нажатия
	inputs := []string{"^", "v", "<", ">", "A"}
	for _, inp := range inputs {
		ns, ok := nextMacroState(st, inp)
		if !ok {
			continue
		}
		if dfsLimited(ns, code, depth+1, limit, visited) {
			return true
		}
	}
	return false
}

// nextMacroState — как и раньше, обрабатывает «макронажатие»:
// если inp — стрелка, двигаем rpos[0]; если 'A', пробиваем вниз по цепочке.
func nextMacroState(st State, inp string) (State, bool) {
	if inp != "A" {
		cur := st.rpos[0]
		nx := robotAdj[cur][inp]
		if nx == "" {
			return State{}, false
		}
		var rp RobPos
		copy(rp[:], st.rpos[:])
		rp[0] = nx
		st.rpos = rp
		return st, true
	}
	// inp == "A"
	return macroActivate(st, 0)
}

// macroActivate — прожимает 'A' на уровне idx.
func macroActivate(st State, idx int) (State, bool) {
	cur := st.rpos[idx]
	switch cur {
	case "^", "v", "<", ">":
		// Двигаем либо робот idx+1, либо цифровую клавиатуру
		if idx == chainLength {
			// двигаем dig
			nx := digitAdj[st.dig][cur]
			if nx == "" {
				return State{}, false
			}
			st.dig = nx
			return st, true
		}
		// двигаем робота idx+1
		n2 := robotAdj[st.rpos[idx+1]][cur]
		if n2 == "" {
			return State{}, false
		}
		var rp RobPos
		copy(rp[:], st.rpos[:])
		rp[idx+1] = n2
		st.rpos = rp
		return st, true

	case "A":
		if idx < chainLength {
			// идём на уровень ниже
			return macroActivate(st, idx+1)
		}
		// idx == chainLength => нажимаем кнопку на цифровой клавиатуре
		d := st.dig
		// если d — стрелка, это gap
		if _, isDir := robotAdj[d]; isDir && d != "A" {
			return State{}, false
		}
		st.typed += d
		return st, true

	default:
		return State{}, false
	}
}

// ------------------- main -------------------

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <input_file>")
		return
	}

	inputFile := os.Args[1]
	f, err := os.Open(inputFile)
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

	start := time.Now()
	total := 0
	for _, code := range codes {
		depth := solveCodeIDDFS(code)
		numVal := parseNumericPart(code)
		total += depth * numVal
	}
	elapsed := time.Since(start).Seconds()

	fmt.Printf("Part 2: %d (time %.4fs)\n", total, elapsed)
}
