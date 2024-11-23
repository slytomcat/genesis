package main

import (
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatch(t *testing.T) {
	c := Creature{
		age:         0,
		chromosomes: []int{3, 8},
	}
	e := Environment{
		Factors: []float64{8},
	}
	require.Zero(t, e.Match(&c))
	e.Factors[0] = 3
	require.Zero(t, e.Match(&c))
	e.Factors[0] = 4
	require.Equal(t, 1, int(e.Match(&c)))
	e.Factors[0] = 11
	require.Equal(t, 3, int(e.Match(&c)))
}

func TestCapFactorForChildren(t *testing.T) {
	cf := 10.0
	fa := 8
	for age := range 10 {
		cf := cf
		if age < fa {
			cf += cf * float64(fa-age) / float64(fa)
		}
		fmt.Println(cf)
	}
}

func TestMutation(t *testing.T) {
	count := 5000
	size := 200
	var hist map[int]int
	for m := range 10 {
		p := NewPopulation(&Pop{
			Chromosomes:    1,
			Mutation_p:     1,
			Mutation_delta: m + 1,
		}, 0, nil)
		hist = make(map[int]int)
		for range count {
			n := p.Mutate(size)
			v, _ := hist[n]
			hist[n] = v + 1
		}
		require.LessOrEqual(t, len(hist), (m+2)*2, m+1)
	}
	xs := make([]int, 0, len(hist))
	for k := range hist {
		xs = append(xs, k)
	}
	slices.Sort(xs)
	for _, x := range xs {
		t.Logf("%2d - %.6f (%d)", x-size, float64(hist[x])/float64(count), hist[x])
	}
}

func TestRun(t *testing.T) {
	run([]string{"", "stored"})
}
