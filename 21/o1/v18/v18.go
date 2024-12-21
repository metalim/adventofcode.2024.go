// Advent of Code 2024, Day 21: "Keypad Conundrum"
// ------------------------------------------------
// Ошибка: "declared and not used: c" означает, что переменная c объявлена, но нигде не используется.
// Нужно либо задействовать её, либо удалить.
//
// Ниже исправленный пример кода, в котором удалена неиспользуемая переменная.
// Вся логика остаётся прежней. Если вы где-то объявляете `c`, но не используете —
// уберите эту объявленную переменную, чтобы не было ошибки компиляции.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Пример демонстрационной логики (упрощённый):
// - chainLengthPart1 = 2 (Part 1)
// - chainLengthPart2 = 25 (Part 2)
// - BFS с макроактивацией
// - Вывод расчётов для каждой части

const (
	chainLengthPart1 = 2
	chainLengthPart2 = 25
)

// Раскладки:
var robotAdj = map[string]map[string]string{
	"^": {"^": "", "v": "v", "<": "", ">": "A"},
	"A": {"^": "", "v": ">", "<": "^", ">": ""},
	"<": {"^": "", "v": "", "<": "", ">": "v"},
	"v": {"^": "^", "v": "", "<": "<", ">": ">"},
	">": {"^": "A", "v": "", "<": "v", ">": ""},
}

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

// RobPos — фиксированный массив для позиций рук роботов.
// Размер 26, чтобы хватило и для 25, и для 2.
type RobPos [26]string

// State описывает текущее состояние.
type State struct {
	rpos  RobPos
	dig   string
	typed string
}

// initState создаёт начальное состояние: все роботы на 'A', digital = 'A', typed = "".
func initState() State {
	var st State
	for i := 0; i < 26; i++ {
		st.rpos[i] = "A"
	}
	st.dig = "A"
	st.typed = ""
	return st
}

func parseNumericPart(code string) int {
	s := strings.TrimSuffix(code, "A")
	s = strings.TrimLeft(s, "0")
	if s == "" {
		return 0
	}
	n, _ := strconv.Atoi(s)
	return n
}

// solveCode выполняет BFS, чтобы найти длину кратчайшей последовательности
// для набора кода code при chainLen роботах.
func solveCode(chainLen int, code string) int {
	if code == "" {
		return 0
	}

	startSt := initState()
	dist := map[State]int{startSt: 0}
	queue := []State{startSt}

	inputs := []string{"^", "v", "<", ">", "A"}

	for len(queue) > 0 {
		st := queue[0]
		queue = queue[1:]
		d := dist[st]

		if st.typed == code {
			return d
		}
		// Префиксная отсечка
		if !strings.HasPrefix(code, st.typed) {
			continue
		}

		for _, inp := range inputs {
			ns, ok := nextMacroState(st, inp, chainLen)
			if !ok {
				continue
			}
			if _, used := dist[ns]; !used {
				dist[ns] = d + 1
				queue = append(queue, ns)
			}
		}
	}
	return 0
}

// nextMacroState обрабатывает одно нажатие inp (стрелка или 'A') на верхнем роботе (rpos[0]).
func nextMacroState(st State, inp string, chainLen int) (State, bool) {
	if inp != "A" {
		cur := st.rpos[0]
		nx := robotAdj[cur][inp]
		if nx == "" {
			return State{}, false
		}
		st.rpos[0] = nx
		return st, true
	}
	// Иначе 'A'
	return macroActivate(st, 0, chainLen)
}

// macroActivate рекурсивно «пробивает» 'A' вниз по цепочке.
func macroActivate(st State, idx, chainLen int) (State, bool) {
	cur := st.rpos[idx]
	switch cur {
	case "^", "v", "<", ">":
		// Передаём стрелку уровню idx+1 или цифровой
		if idx == chainLen {
			nx := digitAdj[st.dig][cur]
			if nx == "" {
				return State{}, false
			}
			st.dig = nx
			return st, true
		} else {
			n2 := robotAdj[st.rpos[idx+1]][cur]
			if n2 == "" {
				return State{}, false
			}
			st.rpos[idx+1] = n2
			return st, true
		}
	case "A":
		if idx < chainLen {
			return macroActivate(st, idx+1, chainLen)
		}
		// Иначе нажать кнопку цифровой
		d := st.dig
		if _, isDir := robotAdj[d]; isDir && d != "A" {
			return State{}, false
		}
		st.typed += d
		return st, true
	default:
		return State{}, false
	}
}

// main читает 5 кодов из входного файла, решает для Part1 и Part2, печатает результат.
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
		t := strings.TrimSpace(sc.Text())
		if t != "" {
			codes = append(codes, t)
		}
	}
	if err := sc.Err(); err != nil {
		panic(err)
	}
	if len(codes) == 0 {
		fmt.Println("No codes in input.")
		return
	}

	// Part 1
	startP1 := time.Now()
	var sumP1 int
	lengthsP1 := make([]int, len(codes))
	numValsP1 := make([]int, len(codes))

	for i, cd := range codes {
		nv := parseNumericPart(cd)
		ln := solveCode(chainLengthPart1, cd)
		sumP1 += nv * ln
		lengthsP1[i] = ln
		numValsP1[i] = nv
	}
	tP1 := time.Since(startP1)

	// Вывод подробностей + общий итог Part 1
	for i := range codes {
		fmt.Printf("%d * %d = %d\n", numValsP1[i], lengthsP1[i], numValsP1[i]*lengthsP1[i])
	}
	fmt.Printf("Part 1: %d (time %.4fs)\n\n", sumP1, tP1.Seconds())

	// Part 2
	startP2 := time.Now()
	var sumP2 int
	lengthsP2 := make([]int, len(codes))
	numValsP2 := make([]int, len(codes))

	for i, cd := range codes {
		nv := parseNumericPart(cd)
		ln := solveCode(chainLengthPart2, cd)
		sumP2 += nv * ln
		lengthsP2[i] = ln
		numValsP2[i] = nv
	}
	tP2 := time.Since(startP2)

	for i := range codes {
		fmt.Printf("%d * %d = %d\n", numValsP2[i], lengthsP2[i], numValsP2[i]*lengthsP2[i])
	}
	fmt.Printf("Part 2: %d (time %.4fs)\n", sumP2, tP2.Seconds())
}
