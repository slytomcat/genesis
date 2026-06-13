package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"gonum.org/v1/plot/plotter"
)

type Environment struct {
	Factors       []float64
	Capacity      int
	OverCapFactor float64
	inc           []float64
	delta         float64
	Next          func()
	Stored        [][]float64
}

func NewRandEnvironment(ce *Env) *Environment {
	e := Environment{
		Factors:       make([]float64, ce.Factors),
		Capacity:      ce.Capacity,
		OverCapFactor: ce.OverCapFactor,
		inc:           make([]float64, ce.Factors),
		delta:         ce.Delta,
		Stored:        make([][]float64, 0, 100),
	}
	e.Next = e.next
	for i := range e.Factors {
		e.inc[i] = ce.Inc + ce.Inc/4*(rand.Float64()*2-1)
		e.Factors[i] = float64(ce.FactorVol)
	}
	return &e
}

func (e *Environment) CapacityFactor(pSize int) float64 {
	return 1 + e.OverCapFactor*math.Pow(math.Exp(float64(pSize-e.Capacity)/float64(e.Capacity)), 8)
}

func (e *Environment) next() {
	st := make([]float64, len(e.Factors))
	for i, f := range e.Factors {
		e.Factors[i] = f + rand.Float64()*e.inc[i] + math.Pow(1.3*(rand.Float64()*2*e.delta-e.delta), 3)
		st[i] = e.Factors[i]
	}
	e.Stored = append(e.Stored, st)
}

func (e *Environment) Match(c *Creature) float64 {
	res := 0.0
	for _, f := range e.Factors {
		// find closest chromosome to factor f
		r := math.MaxFloat64
		for _, g := range c.chromosomes {
			if v := math.Abs(float64(g) - f); v < r {
				r = v
			}
		}
		res += r
	}
	return res / float64(len(e.Factors))
}

var r = strings.Builder{}
func (e *Environment) factorsList() string {
	r.Reset()
	for _, f := range e.Factors {
		fmt.Fprintf(&r, "%4.3f, ", f)
	}
	return r.String()[:r.Len()-2]
}

func (e *Environment) MakeAndStore(fileName string, simAges int) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	for range simAges {
		e.next()
		if _, err := f.Write([]byte(e.factorsList() + "\n")); err != nil {
			return err
		}
	}
	return nil
}

func (e *Environment) SaveHistograms() {
	factors := len(e.Factors)
	XYss := make([]plotter.XYs, factors)
	for i := range factors {
		XYss[i] = make(plotter.XYs, len(e.Stored))
	}
	for year, factors := range e.Stored {
		for i, f := range factors {
			XYss[i][year].X = float64(year)
			XYss[i][year].Y = f
		}
	}
	for i := range factors {
		MakeAndSaveHistogram(fmt.Sprintf("Environment_%d", i), fmt.Sprintf("Environment factor #%d", i), "age", "value", &XYss[i])
	}
	sMax := e.Capacity + e.Capacity/5
	xys := make(plotter.XYs, sMax)
	for s := range sMax {
		xys[s].X = float64(s)
		xys[s].Y = e.CapacityFactor(s)
	}
	MakeAndSaveHistogram("Capacity_factor", "Capacity factor", "population size", "value", &xys)
}

func NewStoredEnvironment(fileName string, capacity int, overCapFactor float64, years int) (*Environment, error) {
	stored, err := readCsv(fileName)
	if err != nil {
		return nil, err
	}
	e := Environment{
		Capacity:      capacity,
		OverCapFactor: overCapFactor,
		Stored:        *stored,
	}
	if len(*stored) < years {
		return nil, fmt.Errorf("stored environment has smaller years (%d) than required (%d) make new env via 'store' or decrease simulation years in settings", len(*stored), years)
	}
	i := 0
	e.Next = func() {
		e.Factors = (e.Stored)[i]
		i++
	}
	return &e, nil
}

func readCsv(fileName string) (*[][]float64, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	stored := make([][]float64, 0, len(data))
	for _, s := range data {
		sl := []float64{}
		for _, fs := range s {
			v, err := strconv.ParseFloat(strings.Trim(fs, " "), 64)
			if err != nil {
				return nil, err
			}
			sl = append(sl, v)
		}
		stored = append(stored, sl)
	}
	if len(stored) == 0 {
		return nil, fmt.Errorf("file '%s' contain no data", fileName)
	}
	return &stored, nil
}
