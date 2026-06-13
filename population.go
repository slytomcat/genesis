package main

import (
	"cmp"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"slices"
)

const chSize = 2048

type Population struct {
	Creatures    []*Creature
	Chromosomes  int
	Mutate       func(int) int
	BirthP       float64
	ChildFactor  float64
	FertilityAge int
	MatchFactor  float64
	AgeFactor    float64
	metrics      *Metrics
	workCh       chan func()
}

func NewPopulation(cp *Pop, ageFactor float64, m *Metrics) *Population {
	p := Population{
		Creatures:    make([]*Creature, cp.InitSize),
		Chromosomes:  cp.Chromosomes,
		Mutate:       MutateFunc(cp.MutationP, float64(cp.Mutation_delta)),
		MatchFactor:  cp.MatchFactor,
		BirthP:       cp.BirthP,
		ChildFactor:  cp.ChildFactor,
		FertilityAge: cp.FertilityAge,
		AgeFactor:    ageFactor,
		metrics:      m,
		workCh:       make(chan func(), chSize), // workers pool input channel
	}
	for i := range cp.InitSize {
		p.Creatures[i] = NewCreature(p.Mutate, cp.Chromosomes, cp.Gens, cp.MutationP)
	}
	cpuCount := runtime.NumCPU()
	if cpuCount > 1 {
		cpuCount-- // leave one cpu free for other processes
	}
	for range cpuCount {
		go func() {
			for f := range p.workCh {
				f()
			}
		}()
	}
	return &p
}

func (p *Population) stop() {
	close(p.workCh)
}

type workResults struct {
	crs      []*Creature
	deathAge *int
	born     int
}

func (p *Population) Next(e *Environment) {
	res := make(chan workResults, chSize)
	go func() {
		envFactor := e.CapacityFactor(p.Size()) * p.MatchFactor
		for _, c := range p.Creatures {
			p.metrics.AgeStore(c.age)
			p.workCh <- func() {
				wr := workResults{crs: make([]*Creature, 0, 2)}
				dead, child := c.Year(e, p, envFactor)
				if !dead {
					wr.crs = append(wr.crs, c)
				} else {
					deathAge := c.age
					wr.deathAge = &deathAge
				}
				if child != nil {
					wr.crs = append(wr.crs, child)
					wr.born = 1
				}
				res <- wr
			}
		}
	}()
	newP := []*Creature{}
	deathCount := 0
	bornCount := 0
	for range p.Creatures {
		wr := <-res
		newP = append(newP, wr.crs...)
		if wr.deathAge != nil {
			p.metrics.dAgeStore(*wr.deathAge)
			deathCount++
		}
		bornCount += wr.born
	}
	p.Creatures = newP
	p.metrics.PopStore(p.Size(), bornCount, deathCount)
}

func (p Population) RandomPartner() *Creature {
	for range 10 {
		partner := p.Creatures[rand.Intn(len(p.Creatures))]
		if partner.age >= p.FertilityAge {
			return partner
		}
	}
	return nil
}

func (p *Population) Size() int {
	return len(p.Creatures)
}

func (p *Population) Dump() {
	file, err := os.Create(p.metrics.name + "_population_dump")
	if err != nil {
		fmt.Printf("dump file creation error: %v", err)
		return
	}
	defer file.Close()
	slices.SortFunc(p.Creatures, func(a, b *Creature) int { return cmp.Compare(a.age, b.age) })
	for _, c := range p.Creatures {
		file.WriteString(c.String() + "\n")
	}
}
