package main

import (
	"math/rand"
)

type Creature struct {
	age         int
	chromosomes []int
}

func NewCreature(mutation func(int) int, chromosomes, gens int, mutP float64) *Creature {
	c := Creature{chromosomes: make([]int, chromosomes), age: rand.Intn(100)}
	for i := range chromosomes {
		if mutP > rand.Float64() {
			c.chromosomes[i] = mutation(gens)
		}
	}
	return &c
}

func (c *Creature) Year(e *Environment, p *Population, capacityFactor float64) (bool, *Creature) {
	c.age++
	var child *Creature
	if c.age > p.FerityAge && p.BernP > rand.Float64() {
		child = c.Bern(p)
	}
	if c.age < p.FerityAge { // increase capacity factor for yang creatures
		capacityFactor += p.ChildFactor * capacityFactor * float64((p.FerityAge-c.age)/p.FerityAge)
	}
	deadP := p.MatchFactor*e.Match(c)*capacityFactor + p.AgeFactor*float64(c.age)
	return deadP > rand.Float64(), child
}

func (c *Creature) Bern(p *Population) *Creature {
	partner := p.RandomPartner()
	child := Creature{chromosomes: []int{}}
	for i, g := range c.chromosomes {
		child.chromosomes = append(child.chromosomes, p.Mutate([]int{g, partner.chromosomes[i]}[rand.Intn(2)]))
	}
	return &child
}
