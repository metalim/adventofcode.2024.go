package main

import (
	"flag"
	"fmt"
	"maps"
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
solving progress:

part 2 is reading comprehension lol
and graph reading...
or just bruteforcing :)

first sample — just demo of ops

second sample — some logic? or sum? no fucking idea

third sample — AND of 2 numbers, and 2 pairs are swapped
swaps for third sample should include z00,z01,z02,z05
but it's AND, not SUM, so we don't care about third sample

input: 4 pairs swapped

DFS? maybe?
3:25:00 in...
naaah...
just DFS the fuck of it?
what to test, though...

candidates from brute force:
{"cmf", "z26"},
{"vpm", "z36"},
{"gsd", "z26"},
{"bbb", "vpm"},
{"dfp", "z26"},
{"nwm", "z32"},
{"htb", "vpm"},
{"tbt", "z32"},
{"wkk", "z36"},
{"nhb", "psw"},
{"kth", "z12"},
{"kth", "nng"},
{"kth", "psw"},
{"qnf", "vpm"},
{"cnp", "vpm"},

probably incorrect...?
check that later, when solved.
UPDATE: yeah, all correct swaps are in the list, but with lot of false positives.

I'll take a nap, and continue later today.
ok. a nap.

...

Alright. After long nap, then walking 10 kilometers, lunch, dinner, shopping...
I'm back.

For now, **I don't want to visualize the graph** (and solve it manually). I'd like to try automatic solutions.

Things to consider:
- if we continue bruteforcing, we need to improve gate calculation
- as we know it's Sum, we can check if nodes link to lower nodes only. For instance z00 should not link to x01, y02, etc
  - this is a good way to locate the incorrect lanes. And even if the logic is not broken, then we will know problem is in same lane only.

lets isolate the laneWires: groups where wires can be swapped

all swaps except lane z26 do not break the set of inputs
z26 is obviously swapped with wire linked directly to x26 and y26
what next?

now we can find groups related to each lane.
and check if they can be swapped or not.

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

hmmm, z11 lane has swap with different lane?
or 2 swaps in same lane???? naaahh... Eric doesn't do that.
yet, it is possible.

my assumption is all swaps are local, no swaps between z26 and z11 for example.
z11 inputs are limited to x/y11 and below (confirmed?),
yet, swap in same lane doesn't work.
so, it's swap with lane 12?

0-7: x00-x80
8,9,10,11: 01,02,04,08

promising... very... wow...
up to 0-30 fixed with 3 swaps.
31-45 should contain a SINGLE swap

why is it so much slower?

Currently I have solution with 5 swaps. But it should be 4.
I suspect double swaps are never used.
So, instead of checking double swaps, I need to check single swap for next lane.

wow! :)
got solution with 4 swaps!
but it's not accepted LMAO!
cmf,kth,psw,qnf,tbt,vpm,z26,z32 <---
WTF????

so, there're are multiple solutions... that's... unexpected.
wtf Eric?

There's a chance that my solution is incorrect, but tests are passing.
But lets just test all combinations: there are just 2x2x1x1 = 4 of them.
First swap, lane 13: [kth psw] or [kth z12]
Second swap, lane 26: [cmf z26] or [gsd z26]
Third swap, lane 33: [tbt z32]
Fourth swap, lane 37: [qnf vpm]

gsd,kth,qnf,tbt,vpm,z12,z26,z32 was accepted.

Now, why first input wasn't accepted?
I'm not testing multiple bit overflows. Like 11+11=110.
That could be the case. Other than that, I'm not sure.

Confirmed, some tests are not perfect. Might revisit this later.

wow. That was a lot of work. And a fun task.

Now, just 4 hours until next task... zZzZz...
*/

type Lane struct {
	Wires map[string]bool
	Valid bool
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
	xs := parsed.Xs
	zs := parsed.Zs

	var swaps []string // the result
	// z00: [x00, y00, ...], z01: [x01, y01, ... (but not x00, y00)], ...
	lanes := make(map[string]*Lane)
	var minTest int

	var maxLaneTested int
	for iLane := range xs {
		// we are iterating over xs, because zs has 1 more lane, which is not explicitly tested.
		// last lane (z45) doesn't have direct inputs, and is tested in previous loop (z44)
		// via 1+1=10 overflow
		z := zs[iLane]

		maxLaneTested = iLane

		wires := getWires(gates, z)
		// exclude wires of previous VALID lanes
		for j := 0; j < iLane; j++ {
			lane := lanes[zs[j]]
			if !lane.Valid {
				continue
			}
			for w := range wires {
				if lane.Wires[w] {
					delete(wires, w)
				}
			}
		}
		lane := &Lane{Wires: wires}
		// TODO: this is not updated after swaps, but it's ok unless two consecutive lanes need swaps
		lanes[z] = lane
		Verbosef("lane %s: %v\n", z, toSlice(wires))

		// all lanes below iMin should be correct
		// first, test if current lane is also correct
		// then do fuckery with groups (if needed)
		err := testWire(parsed, minTest, iLane)
		if err == 0 {
			// good lane, confirmed!
			lane.Valid = true
			minTest = iLane
			continue
		}

		fmt.Printf("%s has error: %f\n", z, err)

		// now, get the group of wires to swap with, and fix the lane.
		group := maps.Clone(wires)
		// iMin is for wire testing, don't use it for grouping
		for j := 0; j < iLane; j++ {
			z2 := zs[j]
			l2 := lanes[z2]
			if l2.Valid {
				continue
			}
			maps.Copy(group, l2.Wires)
		}
		fmt.Printf("group: %v\n", toSlice(group))
		pairs := getPairs(group)
		var candidates [][2]string
		for _, pair := range pairs {
			err := testWireSwap(parsed, minTest, iLane, pair)
			if err != 0 {
				continue
			}
			fmt.Printf("found swap candidate: %v\n", pair)
			candidates = append(candidates, pair)
		}
		if len(candidates) == 0 {
			fmt.Printf("%s: swap not found :(\n", z)
			continue
		}
		fmt.Printf("found %d swap candidates\n", len(candidates))
		selected := candidates[len(candidates)-1] // TODO: multiple candidates? choose the last one for now
		swaps = append(swaps, selected[:]...)
		fmt.Printf("%s: swapping: %v\n", z, selected)
		gates[selected[0]], gates[selected[1]] = gates[selected[1]], gates[selected[0]]
		lane.Valid = true
		minTest = iLane
		// TODO: another thing to implement: update lane wires after swaps,
		// so next lane will be tested with updated wires
	}

	err := testWire(parsed, 0, maxLaneTested)
	if err == 0 {
		fmt.Printf("%s: full test passed!!!\n", zs[maxLaneTested])
	} else {
		fmt.Printf("%s: error: %f\n", zs[maxLaneTested], err)
	}

	slices.Sort(swaps)
	fmt.Printf("Part 2: %s\t\tin %v\n", strings.Join(swaps, ","), time.Since(timeStart))
}

func getPairs(gates map[string]bool) [][2]string {
	pairs := [][2]string{}

	for g1 := range gates {
		for g2 := range gates {
			if g1 < g2 {
				pairs = append(pairs, [2]string{g1, g2})
			}
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

// get non-leaf (non-input) wires
func getWires(gates Gates, name string) map[string]bool {
	wires := make(map[string]bool)
	if gate, ok := gates[name]; ok {
		// it IS a gate
		wires[name] = true
		for _, wire := range gate.Inputs {
			maps.Copy(wires, getWires(gates, wire))
		}
		return wires
	}
	return wires
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

func hasCommonWire(a, b [2]string) bool {
	for _, wire := range a {
		if slices.Contains(b[:], wire) {
			return true
		}
	}
	return false
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

// returns error from 0 to 1. 0 is no error, 1 is 100% error
func testWireSwap(parsed *Parsed, iMin, iMax int, swaps ...[2]string) float64 {
	gates := parsed.Gates
	zs := parsed.Zs
	for _, swap := range swaps {
		a, b := swap[0], swap[1]
		Verbosef("swapping %s and %s\n", a, b)
		gates[a], gates[b] = gates[b], gates[a]
		defer func() {
			gates[a], gates[b] = gates[b], gates[a]
		}()
	}
	if hasLoops(gates, zs) {
		return 1
	}
	return testWire(parsed, iMin, iMax)
}

func op(a, b int) int {
	return a + b
}

func opSample(a, b int) int {
	return a & b
}

func testWire(parsed *Parsed, iMin, iMax int) (err float64) {
	gates := parsed.Gates
	zs := parsed.Zs
	tests := [][2]int{} // actual values, not bit indices
	for i := iMin; i <= iMax; i++ {
		for a := 0; a <= 1; a++ {
			for b := 0; b <= 1; b++ {
				tests = append(tests, [2]int{a << i, b << i})
			}
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
