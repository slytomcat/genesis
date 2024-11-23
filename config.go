package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

const cfgFileName = "settings.yaml"

type Env struct {
	Factors         int     // factors count
	Factor_vol      int     // factor initial volume
	Inc             float64 // factor incretion
	Delta           float64 // factor deviation
	Capacity        int     // env capacity limit
	Over_cap_factor float64 // over capacity penalty
}

type Pop struct {
	Init_size      int     // initial population size
	Chromosomes    int     // number of chromosomes
	Gens           int     // initial chromosome length = env factor initial volume
	Match_factor   float64 // env match factor
	Bern_p         float64 // bern probability
	Child_factor   float64 // capacity factor incretion for children
	Mutation_p     float64 // mutation probability
	Mutation_delta int     // mutation size
	Ferity_age     int     // ferity age
	Age_factor     float64 // age factor
}

type Sym struct {
	Ages int
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
