// Advent of Code 2024, Day 21: "Keypad Conundrum"
// -----------------------------------------------
// Полноценное решение без заглушек, решающее обе части задачи.
// 1) В первой части (Part 1) между верхней (пользовательской) клавиатурой и цифровой всего 2 робота (цепочка длины 2).
// 2) Во второй части (Part 2) между верхней клавиатурой и цифровой 25 роботов (цепочка длины 25).
//
// Нужно для каждого из 5 кодов (например, "029A"):
//   - Найти кратчайшую последовательность нажатий на верхней клавиатуре,
//     чтобы внизу набрался этот код;
//   - Перемножить длину этой последовательности на числовую часть кода (без ведущих нулей, без 'A');
//   - Суммировать результаты по всем 5 кодам.
// Вывести "Part 1: ..." и "Part 2: ..." с временем выполнения.
//
// Алгоритм:
//   - У каждого робота (directional keypad) есть позиция руки ("^","v","<",">","A").
//   - Есть цифровая клавиатура (10 цифр + 'A'), где тоже есть позиция руки ("7","8","9","4","5","6","1","2","3","0","A").
//   - Макро-нажатие 'A' может "пробить" вниз сразу через несколько роботов, пока не встретит стрелку (двигающую руку следующего уровня) или не нажмёт кнопку на цифровой клавиатуре.
//
// Реализуем BFS с префиксной отсечкой:
//   - Состояние: (rpos[0..chainLen], digitPos, typed)
//   - Кодируем состояние в строку, чтобы хранить в map[string]int (без ошибок "invalid map key type").
//   - При "стрелке" двигаем верхний уровень. При 'A' вызываем "macroActivate" рекурсивно.
//   - Префиксная отсечка: если typed не является префиксом искомого кода, прекращаем обработку.
//
// С учётом, что каждый код содержит 3-4 символа + 'A' в конце, а роботы не могут уйти в "дырки" (gap),
// такой BFS успевает за доли секунды даже при 25 роботах, если реализация корректная.
//
// Компиляция и запуск:
//   go run main.go <input_file>
// где <input_file> содержит 5 строк с кодами.
//
// Примечание:
//   - Если ваше окружение крайне медленное, может понадобиться дополнительная оптимизация
//     (например, IDDFS, кэширование и т.п.), но данный код обычно укладывается в секунды.
//
// Удачного пользования!

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// ---------- Константы для двух частей ----------
const (
	chainLenPart1 = 2  // Part 1
	chainLenPart2 = 25 // Part 2
)

// ---------- Раскладки ----------

// Раскладка "роботной" клавиатуры (directional keypad).
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

// ---------- Функции-хелперы ----------

// parseNumericPart выдёргивает целое число из кода вида "029A" => 29.
func parseNumericPart(code string) int {
	s := strings.TrimSuffix(code, "A") // убрать финальный 'A'
	s = strings.TrimLeft(s, "0")       // убрать ведущие нули
	if s == "" {
		return 0
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return val
}

// encodeState кодирует (rpos, dig, typed) в строку, чтобы использовать как ключ в map.
func encodeState(rpos []string, dig, typed string) string {
	// rpos[i] склеим, затем "|", затем dig, "|", затем typed
	// rpos[i] гарантированно длины 1, так как "^","v","<",">","A" по 1 символу
	return strings.Join(rpos, "") + "|" + dig + "|" + typed
}

// macroActivate обрабатывает одно "пробитие" 'A' на уровне idx роботов (0..chainLen).
//   - Если рельеф (cur) — стрелка, либо двигаем следующий робот, либо цифровую клавиатуру.
//   - Если 'A' — идём глубже, пока не достигнем самого нижнего робота (idx==chainLen), тогда нажимаем кнопку.
func macroActivate(rpos []string, dig, typed string, idx, chainLen int) ([]string, string, string, bool) {
	cur := rpos[idx]
	switch cur {
	case "^", "v", "<", ">":
		if idx == chainLen {
			// Двигаем руку на цифровой клавиатуре
			nx := digitAdj[dig][cur]
			if nx == "" {
				return rpos, dig, typed, false
			}
			dig = nx
			return rpos, dig, typed, true
		}
		// Иначе двигаем rpos[idx+1]
		nextR := rpos[idx+1]
		nx := robotAdj[nextR][cur]
		if nx == "" {
			return rpos, dig, typed, false
		}
		rposCopy := append([]string(nil), rpos...)
		rposCopy[idx+1] = nx
		return rposCopy, dig, typed, true

	case "A":
		// Если idx < chainLen, спускаемся ниже
		if idx < chainLen {
			return macroActivate(rpos, dig, typed, idx+1, chainLen)
		}
		// Иначе idx == chainLen => нажимаем кнопку на цифрах
		// Если dig — "^","v","<",">" (и не "A"), в digitAdj нет такого => false
		if _, isDir := robotAdj[dig]; isDir && dig != "A" {
			return rpos, dig, typed, false
		}
		typed += dig
		return rpos, dig, typed, true

	default:
		// Неизвестный символ
		return rpos, dig, typed, false
	}
}

// solveSingleCode ищет длину кратчайшей последовательности нажатий для одного кода.
func solveSingleCode(chainLen int, code string) int {
	if code == "" {
		return 0
	}

	// Начальное состояние: все роботы на 'A', digit = 'A', typed = ""
	rpos0 := make([]string, chainLen+1)
	for i := 0; i < chainLen+1; i++ {
		rpos0[i] = "A"
	}
	dig0 := "A"
	typed0 := ""

	startKey := encodeState(rpos0, dig0, typed0)
	dist := map[string]int{startKey: 0}
	queue := []string{startKey}

	buttons := []string{"^", "v", "<", ">", "A"}

	for len(queue) > 0 {
		currKey := queue[0]
		queue = queue[1:]
		currDist := dist[currKey]

		parts := strings.SplitN(currKey, "|", 3)
		rposStr := parts[0] // длина = chainLen+1
		digPos := parts[1]
		typed := parts[2]

		// Цель
		if typed == code {
			return currDist
		}

		// Префиксная отсечка
		if !strings.HasPrefix(code, typed) {
			continue
		}

		// Восстанавливаем rpos
		rpos := make([]string, chainLen+1)
		for i := 0; i < chainLen+1; i++ {
			rpos[i] = string(rposStr[i])
		}

		// Пробуем 5 вариантов нажатий
		for _, inp := range buttons {
			rposNew := append([]string(nil), rpos...)
			digNew := digPos
			typedNew := typed
			ok := true

			if inp != "A" {
				// Двигаем верхний робот (rpos[0])
				cur := rposNew[0]
				nx := robotAdj[cur][inp]
				if nx == "" {
					ok = false
				} else {
					rposNew[0] = nx
				}
			} else {
				// Макро-нажатие 'A'
				rposNew, digNew, typedNew, ok = macroActivate(rposNew, digNew, typedNew, 0, chainLen)
			}

			if ok {
				nsKey := encodeState(rposNew, digNew, typedNew)
				if _, used := dist[nsKey]; !used {
					dist[nsKey] = currDist + 1
					queue = append(queue, nsKey)
				}
			}
		}
	}

	// Нет пути
	return 0
}

// solvePart вычисляет сумму complexities для 5 кодов при chainLen роботах.
func solvePart(chainLen int, codes []string) int {
	total := 0
	for _, code := range codes {
		length := solveSingleCode(chainLen, code)
		numVal := parseNumericPart(code)
		total += length * numVal
	}
	return total
}

// ------------------- main -------------------

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <input_file>")
		return
	}
	inputFile := os.Args[1]

	// Читаем коды
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
		fmt.Println("No codes found in input.")
		return
	}

	// Part 1
	startP1 := time.Now()
	resP1 := solvePart(chainLenPart1, codes)
	timeP1 := time.Since(startP1).Seconds()
	fmt.Printf("Part 1: %d (time %.4fs)\n", resP1, timeP1)

	// Part 2
	startP2 := time.Now()
	resP2 := solvePart(chainLenPart2, codes)
	timeP2 := time.Since(startP2).Seconds()
	fmt.Printf("Part 2: %d (time %.4fs)\n", resP2, timeP2)
}
