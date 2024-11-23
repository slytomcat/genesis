package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

const (
	envFile = "env.csv"
	usage   = `Usage:
	%s <cmd>

Commands:
	store - make random environment and store it to file
	random - start simulation with random environment
	stored - start simulation with stored environment
`
)

func main() {
	run(os.Args)
}

func run(args []string) {
	var e *Environment
	var err error
	if len(args) != 2 {
		fmt.Printf(usage, args[0])
		return
	}
	c, err := readConfig()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	switch args[1] {
	case "store":
		e = NewRandEnvironment(&c.Environment)
		err = e.MakeAndStore(envFile, c.Simulation.Ages)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		return
	case "random":
		e = NewRandEnvironment(&c.Environment)
	case "stored":
		e, err = NewStoredEnvironment(envFile, c.Environment.Capacity, c.Environment.Over_cap_factor)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			fmt.Printf(usage, args[0])
			return
		}
	default:
		fmt.Printf("Error: unknown command %s\n", args[1])
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go InterruptHandler(cancel)
	pM := NewPopulation(&c.Population, c.Population.Age_factor, metricsM)
	pI := NewPopulation(&c.Population, 0, metricsI)
	copy(pI.Creatures, pM.Creatures)
	defer results()
	defer close(pM.workCh)
	defer close(pI.workCh)
	for year := range c.Simulation.Ages {
		select {
		case <-ctx.Done():
			return
		default:
		}
		pM.Next(e)
		pI.Next(e)
		e.Next()
		fmt.Printf("year %5d,\tMortal: %6d\tdiff: %+5d,\tImmortal: %6d\tfactor: %s\n", year, pM.Size(), pM.Size()-pI.Size(), pI.Size(), e.factorsList())
		if pM.Size() == 0 || pI.Size() == 0 {
			return
		}
	}
}

func InterruptHandler(cancel func()) {
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	cancel()
}

func results() {
	metricsM.Store()
	metricsI.Store()
}
