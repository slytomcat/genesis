package main

import (
	"math"
	"math/rand"
	"runtime"
)

const chSize = 2024

type Population struct {
	Creatures   []*Creature
	Gens        int
	Mutate      func(int) int
	BernP       float64
	ChildFactor float64
	FerityAge   int
	MatchFactor float64
	AgeFactor   float64
	metrics     *Metrics
	workCh      chan func()
}

func NewPopulation(cp *Pop, ageFactor float64, m *Metrics) *Population {
	p := Population{
		Creatures: make([]*Creature, cp.Init_size),
		Gens:      cp.Chromosomes,
		Mutate: func(g int) int {
			if cp.Mutation_p > rand.Float64() {
				v := 1 + rand.ExpFloat64()/12.5*float64(cp.Mutation_delta)
				if v > float64(g)-1 {
					v = float64(g - 1)
				}
				if rand.Intn(2) == 0 {
					v = -v
				}
				return g + int(math.Round(v))
			}
			return g
		},
		MatchFactor: cp.Match_factor,
		BernP:       cp.Bern_p,
		ChildFactor: cp.Child_factor,
		FerityAge:   cp.Ferity_age,
		AgeFactor:   ageFactor,
		metrics:     m,
		workCh:      make(chan func(), chSize),
	}
	for i := range cp.Init_size {
		p.Creatures[i] = NewCreature(p.Mutate, cp.Chromosomes, cp.Gens, cp.Mutation_p)
	}
	for range runtime.NumCPU() / 2 {
		go func() {
			for f := range p.workCh {
				f()
			}
		}()
	}
	return &p
}

type workResults struct {
	crs      []*Creature
	deathAge *int
	bern     int
}

func (p *Population) Next(e *Environment) {
	res := make(chan workResults, chSize)
	go func() {
		capacityFactor := e.CapacityFactor(p.Size())
		for _, c := range p.Creatures {
			p.workCh <- func() {
				wr := workResults{crs: make([]*Creature, 0, 2)}
				dead, child := c.Year(e, p, capacityFactor)
				if !dead {
					wr.crs = append(wr.crs, c)
				} else {
					deathAge := c.age
					wr.deathAge = &deathAge
				}
				if child != nil {
					wr.crs = append(wr.crs, child)
					wr.bern = 1
				}
				res <- wr
			}
		}
	}()
	newP := []*Creature{}
	deathCount := 0
	bernCount := 0
	for range p.Creatures {
		wr := <-res
		newP = append(newP, wr.crs...)
		if wr.deathAge != nil {
			p.metrics.dAgeStore(*wr.deathAge)
			deathCount++
		}
		bernCount += wr.bern
	}
	p.Creatures = newP
	p.metrics.PopStore(p.Size(), bernCount, deathCount)
}

func (p Population) RandomPartner() *Creature {
	return p.Creatures[rand.Intn(len(p.Creatures))]
}

func (p *Population) Size() int {
	return len(p.Creatures)
}
