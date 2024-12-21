// Advent of Code 2024, Day 21, Part Two: "Keypad Conundrum" — решение с 25 промежуточными роботами.
// ----------------------------------------------------------------------
// Проблема "invalid map key type State" возникает, если в структуре State
// есть поле типа slice (например, []string). Слайсы в Go не сравнимы, и их
// нельзя использовать как ключи в map.
//
// Решение: вместо []string использовать массив фиксированной длины (который
// полностью сравним), ведь число роботов (25) константно. Тогда State
// станет сравнимым, и map[State] будет работать.
//
// В примере ниже:
//   - chainLength = 25 (число промежуточных роботов);
//   - Определяем RobPosArray [chainLength+1]string для позиций рук;
//   - State содержит такой массив, а также поле digPos (цифровая клавиатура) и typed.
//
// Код вычисляет кратчайшие последовательности нажатий для 5 кодов
// (аналогично первой части, но с 25 роботами). Выводим итоговую сумму complexities.
//
// Запуск:
//   go run main.go input.txt

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Число промежуточных роботов
const chainLength = 25

// Массив позиций рук роботов (фикс. размер, чтобы можно было сравнивать State).
type RobPosArray [chainLength + 1]string

// Описание состояния: где рука каждого из (chainLength+1) роботов, где рука на цифровой клавиатуре, что набрано.
type State struct {
	robPos RobPosArray
	digPos string
	typed  string
}

// Раскладка "directional keypad" для роботов
var robotKeyAdj = map[string]map[string]string{
	"^": {
		"^": "",
		"v": "v",
		"<": "",
		">": "A",
	},
	"A": {
		"^": "",
		"v": ">",
		"<": "^",
		">": "",
	},
	"<": {
		"^": "",
		"v": "",
		"<": "",
		">": "v",
	},
	"v": {
		"^": "^",
		"v": "",
		"<": "<",
		">": ">",
	},
	">": {
		"^": "A",
		"v": "",
		"<": "v",
		">": "",
	},
}

// Раскладка цифровой клавиатуры (нижний уровень).
var digitKeyAdj = map[string]map[string]string{
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

// startState создаёт начальное состояние: все роботы на 'A', digital = 'A', typed = "".
func startState() State {
	var rob RobPosArray
	for i := 0; i < chainLength+1; i++ {
		rob[i] = "A"
	}
	return State{robPos: rob, digPos: "A", typed: ""}
}

// solveCodeN делает BFS по состояниям и находит длину кратчайшей последовательности нажатий для code.
func solveCodeN(code string) int {
	if code == "" {
		return 0
	}
	st0 := startState()

	dist := map[State]int{st0: 0}
	queue := []State{st0}

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
			ns, ok := nextState(st, inp)
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

// nextState обрабатывает одно нажатие inp на уровне 0 (верхний).
func nextState(st State, inp string) (State, bool) {
	rob := st.robPos
	dig := st.digPos
	t := st.typed

	if inp != "A" {
		// Двигаем верхний робот (rob[0])
		nextPos := robotKeyAdj[rob[0]][inp]
		if nextPos == "" {
			return State{}, false
		}
		var newRob RobPosArray
		copy(newRob[:], rob[:])
		newRob[0] = nextPos
		return State{robPos: newRob, digPos: dig, typed: t}, true
	}

	// Иначе inp == "A"
	return activateLevel(st, 0)
}

// activateLevel пытается "активировать" робота на уровне idx.
func activateLevel(st State, idx int) (State, bool) {
	r := st.robPos[idx]
	// Если r — стрелка, то отдаём её следующему уровню (или цифровой клавиатуре, если idx — последний).
	if r == "^" || r == "v" || r == "<" || r == ">" {
		if idx == chainLength {
			// Двигаем цифровую клавиатуру
			nextDig := digitKeyAdj[st.digPos][r]
			if nextDig == "" {
				return State{}, false
			}
			ns := st
			ns.digPos = nextDig
			return ns, true
		} else {
			// Двигаем робота idx+1
			r2 := st.robPos[idx+1]
			nextPos := robotKeyAdj[r2][r]
			if nextPos == "" {
				return State{}, false
			}
			ns := st
			var newRob RobPosArray
			copy(newRob[:], st.robPos[:])
			newRob[idx+1] = nextPos
			ns.robPos = newRob
			return ns, true
		}
	}

	// Если r == 'A', спускаемся на уровень idx+1 или нажимаем цифровую кнопку, если idx — последний.
	if r == "A" {
		if idx == chainLength {
			// Нажимаем кнопку на цифровой клавиатуре => typed += digPos
			d := st.digPos
			// Если это стрелка (вдруг?), то gap. Но в digitKeyAdj нет '^','v','<','>' ключей.
			// Если 'A' или цифра => добавляем.
			// Проверим, нет ли в robotKeyAdj[d], кроме 'A' (которая совпадает у обоих).
			if _, found := robotKeyAdj[d]; found && d != "A" {
				// Значит '^','v','<','>' => gap
				return State{}, false
			}
			ns := st
			ns.typed = ns.typed + d
			return ns, true
		} else {
			// Активируем следующий уровень
			return activateLevel(st, idx+1)
		}
	}

	// Иначе непонятная кнопка => false
	return State{}, false
}

// parseNumericPart: "029A" => 29, "456A" => 456, "000A" => 0
func parseNumericPart(code string) int {
	s := strings.TrimSuffix(code, "A")
	s = strings.TrimLeft(s, "0")
	if s == "" {
		return 0
	}
	n, _ := strconv.Atoi(s)
	return n
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <input_file>")
		return
	}

	// Считаем 5 кодов из файла
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

	// Подсчёт общей сложности
	start := time.Now()
	sumComplexities := 0
	for _, code := range codes {
		seqLen := solveCodeN(code)
		numVal := parseNumericPart(code)
		sumComplexities += seqLen * numVal
	}
	elapsed := time.Since(start).Seconds()

	fmt.Printf("Part 2: %d (time %.4fs)\n", sumComplexities, elapsed)
}
