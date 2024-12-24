package main

import (
	"cmp"
	"flag"
	"fmt"
	"maps"
	"math"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

func catch(err error) {
	if err != nil {
		panic(err)
	}
}

var Verbose bool

func main() {
	flag.BoolVar(&Verbose, "v", false, "verbose output")
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("Usage: go run . input.txt")
		os.Exit(1)
	}

	bs, err := os.ReadFile(flag.Arg(0))
	catch(err)

	parsed := parseInput(string(bs))
	part1(parsed)
	part2(parsed)
}

type WireVal int
type Inputs map[string]WireVal
type GateOp []string
type Gates map[string]GateOp

type Parsed struct {
	Inputs Inputs // actually input wires
	Gates  Gates
	Xs     []string
	Ys     []string
	Zs     []string
}

var reWire = regexp.MustCompile(`^(\w+): (\d+)$`)
var reGate = regexp.MustCompile(`^(\w+) (\w+) (\w+) -> (\w+)$`)

func parseInput(input string) Parsed {
	parts := strings.Split(input, "\n\n")

	p := Parsed{
		Inputs: make(Inputs),
		Gates:  make(Gates),
	}

	for _, line := range strings.Split(parts[0], "\n") {
		m := reWire.FindStringSubmatch(line)
		name := m[1]
		sValue := m[2]
		value, err := strconv.Atoi(sValue)
		catch(err)

		if _, ok := p.Inputs[name]; ok {
			panic("oops, double wire in input file")
		}
		p.Inputs[name] = WireVal(value)
	}

	for _, line := range strings.Split(parts[1], "\n") {
		if len(line) == 0 {
			continue
		}
		m := reGate.FindStringSubmatch(line)
		wire1 := m[1]
		op := m[2]
		wire2 := m[3]
		out := m[4]
		if _, ok := p.Gates[out]; ok {
			panic("oops, double gate in input file")
		}
		p.Gates[out] = GateOp{wire1, op, wire2}
	}

	for i := 0; ; i++ {
		x := nameWire("x", i)
		if _, ok := p.Inputs[x]; !ok {
			break
		}
		p.Xs = append(p.Xs, x)
	}
	for i := 0; ; i++ {
		y := nameWire("y", i)
		if _, ok := p.Inputs[y]; !ok {
			break
		}
		p.Ys = append(p.Ys, y)
	}
	for i := 0; ; i++ {
		z := nameWire("z", i)
		if _, ok := p.Gates[z]; !ok {
			break
		}
		p.Zs = append(p.Zs, z)
	}
	slices.Sort(p.Xs)
	slices.Sort(p.Ys)
	slices.Sort(p.Zs)
	Verbosef("xs: %v\n", p.Xs)
	Verbosef("ys: %v\n", p.Ys)
	Verbosef("zs: %v\n", p.Zs)
	return p
}
func nameWire(prefix string, i int) string {
	return fmt.Sprintf("%s%02d", prefix, i)
}

func part1(parsed Parsed) {
	timeStart := time.Now()
	gates := parsed.Gates
	wires := maps.Clone(parsed.Inputs)
	if Verbose {
		fmt.Println("Wires:")
		for wire, val := range wires {
			fmt.Println(wire, val)
		}
		fmt.Println("Gates:")
		for wire, expr := range gates {
			fmt.Println(wire, expr)
		}
	}

	// determine their values
	z := getZ(gates, wires, parsed.Zs)
	fmt.Printf("Part 1: %d\t\tin %v\n", z, time.Since(timeStart))
}

func getZ(gates Gates, wires Inputs, zs []string) int {
	var zVal int
	for i, z := range zs {
		val := getVal(gates, wires, z)
		zVal += int(val) << i
	}
	return zVal
}

func getVal(gates Gates, wires Inputs, name string) WireVal {
	if gate, ok := gates[name]; ok {
		switch gate[1] {
		case "AND":
			return getVal(gates, wires, gate[0]) & getVal(gates, wires, gate[2])
		case "OR":
			return getVal(gates, wires, gate[0]) | getVal(gates, wires, gate[2])
		case "XOR":
			return getVal(gates, wires, gate[0]) ^ getVal(gates, wires, gate[2])
		default:
			panic(fmt.Sprintf("unknown op: %s", gate[1]))
		}
	}
	return wires[name]
}

/*
reading comprehension lol
and graph reading...
or just bruteforcing :)

sample: just demo of ops

sample2: some logic? or sum?
no fucking idea

sample3: AND of 2 numbers, and 2 pairs swapped
candidates for sample3.txt should include z00,z01,z02,z05
but it's AND, not SUM, so we don't care about the sample3

input: 4 pairs swapped

DFS? maybe?
3:25 in...
naaah...
just DFS the fuck of it?
what to test, though...

I'll take a nap, and continue later today
*/

type Candidate struct {
	Metric float64
	Pair   [2]string
}

func part2(parsed Parsed) {
	timeStart := time.Now()
	gates := parsed.Gates

	// test the wires just for fun
	initialMetric := testWireSwap(gates, parsed.Xs, parsed.Ys, parsed.Zs, "z00", "z00")
	fmt.Printf("initial incorrect: %f\n", initialMetric)

	Verbose = false
	// do we need to draw a graph? and visually check the wires?
	// or we can do loops and test the wires?
	allPairs := [][2]string{}

	for g1 := range gates {
		for g2 := range gates {
			if g1 >= g2 {
				continue
			}
			allPairs = append(allPairs, [2]string{g1, g2})
		}
	}
	fmt.Printf("pairs: %d\n", len(allPairs))

	var candidates []Candidate
	bestMetric := initialMetric
	// parallelize this?
	for _, p1 := range allPairs {
		metric := testWireSwap(gates, parsed.Xs, parsed.Ys, parsed.Zs, p1[0], p1[1])
		if metric < bestMetric {
			bestMetric = metric
			fmt.Printf("new best incorrect: %f\n", bestMetric)
		}
		if metric < initialMetric {
			fmt.Printf("new candidate: %s\n", p1)
			candidates = append(candidates, Candidate{Metric: metric, Pair: p1})
			slices.SortFunc(candidates, func(a, b Candidate) int {
				return cmp.Compare(a.Metric, b.Metric)
			})
			fmt.Printf("candidates: %v\n", candidates)
		}
	}
	slices.SortFunc(candidates, func(a, b Candidate) int {
		return cmp.Compare(a.Metric, b.Metric)
	})
	fmt.Printf("candidates: %v\n", candidates)

	fmt.Printf("Part 2: %d\t\tin %v\n", 0, time.Since(timeStart))
}

func sign(a float64) int {
	if a < 0 {
		return -1
	}
	if a > 0 {
		return 1
	}
	return 0
}

func testLoop(gates Gates, visited map[string]bool, name string) bool {
	if visited[name] {
		return true
	}
	if gate, ok := gates[name]; ok {
		visited[name] = true
		if testLoop(gates, visited, gate[0]) {
			return true
		}
		if testLoop(gates, visited, gate[2]) {
			return true
		}
		visited[name] = false
	}
	return false
}

func hasLoops(gates Gates, zs []string) bool {
	for _, z := range zs {
		if testLoop(gates, make(map[string]bool), z) {
			return true
		}
	}
	return false
}

func testWireSwap(gates Gates, xs, ys, zs []string, a, b string) float64 {
	Verbosef("swapping %s and %s\n", a, b)
	gates[a], gates[b] = gates[b], gates[a]
	defer func() {
		gates[a], gates[b] = gates[b], gates[a]
	}()
	if hasLoops(gates, zs) {
		return math.MaxInt
	}

	// TODO: Need better tests here, with some bit patterns, overflows, etc
	// instead of single bit patterns
	tests := [][2]int{}
	for i := range xs {
		tests = append(tests, [2]int{i, -1})
		tests = append(tests, [2]int{-1, i})
		tests = append(tests, [2]int{i, i})
	}

	var incorrect int
	total := len(tests)
	wires := make(Inputs)
	for _, test := range tests {
		var xVal, yVal int
		if test[0] >= 0 {
			wires[xs[test[0]]] = 1
			xVal = 1 << test[0]
		}
		if test[1] >= 0 {
			wires[ys[test[1]]] = 1
			yVal = 1 << test[1]
		}
		zVal := getZ(gates, wires, zs)
		if zVal != xVal+yVal {
			incorrect++
			Verbosef("x: %x, y: %x, z: %x\n", xVal, yVal, zVal)
		}
		if test[0] >= 0 {
			wires[xs[test[0]]] = 0
		}
		if test[1] >= 0 {
			wires[ys[test[1]]] = 0
		}
	}
	return float64(incorrect) / float64(total)
}

func Verbosef(format string, a ...any) {
	if Verbose {
		fmt.Printf(format, a...)
	}
}

// probably incorrect...?
// hardcodedPairs := [][2]string{
// 	{"cmf", "z26"},
// 	{"vpm", "z36"},
// 	{"gsd", "z26"},
// 	{"bbb", "vpm"},
// 	{"dfp", "z26"},
// 	{"nwm", "z32"},
// 	{"htb", "vpm"},
// 	{"tbt", "z32"},
// 	{"wkk", "z36"},
// 	{"nhb", "psw"},
// 	{"kth", "z12"},
// 	{"kth", "nng"},
// 	{"kth", "psw"},
// 	{"qnf", "vpm"},
// 	{"cnp", "vpm"},
// }
