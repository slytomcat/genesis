package main

import (
	"fmt"
	"math"
	"slices"

	"gonum.org/v1/plot/plotter"
)

type Metrics struct {
	name      string
	popCount  map[int]int
	dAgeCount map[int]int
	popHist   []int
	bornHist  []int
	deathHist []int
}

var (
	metricsM = &Metrics{name: "Mortal_Population", popCount: map[int]int{}, dAgeCount: map[int]int{}}
	metricsI = &Metrics{name: "Immortal_Population", popCount: map[int]int{}, dAgeCount: map[int]int{}}
)

func (m *Metrics) PopStore(size, bern, dead int) {
	v := m.popCount[size]
	m.popCount[size] = v + 1
	m.popHist = append(m.popHist, size)
	m.bornHist = append(m.bornHist, bern)
	m.deathHist = append(m.deathHist, dead)
}

func (m *Metrics) dAgeStore(age int) {
	v := m.dAgeCount[age]
	m.dAgeCount[age] = v + 1
}

func minMidMaxTotal(m map[int]int) (float64, float64, float64, float64) {
	min, max := math.MaxInt64, 0
	sum, count := 0.0, 0.0
	for k, v := range m {
		if k > max {
			max = k
		} else if k < min {
			min = k
		}
		sum += float64(k * v)
		count += float64(v)
	}
	return float64(min), sum / count, float64(max), count
}

// xyVals makes *plotter.XYs from m[x]y source.
// Count of x values can be reduced by xDiv. The y-values are collected on new x values.
// the resulted y values can be normalized by yDiv (y=y/yDiv)
func xyVals(m map[int]int, xDiv int, yDiv float64) *plotter.XYs {
	n := len(m)
	short := map[int]int{}
	for x, y := range m {
		x = (x / xDiv) * xDiv
		v, _ := short[x]
		short[x] = v + y
	}
	n = len(short)
	xs := make([]int, 0, n)
	for x := range short {
		xs = append(xs, x)
	}
	xys := make(plotter.XYs, n)
	slices.Sort(xs)
	for i, x := range xs {
		xys[i].X = float64(x)
		xys[i].Y = float64(short[x]) / yDiv
	}
	return &xys
}

func HistXY(d []int, xDiv int, yDiv float64) *plotter.XYs {
	n := len(d)
	xys := make(map[int]int, n)
	for i, v := range d {
		xys[i] = v
	}
	return xyVals(xys, xDiv, yDiv)
}

func (m *Metrics) Store() {
	min, avg, max, _ := minMidMaxTotal(m.popCount)
	xDesc := fmt.Sprintf("min/avg/max population size: %v / %v / %v", min, avg, max)
	MakeAndSaveHistogram(m.name, "Population size by age", xDesc, "size", HistXY(m.popHist, 1, 1))
	MakeAndSaveHistogram(m.name+"_bern", "Bern by age", "age", "size", HistXY(m.bornHist, 20, 1))
	MakeAndSaveHistogram(m.name+"_dead", "Dead by age", "age", "size", HistXY(m.deathHist, 20, 1))
	MakeAndSaveHistogram(m.name+"_sizes", "Count of years by population size", xDesc, "years", xyVals(m.popCount, len(m.popCount)/20, 1))
	min, avg, max, total := minMidMaxTotal(m.dAgeCount)
	xDesc = fmt.Sprintf("min/avg/max death age: %v / %v / %v\ntotal creatures: %v", min, avg, max, total)
	MakeAndSaveHistogram(m.name+"_deaths", "Death probability by age", xDesc, "percent", xyVals(m.dAgeCount, 1, total))
}
