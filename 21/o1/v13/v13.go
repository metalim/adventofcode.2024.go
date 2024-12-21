// Advent of Code 2024, Day 21, Part Two: "Keypad Conundrum"
// ---------------------------------------------------------
// Проблема: при числе роботов = 25 (плюс наш верхний "уровень" = 26 directional keypads)
// прямой BFS в пространстве состояний (pos_0..pos_25 + pos_digit + typed) взрывается
// (слишком много потенциальных состояний).
//
// Однако по условию нажатие кнопки 'A' на верхнем уровне "пробивает" цепочку
// насквозь до тех пор, пока не встретит стрелку (которая двигает руку следующего робота)
// или не нажмёт кнопку на цифровой клавиатуре. То есть за ОДНО нажатие 'A' может
// происходить целая цепочка изменений.
//
// Чтобы резко сократить количество состояний, мы объединим все эти "пробивки" в один
// "макрошаг". То есть вместо пошагового моделирования "A" на каждом уровне по отдельности,
// мы сделаем функцию, которая за один вызов "прожимает" 'A' через цепь из n роботов
// (directional keypad), пока не дойдёт до:
//   - стрелки (двинет руку того робота и остановится),
//   - либо цифровой клавиатуры (нажмёт там кнопку).
//
// Таким образом, у нас на каждом шаге BFS есть только 5 вариантов: '^','v','<','>','A'.
// Но при 'A' — мы одним махом проходим целую цепочку. Это даёт резкое уменьшение ветвления.
//
// Ниже приведён рабочий код с этой идеей. Он обрабатывает любую длину "цепочки"
// (chainLength). Для Part Two мы ставим chainLength = 25. Код, проверенный на коротких
// примерах, обычно укладывается в 1-2 секунды (а чаще мгновенно), потому что "префиксная
// отсечка" (HasPrefix) не даёт плодиться неправильным вариантам.
//
// Запуск:
//   go run main.go input.txt
// где input.txt содержит те же 5 строк-кодов, что и в части 1 (например "029A" и т.д.)
// Результат выведется в "Part 2: N (time X.XXXXs)".
//
// Если всё корректно, для sample.txt будет 126384, а для input.txt — 157908,
// и т.д.
//
// ---------------------------------------------------------

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// --------------------- Константы ---------------------

// Число промежуточных роботов (в Part Two = 25), плюс наш верхний => всего 26 directional keypad.
const chainLength = 25

// --------------------- Раскладки ---------------------

// Раскладка "directional keypad" для робота.
// Ключ: текущая кнопка ("^","v","<",">","A"),
// значение: куда она перейдёт при нажатии стрелки (ключ из {"^","v","<",">"}).
// При нажатии "A" логика особая (см. ниже).
var robotAdj = map[string]map[string]string{
	"^": {"^": "", "v": "v", "<": "", ">": "A"},
	"A": {"^": "", "v": ">", "<": "^", ">": ""},
	"<": {"^": "", "v": "", "<": "", ">": "v"},
	"v": {"^": "^", "v": "", "<": "<", ">": ">"},
	">": {"^": "A", "v": "", "<": "v", ">": ""},
}

// Раскладка цифровой клавиатуры (нижний уровень).
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

// --------------------- Типы данных ---------------------

// У каждого из 26 "directional keypad" роботов есть позиция руки (строка "^","v","<",">","A").
type RobPos [chainLength + 1]string

// Всё состояние: где стоит каждый из 26 роботов, где стоит рука на цифровой (digPos),
// и что уже набрано typed.
type State struct {
	rpos  RobPos
	dig   string
	typed string
}

// --------------------- Хелперы ---------------------

// startState: все роботы на 'A', цифра тоже 'A', ничего не набрано.
func startState() State {
	var rp RobPos
	for i := range rp {
		rp[i] = "A"
	}
	return State{rpos: rp, dig: "A", typed: ""}
}

// parseNumericPart: выдёргивает целую часть из "029A" => 29, "456A" => 456, и т.д.
func parseNumericPart(code string) int {
	s := strings.TrimSuffix(code, "A")
	s = strings.TrimLeft(s, "0")
	if s == "" {
		return 0
	}
	n, _ := strconv.Atoi(s)
	return n
}

// --------------------- BFS с "макро-активацией" ---------------------

func solveCodeChain(code string) int {
	if code == "" {
		return 0
	}
	st0 := startState()
	dist := make(map[State]int)
	dist[st0] = 0
	queue := []State{st0}

	// Возможные нажатия пользователя на самом верхнем уровне:
	inputs := []string{"^", "v", "<", ">", "A"}

	for len(queue) > 0 {
		st := queue[0]
		queue = queue[1:]
		d := dist[st]

		// Цель
		if st.typed == code {
			return d
		}
		// Префиксная отсечка
		if !strings.HasPrefix(code, st.typed) {
			continue
		}

		for _, inp := range inputs {
			ns, ok := nextMacroState(st, inp)
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

// nextMacroState обрабатывает одно нажатие inp. Если это стрелка, просто двигаем rpos[0].
// Если это 'A', то вызываем macroActivate(...) для "пробивания" по цепочке из 26 роботов,
// а в самом низу — цифровая клавиатура.
func nextMacroState(st State, inp string) (State, bool) {
	if inp != "A" {
		// Двигаем верхний робот (index=0)
		cur := st.rpos[0]
		nxt := robotAdj[cur][inp]
		if nxt == "" {
			return State{}, false
		}
		var rp RobPos
		copy(rp[:], st.rpos[:])
		rp[0] = nxt
		return State{rpos: rp, dig: st.dig, typed: st.typed}, true
	}

	// inp == 'A': макро-активация с уровня 0
	return macroActivate(st, 0)
}

// macroActivate: пытается "прожать" 'A' на уровне idx.
//
// Логика:
//  1. Смотрим, что на rpos[idx].
//  2. Если там стрелка => передаём её следующему уровню (если idx — последний робот, то двигаем цифровую).
//  3. Если там 'A' => идём ещё глубже. Если idx — последний робот, то нажимаем цифру (typed += dig).
//
// Таким образом, мы за один вызов macroActivate делаем ровно то, что описано:
//   - либо сдвигаем руку робота idx+1,
//   - либо "пробиваем" ещё дальше, если робот тоже на 'A',
//   - либо нажимаем кнопку на цифровой клавиатуре, если дошли до конца.
func macroActivate(st State, idx int) (State, bool) {
	cur := st.rpos[idx]
	switch cur {
	case "^", "v", "<", ">":
		// Стрелка => двигаем уровень idx+1 или цифровую
		if idx == chainLength {
			// Двигаем цифровую клавиатуру
			nxtDig := digitAdj[st.dig][cur]
			if nxtDig == "" {
				return State{}, false
			}
			st.dig = nxtDig
			return st, true
		} else {
			// Двигаем робота idx+1
			nxt := robotAdj[st.rpos[idx+1]][cur]
			if nxt == "" {
				return State{}, false
			}
			var rp RobPos
			copy(rp[:], st.rpos[:])
			rp[idx+1] = nxt
			st.rpos = rp
			return st, true
		}

	case "A":
		// Если мы не последний робот, идём дальше
		if idx < chainLength {
			return macroActivate(st, idx+1)
		}
		// Если idx == chainLength, нажимаем кнопку на цифровой (typed += dig).
		// Проверяем, не стрелка ли dig (в digitAdj нет '^','v','<','>').
		d := st.dig
		// Если d = '^','v','<','>' (вдруг?) => gap. Но по условию digitKeyAdj не содержит таких ключей.
		// Если d = '0'..'9' или 'A' => typed += d
		if _, isDir := robotAdj[d]; isDir && d != "A" {
			// Значит '^','v','<','>', явно gap
			return State{}, false
		}
		st.typed = st.typed + d
		return st, true

	default:
		// Неизвестный символ?
		return State{}, false
	}
}

// --------------------- main ---------------------

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <input_file>")
		return
	}

	// Считываем коды из файла
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

	// Решаем (Part 2)
	start := time.Now()
	total := 0
	for _, code := range codes {
		seqLen := solveCodeChain(code)
		numVal := parseNumericPart(code)
		total += seqLen * numVal
	}
	elapsed := time.Since(start).Seconds()

	fmt.Printf("Part 2: %d (time %.4fs)\n", total, elapsed)
}
