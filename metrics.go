package main

import (
	"fmt"
	"math"
	"slices"

	"gonum.org/v1/plot/plotter"
)

type Metrics struct {
	name      string
	ageCount  map[int]int
	popCount  map[int]int
	dAgeCount map[int]int
	popHist   []int
	bornHist  []int
	deathHist []int
}

var (
	metricsM = &Metrics{name: "Mortal_Population", popCount: map[int]int{}, dAgeCount: map[int]int{}, ageCount: map[int]int{}}
	metricsI = &Metrics{name: "Immortal_Population", popCount: map[int]int{}, dAgeCount: map[int]int{}, ageCount: map[int]int{}}
)

func (m *Metrics) AgeStore(age int) {
	m.ageCount[age] += 1
}

func (m *Metrics) PopStore(size, born, dead int) {
	m.popCount[size] += 1
	m.popHist = append(m.popHist, size)
	m.bornHist = append(m.bornHist, born)
	m.deathHist = append(m.deathHist, dead)
}

func (m *Metrics) dAgeStore(age int) {
	m.dAgeCount[age] += 1
}

func minMidMaxTotal(m map[int]int) (int, float64, int, int) {
	min, max := math.MaxInt64, 0
	sum, count := 0, 0
	for k, v := range m {
		if k > max {
			max = k
		}
		if k < min {
			min = k
		}
		sum += k * v
		count += v
	}
	return min, float64(sum) / float64(count), max, count
}

// minMidMax returns min, avg and max values from slice d
func minMidMax(d []int) (int, float64, int) {
	min, max := math.MaxInt64, 0
	sum := 0.0
	for _, v := range d {
		if v > max {
			max = v
		} else if v < min {
			min = v
		}
		sum += float64(v)
	}
	return min, sum / float64(len(d)), max
}

// xyVals makes *plotter.XYs from m[x]y source.
// Count of x values can be reduced by xDiv. The y-values are collected on new x values.
// the resulted y values can be normalized by yDiv (y=y/yDiv)
func xyVals(m map[int]int, xDiv int, yDiv float64) *plotter.XYs {
	n := len(m)
	short := map[int]float64{}
	counts := map[int]int{}
	for x, y := range m {
		x = (x / xDiv) * xDiv
		v, _ := short[x]
		short[x] = v + float64(y)
		counts[x] += 1
	}
	for x, y := range short {
		short[x] = y / float64(counts[x])
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

func HistXY(d []int, xDiv, yDiv int) *plotter.XYs {
	n := len(d)
	xys := make(map[int]int, n)
	for i, v := range d {
		xys[i] = v
	}
	return xyVals(xys, xDiv, float64(yDiv))
}

func (m *Metrics) Store() {
	min, avg, max := minMidMax(m.popHist)
	years := len(m.popHist)
	xDiv := years / 300
	xDesc := fmt.Sprintf("min/avg/max population size by ages: %v / %v / %v", min, avg, max)
	MakeAndSaveHistogram(m.name+"_size", "Population size by ages", xDesc, "size", HistXY(m.popHist, xDiv, 1))
	min, avg, max = minMidMax(m.bornHist)
	xDesc = fmt.Sprintf("min/avg/max born by ages: %v / %v / %v", min, avg, max)
	MakeAndSaveHistogram(m.name+"_born", "Born by ages", xDesc, "size", HistXY(m.bornHist, xDiv, 1))
	min, avg, max = minMidMax(m.deathHist)
	xDesc = fmt.Sprintf("min/avg/max death by ages: %v / %v / %v", min, avg, max)
	MakeAndSaveHistogram(m.name+"_dead", "Death by ages", xDesc, "size", HistXY(m.deathHist, xDiv, 1))
	min, avg, max, _ = minMidMaxTotal(m.popCount)
	xDesc = fmt.Sprintf("min/avg/max population size: %v / %v / %v", min, avg, max)
	MakeAndSaveHistogram(m.name+"_sizes", "Count of years by population size", xDesc, "years", xyVals(m.popCount, len(m.popCount)/50, 1))
	min, avg, max, total := minMidMaxTotal(m.dAgeCount)
	xDesc = fmt.Sprintf("min/avg/max death age: %v (p= %1.4f) / %v / %v\ntotal creatures: %v", min, float64(m.dAgeCount[min])/float64(total)*100, avg, max, total)
	MakeAndSaveHistogram(m.name+"_deaths", "Relative deaths by age", xDesc, "percent", xyVals(m.dAgeCount, 1, float64(total)/100))
	min, _, max, total = minMidMaxTotal(m.ageCount)
	xDesc = fmt.Sprintf("min/max age:  %v (%2.3f%%) / %v \navg creatures per year: %v", min, float64(m.ageCount[min])/float64(total)*100, max, float64(total)/float64(years))
	MakeAndSaveHistogram(m.name+"_ages", "Headcount by age", xDesc, "percent", xyVals(m.ageCount, 1, float64(total)/100))
	dProbability := map[int]float64{}
	xs := make([]int, 0, len(m.dAgeCount))
	for age, count := range m.dAgeCount {
		hCount := m.ageCount[age]
		if hCount == 0 {
			continue
		}
		dProbability[age] = float64(count) / float64(hCount) * 100
		xs = append(xs, age)
	}
	slices.Sort(xs)
	xys := make(plotter.XYs, len(xs))
	for i, x := range xs {
		xys[i].X = float64(x)
		xys[i].Y = dProbability[x]
	}
	xDesc = fmt.Sprintf("P(%d) = %2.3f%%, P(%d) = %2.3f%%, P(%d) = %2.3f%%, P(%d) = %2.3f%%", xs[0], dProbability[xs[0]], xs[5], dProbability[xs[5]], xs[20], dProbability[xs[20]], xs[60], dProbability[xs[60]])
	MakeAndSaveHistogram(m.name+"_death_probability", "Death probability by age", xDesc, "percent", &xys)
}
