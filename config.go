package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

const cfgFileName = "settings.yaml"

type Env struct {
	Factors       int     `yaml:"factors"`         // factors count
	FactorVol     int     `yaml:"factor_vol"`      // factor initial volume
	Inc           float64 `yaml:"inc"`             // factor incretion
	Delta         float64 `yaml:"delta"`           // factor deviation
	Capacity      int     `yaml:"capacity"`        // env capacity limit
	OverCapFactor float64 `yaml:"over_cap_factor"` // over capacity penalty
}

type Pop struct {
	InitSize       int     `yaml:"init_size"`      // initial population size
	Chromosomes    int     `yaml:"chromosomes"`    // number of chromosomes
	Gens           int     `yaml:"gens"`           // initial chromosome length = env factor initial volume
	MatchFactor    float64 `yaml:"match_factor"`   // env match factor
	BirthP         float64 `yaml:"birth_p"`        // birth probability
	ChildFactor    float64 `yaml:"child_factor"`   // capacity factor incretion for children
	MutationP      float64 `yaml:"mutation_p"`     // mutation probability
	Mutation_delta int     `yaml:"mutation_delta"` // mutation size
	FertilityAge   int     `yaml:"fertility_age"`  // ferity age
	AgeFactor      float64 `yaml:"age_factor"`     // age factor
}

type Sym struct {
	Years int
}

type Cfg struct {
	Environment Env
	Population  Pop
	Simulation  Sym
}

func readConfig() (*Cfg, error) {
	data, err := os.ReadFile(cfgFileName)
	if err != nil {
		return nil, err
	}
	cfg := Cfg{}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
