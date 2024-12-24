package main

import (
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

var Verbose, Verbose1, Verbose2 bool

func main() {
	flag.BoolVar(&Verbose, "v", false, "verbose output")
	flag.BoolVar(&Verbose1, "v1", false, "verbose output for part 1")
	flag.BoolVar(&Verbose2, "v2", false, "verbose output for part 2")
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
type GateOp struct {
	Op     string
	Inputs [2]string
}
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

func parseInput(input string) *Parsed {
	parts := strings.Split(input, "\n\n")

	p := &Parsed{
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
		p.Gates[out] = GateOp{Op: op, Inputs: [2]string{wire1, wire2}}
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

func part1(parsed *Parsed) {
	if Verbose1 {
		old := Verbose
		Verbose = true
		defer func() {
			Verbose = old
		}()
	}
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
		switch gate.Op {
		case "AND":
			return getVal(gates, wires, gate.Inputs[0]) & getVal(gates, wires, gate.Inputs[1])
		case "OR":
			return getVal(gates, wires, gate.Inputs[0]) | getVal(gates, wires, gate.Inputs[1])
		case "XOR":
			return getVal(gates, wires, gate.Inputs[0]) ^ getVal(gates, wires, gate.Inputs[1])
		default:
			panic(fmt.Sprintf("unknown op: %s", gate.Op))
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

I'll take a nap, and continue later today.
ok. a nap.

...

Alright. After long nap, then walk for 10 kilometers, lunch, dinner, shopping...
I'm back.

For now, I don't want to visualize the graph (and solve it manually). I'd like to try automatic solutions.

Things to consider:
- if we continue bruteforcing, we need to improve gate calculation
- as we know it's Sum, we can check if nodes link to lower nodes only. For instance z00 should not link to x01, y02, etc
	- this is a good way to locate the incorrect lanes. And even if the logic is not broken, then we will know problem is in same lane only.
*/

func getPairs(gates map[string]bool) [][2]string {
	pairs := [][2]string{}

	for g1 := range gates {
		for g2 := range gates {
			if g1 >= g2 {
				continue
			}
			pairs = append(pairs, [2]string{g1, g2})
		}
	}
	slices.SortFunc(pairs, func(a, b [2]string) int {
		// compare [0], then [1]
		if a[0] != b[0] {
			return strings.Compare(a[0], b[0])
		}
		return strings.Compare(a[1], b[1])
	})
	return pairs
}

func getLeaves(gates Gates, name string) map[string]bool {
	leaves := make(map[string]bool)
	if gate, ok := gates[name]; ok {
		for _, wire := range gate.Inputs {
			maps.Copy(leaves, getLeaves(gates, wire))
		}
		return leaves
	}
	leaves[name] = true
	return leaves
}

func getWires(gates Gates, name string) map[string]bool {
	wires := make(map[string]bool)
	if gate, ok := gates[name]; ok {
		for _, wire := range gate.Inputs {
			wires[wire] = true
			maps.Copy(wires, getWires(gates, wire))
		}
		return wires
	}
	wires[name] = true
	return wires
}
func toSet(s []string) map[string]bool {
	seen := make(map[string]bool)
	for _, v := range s {
		seen[v] = true
	}
	return seen
}

func toSlice(s map[string]bool) []string {
	result := []string{}
	for k, v := range s {
		if v {
			result = append(result, k)
		}
	}
	slices.Sort(result)
	return result
}

func checkMissingAndExtraLeaves(zs []string, i int, leaves map[string]bool) bool {
	// 1. all x_j and y_j with j = (0..i) should be in leaves
	// except the largest z_i, which has only x_(j-1) and y_(j-1)
	var missingLeaves bool
	maxJ := i
	if i == len(zs)-1 {
		maxJ = i - 1
	}
	for j := 0; j <= maxJ; j++ {
		if !leaves[nameWire("x", j)] && !leaves[nameWire("y", j)] {
			missingLeaves = true
			break
		}
	}
	if missingLeaves {
		fmt.Printf("%s missing some leaves: %v\n", zs[i], toSlice(leaves))
	}

	// 2. all leaves are i or less
	var extraLeaves bool
	for l := range leaves {
		n, err := strconv.Atoi(l[1:])
		catch(err)
		if n > i {
			extraLeaves = true
			break
		}
	}
	if extraLeaves {
		fmt.Printf("%s has extra leaves: %v\n", zs[i], toSlice(leaves))
	}
	return missingLeaves || extraLeaves
}

func part2(parsed *Parsed) {
	if Verbose2 {
		old := Verbose
		Verbose = true
		defer func() {
			Verbose = old
		}()
	}
	timeStart := time.Now()
	gates := parsed.Gates
	zs := parsed.Zs
	/*
		btw, should we use values of inputs from input file for testing?
		x: 100111000100001000000001010000101000101101111 low -> high
		y: 111101010111110000000010010111000100010111001 low -> high
		nah...
	*/

	var swaps []string
	// ok, let isolate the lanes: groups where wires can be swapped
	lanes := make(map[string]map[string]bool) // z00: [x00, y00, ...], z01: [x01, y01, ... (but not x00, y00)], ...
	for i, z := range zs {
		if i > 12 {
			break
		}

		leaves := getLeaves(gates, z)
		if checkMissingAndExtraLeaves(zs, i, leaves) {
			fmt.Printf("lane %s is incorrect\n", z)
			// what to do with z26?
			// ...
		}
		// ok, all swaps except lane z26 do not break the logic
		// z26 is obviously swapped with wire linked directly to x26 and y26
		// what next?

		// now we can find groups related to each i. and check if they are swapped
		// or not.

		wires := getWires(gates, z)
		// exclude wires from previous lanes
		for j := 0; j < i; j++ {
			for w := range wires {
				if lanes[zs[j]][w] {
					delete(wires, w)
					// exclude from the lane (and from the swap group)
				}
			}
		}
		// exclude input wires
		for leaf := range leaves {
			delete(wires, leaf)
		}
		lanes[z] = wires
		Verbosef("lane %s: %v\n", z, toSlice(wires))
		/*
			WOOHOO!
			so there are different kinds of swaps.

			assumptions (based on illegal observations of output lol):
			z00-z11 are potentially correct.
			z12 has some wires swapped (z13 can be correct or incorrect, but it excludes wires from z12)
			z26 obviously (z27 can be correct or incorrect, but it excludes wires from z26)
			...
			in total minimum 2 swaps. But we need 4.
			next?

			go i=0..max, and test each bit.
			minimum there are 4 checks: (0,0)=0, (0,1)=1, (1,0)=1, (1,1)=10

			good. we are splitting the problem into smaller parts.
			now lets test lanes, with upper limit. So lower lanes will be tested first, and confirmed to be correct.
			then we can test higher lanes, and find the incorrect ones.
			Dam! the solution will be almost instant lol! (if swaps are limited to same lane)

			(I'm drinking tea with cookies... Tasty!)
		*/

		// all lower lanes should be correct
		// first test if this lane is also correct
		// then do fuckery with groups (if needed)
		err := testWire(parsed, i)
		if err == 0 {
			// good lane, confirmed!
			continue
		}

		fmt.Printf("lane %s has error: %f\n", z, err)

		// now, get the group of wires to swap with, and fix the lane.
		group := maps.Clone(wires)
		group[z] = true
		Verbosef("group: %v\n", toSlice(group))
		pairs := getPairs(group)
		for _, pair := range pairs {
			err := testWireSwap(parsed, i, pair[0], pair[1])
			if err == 0 {
				swaps = append(swaps, pair[:]...)
				fmt.Printf("swap %s and %s works\n", pair[0], pair[1])
			}
		}
		/*
			hmmm, z11 lane has swap with different lane?
			or 2 swaps in same lane???? naaahh... Eric doesn't do that.
			yet, it is possible.

			my assumption is all swaps are local, no swaps between z26 and z11 for example.
			z11 inputs are limited to x/y11 and below (confirmed?),
			yet, swap in same lane doesn't work.
			so, it's swap with lane 12?

			taking a break to think.
		*/

	} // for each lane

	slices.Sort(swaps)
	fmt.Printf("Part 2: %s\t\tin %v\n", strings.Join(swaps, ","), time.Since(timeStart))
}

func testLoop(gates Gates, visited map[string]bool, name string) bool {
	if visited[name] {
		return true
	}
	if gate, ok := gates[name]; ok {
		visited[name] = true
		for _, input := range gate.Inputs {
			if testLoop(gates, visited, input) {
				return true
			}
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

func testWireSwap(parsed *Parsed, i int, a, b string) float64 {
	gates := parsed.Gates
	zs := parsed.Zs
	Verbosef("swapping %s and %s\n", a, b)
	gates[a], gates[b] = gates[b], gates[a]
	defer func() {
		gates[a], gates[b] = gates[b], gates[a]
	}()
	if hasLoops(gates, zs) {
		return math.MaxInt
	}
	return testWire(parsed, i)
}

func op(a, b int) int {
	return a + b
}

func opSample(a, b int) int {
	return a & b
}

func testWire(parsed *Parsed, i int) (err float64) {
	gates := parsed.Gates
	zs := parsed.Zs
	tests := [][2]int{} // actual values, not bit indices
	for a := 0; a <= 1; a++ {
		for b := 0; b <= 1; b++ {
			tests = append(tests, [2]int{a << i, b << i})
		}
	}

	var incorrect int
	total := len(tests)
	wires := make(Inputs)
	for _, test := range tests {
		xVal := test[0]
		yVal := test[1]
		setInput(wires, xVal, yVal, parsed)
		zVal := getZ(gates, wires, zs)
		if zVal != op(xVal, yVal) {
			incorrect++
			Verbosef("x: %x, y: %x, z: %x\n", xVal, yVal, zVal)
		}
	}
	return float64(incorrect) / float64(total)
}

func setInput(wires Inputs, xVal, yVal int, parsed *Parsed) {
	xs := parsed.Xs
	ys := parsed.Ys
	for i := range xs {
		if xVal&(1<<i) != 0 {
			wires[xs[i]] = 1
		} else {
			wires[xs[i]] = 0
		}
	}
	for i := range ys { // just for kicks
		if yVal&(1<<i) != 0 {
			wires[ys[i]] = 1
		} else {
			wires[ys[i]] = 0
		}
	}
}

func Verbosef(format string, a ...any) {
	if Verbose {
		fmt.Printf(format, a...)
	}
}

// probably incorrect...?
// check that later, when solved.
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
