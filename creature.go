package main

import (
	"fmt"
	"math"
	"math/rand"
)

type Creature struct {
	age         int
	chromosomes []int
}

func (c *Creature) String() string {
	return fmt.Sprintf(`age: %3d, chromosomes: %v`, c.age, c.chromosomes)
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

func (c *Creature) Year(e *Environment, p *Population, envFactor float64) (bool, *Creature) {
	c.age++
	// check if creature is fertile and can birth a child
	var child *Creature
	if c.age >= p.FertilityAge && p.BirthP > rand.Float64() {
		child = c.Bern(p)
	}
	if c.age < p.FertilityAge { // increase environment factor for yang creatures
		envFactor -= p.ChildFactor * float64((p.FertilityAge-c.age)/p.FertilityAge)
	}
	dead := 0.1+rand.Float64()*0.9 < e.Match(c)*envFactor+p.AgeFactor*float64(c.age)
	return dead, child
}

func (c *Creature) Bern(p *Population) *Creature {
	partner := p.RandomPartner()
	if partner == nil {
		return nil
	}
	child := &Creature{chromosomes: make([]int, len(c.chromosomes))}
	for i, g := range c.chromosomes {
		child.chromosomes[i] = p.Mutate([]int{g, partner.chromosomes[i]}[rand.Intn(2)])
	}
	return child
}

func (c *Creature) Copy() *Creature {
	n := Creature{age: c.age, chromosomes: make([]int, len(c.chromosomes))}
	copy(n.chromosomes, c.chromosomes)
	return &n
}

func MutateFunc(MutationP, Mutation_delta float64) func(int) int {
	return func(g int) int {
		if MutationP > rand.Float64() {
			v := int(math.Round(1 + rand.ExpFloat64()/12.5*Mutation_delta))
			if rand.Float64() > 0.6 {
				v = -v
			}
			g += v
			if g <= 0 {
				return 1
			}
		}
		return g
	}
}
