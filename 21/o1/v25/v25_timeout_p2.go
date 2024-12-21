// Advent of Code 2024, Day 21: "Keypad Conundrum"
// ------------------------------------------------
// Реализация с посимвольным набором кода и BFS для каждого символа,
// без неиспользуемых переменных.
//
// Ключевая идея (коротко):
//  1. Код "XYZ...A" разбиваем на отдельные символы (например, ['0','2','9','A']).
//  2. Для каждого символа запускаем BFS, чтобы нажать его (начиная с текущего положения рук роботов и цифровой).
//  3. После нажатия символа, извлекаем из итогового состояния (ключа) текущее расположение рук (роботов и цифр),
//     чтобы продолжить набор следующего символа.
//  4. Так по одному символу набираем весь код, складывая длины. Умножаем результат на числовую часть кода.
//  5. Повторяем для всех 5 кодов, выводим сумму.
//
// При таком подходе не нужны лишние переменные типа finalRob/finalDig, если мы сразу извлекаем
// итоговое положение из конечного состояния, когда BFS находит нужный символ.
//
// ------------------------------------------------
// Компиляция/запуск:
//   go run main.go sample.txt
//   go run main.go input.txt
//
// Для sample.txt должно выдать Part 1: 126384 и Part 2: 154115708116294,
// при условии нормальной производительности среды.
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

// Part 1: 2 робота, Part 2: 25 роботов
const (
	chainLenPart1 = 2
	chainLenPart2 = 25
)

// Раскладки для роботной и цифровой клавиатур:
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

// parseNumericPart выдёргивает целую часть кода (без ведущих нулей, убирает финальный 'A').
func parseNumericPart(code string) int {
	s := strings.TrimSuffix(code, "A")
	s = strings.TrimLeft(s, "0")
	if s == "" {
		return 0
	}
	n, _ := strconv.Atoi(s)
	return n
}

// macroActivate: прожимаем 'A' на уровне idx.
// Возвращаем: (новоеRob, новоеDigit, символЕслиНажалиКнопку, ок?).
func macroActivate(rob []string, digit string, idx, chainLen int) ([]string, string, string, bool) {
	cur := rob[idx]
	switch cur {
	case "^", "v", "<", ">":
		// Двигаем нижний уровень или цифру
		if idx == chainLen-1 {
			nx := digitAdj[digit][cur]
			if nx == "" {
				return rob, digit, "", false
			}
			return rob, nx, "", true
		}
		n2 := robotAdj[rob[idx+1]][cur]
		if n2 == "" {
			return rob, digit, "", false
		}
		cp := append([]string(nil), rob...)
		cp[idx+1] = n2
		return cp, digit, "", true
	case "A":
		// Идём глубже или нажимаем кнопку
		if idx < chainLen-1 {
			return macroActivate(rob, digit, idx+1, chainLen)
		}
		// idx==chainLen-1 => нажать кнопку digit
		if _, isDir := robotAdj[digit]; isDir && digit != "A" {
			return rob, digit, "", false
		}
		return rob, digit, digit, true
	}
	return rob, digit, "", false
}

// stateOneChar хранит текущее состояние при наборе одного символа.
type stateOneChar struct {
	rob []string
	dig string
	chr string // если chr != "", значит мы уже нажали этот символ
}

func encodeOneChar(s stateOneChar) string {
	return strings.Join(s.rob, "") + "|" + s.dig + "|" + s.chr
}

// bfsOneCharFull: находим кратчайшую последовательность, чтобы нажать targetChar,
// возвращаем (число нажатий, конечное положение rob, конечное положение digit).
func bfsOneCharFull(chainLen int, targetChar string, startRob []string, startDig string) (int, []string, string) {
	start := stateOneChar{
		rob: append([]string(nil), startRob...),
		dig: startDig,
		chr: "",
	}
	dist := map[string]int{encodeOneChar(start): 0}
	queue := []stateOneChar{start}

	buttons := []string{"^", "v", "<", ">", "A"}

	var finalKey string
	found := false

	for len(queue) > 0 && !found {
		st := queue[0]
		queue = queue[1:]
		stk := encodeOneChar(st)
		d := dist[stk]

		// Если нажали нужный символ
		if st.chr == targetChar {
			finalKey = stk
			found = true
			break
		}

		// Генерируем соседей
		for _, inp := range buttons {
			rn := append([]string(nil), st.rob...)
			dn := st.dig
			ch := st.chr
			ok := true

			if inp != "A" {
				cur := rn[0]
				nx := robotAdj[cur][inp]
				if nx == "" {
					ok = false
				} else {
					rn[0] = nx
				}
			} else {
				r2, d2, pressed, ok2 := macroActivate(rn, dn, 0, chainLen)
				if !ok2 {
					ok = false
				} else {
					rn = r2
					dn = d2
					if pressed != "" {
						ch = pressed // символ, который нажали
					}
				}
			}

			if ok {
				nst := stateOneChar{rn, dn, ch}
				nk := encodeOneChar(nst)
				if _, used := dist[nk]; !used {
					dist[nk] = d + 1
					queue = append(queue, nst)
				}
			}
		}
	}

	if !found {
		// Не смогли нажать targetChar
		// Возвращаем 0, начальное положение
		return 0, startRob, startDig
	}

	// Ищем dist
	steps := dist[finalKey]

	// Восстанавливаем конечное положение
	parts := strings.SplitN(finalKey, "|", 3)
	robStr := parts[0]
	digStr := parts[1]
	rpos := make([]string, chainLen)
	for i := 0; i < chainLen; i++ {
		rpos[i] = string(robStr[i])
	}
	return steps, rpos, digStr
}

// solveCodeByCharsFull: разбиваем код на символы, набираем по одному.
func solveCodeByCharsFull(chainLen int, code string) int {
	chs := strings.Split(code, "")
	rob := make([]string, chainLen)
	for i := range rob {
		rob[i] = "A"
	}
	dig := "A"

	totalPresses := 0
	for _, c := range chs {
		steps, newRob, newDig := bfsOneCharFull(chainLen, c, rob, dig)
		totalPresses += steps
		rob = newRob
		dig = newDig
	}
	return totalPresses
}

// solvePartChars: суммируем complexities для 5 кодов.
func solvePartChars(chainLen int, codes []string) int {
	res := 0
	for _, cd := range codes {
		presses := solveCodeByCharsFull(chainLen, cd)
		numVal := parseNumericPart(cd)
		res += presses * numVal
	}
	return res
}

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
		t := strings.TrimSpace(sc.Text())
		if t != "" {
			codes = append(codes, t)
		}
	}
	if err := sc.Err(); err != nil {
		panic(err)
	}
	if len(codes) == 0 {
		fmt.Println("No codes found in input.")
		return
	}

	// Part 1
	startP1 := time.Now()
	resP1 := solvePartChars(chainLenPart1, codes)
	tP1 := time.Since(startP1).Seconds()
	fmt.Printf("Part 1: %d (time %.4fs)\n", resP1, tP1)

	// Part 2
	startP2 := time.Now()
	resP2 := solvePartChars(chainLenPart2, codes)
	tP2 := time.Since(startP2).Seconds()
	fmt.Printf("Part 2: %d (time %.4fs)\n", resP2, tP2)
}
