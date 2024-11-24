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
}

func NewRandEnvironment(ce *Env) *Environment {
	e := Environment{
		Factors:       make([]float64, ce.Factors),
		Capacity:      ce.Capacity,
		OverCapFactor: ce.OverCapFactor,
		inc:           make([]float64, ce.Factors),
		delta:         ce.Delta,
	}
	e.Next = e.next
	for i := range e.Factors {
		e.inc[i] = ce.Inc + ce.Inc/4*float64(i)
		e.Factors[i] = float64(ce.FactorVol)
	}
	return &e
}

func (e *Environment) CapacityFactor(pSize int) float64 {
	return 1 + e.OverCapFactor*math.Pow(math.Exp(float64(pSize-e.Capacity)/float64(e.Capacity)), 10)
}

func (e *Environment) next() {
	for i, f := range e.Factors {
		e.Factors[i] = f + rand.Float64()*e.inc[i] + math.Pow(rand.Float64()*2*e.delta-e.delta, 3)
	}
}

func (e *Environment) Match(c *Creature) float64 {
	res := 0.0
	for _, f := range e.Factors {
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

func (e *Environment) factorsList() string {
	r := ""
	for _, f := range e.Factors {
		r += fmt.Sprintf("%4.3f, ", f)
	}
	return r[:len(r)-2]
}

func (e *Environment) MakeAndStore(fileName string, simAges int) error {
	factors := len(e.Factors)
	XYss := make([]plotter.XYs, factors)
	for i := range factors {
		XYss[i] = make(plotter.XYs, simAges)
	}
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	for year := range simAges {
		e.next()
		_, err := f.Write([]byte(e.factorsList() + "\n"))
		if err != nil {
			return err
		}
		for i, f := range e.Factors {
			XYss[i][year].X = float64(year)
			XYss[i][year].Y = f
		}
	}
	for i := range len(e.Factors) {
		MakeAndSaveHistogram(fmt.Sprintf("Environment_%d", i), fmt.Sprintf("Environment factor #%d", i), "age", "value", &XYss[i])
	}
	sMax := e.Capacity + e.Capacity/5
	xys := make(plotter.XYs, sMax)
	for s := range sMax {
		xys[s].X = float64(s)
		xys[s].Y = e.CapacityFactor(s)
	}
	MakeAndSaveHistogram("Capacity_factor", "Capacity factor", "population size", "value", &xys)
	return nil
}

func NewStoredEnvironment(fileName string, capacity int, overCapFactor float64) (*Environment, error) {
	stored, err := readCsv(fileName)
	if err != nil {
		return nil, err
	}
	i := 0
	e := Environment{
		Capacity:      capacity,
		OverCapFactor: overCapFactor,
	}
	e.Next = func() {
		e.Factors = (*stored)[i]
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
	stored := [][]float64{}
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
