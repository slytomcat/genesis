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
)

var (
	version = "local"
)

func usage(app string) {
	fmt.Printf(`%s v. %s
Usage:
	%s <cmd>

Commands:
	store - make random environment and store it to file
	random - start simulation with random environment
	stored - start simulation with stored environment
	version, -v, --version - prints version and exits
`, app, version, app)
}

func main() {
	run(os.Args)
}

func run(args []string) {
	var e *Environment
	var err error
	if len(args) != 2 {
		usage(args[0])
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
		err = e.MakeAndStore(envFile, c.Simulation.Years)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		return
	case "random":
		e = NewRandEnvironment(&c.Environment)
	case "stored":
		e, err = NewStoredEnvironment(envFile, c.Environment.Capacity, c.Environment.OverCapFactor)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
	case "version", "-v", "--version":
		fmt.Printf("%s v. %s\n", args[0], version)
		return
	default:
		fmt.Printf("Error: unknown command %s\n", args[1])
		usage(args[0])
		return
	}
	pM := NewPopulation(&c.Population, c.Population.AgeFactor, metricsM)
	pI := NewPopulation(&c.Population, 0, metricsI)
	for i, c := range pM.Creatures {
		pI.Creatures[i] = c.Copy()
	}
	defer close(pM.workCh)
	defer close(pI.workCh)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go InterruptHandler(cancel)
	defer results()
	for year := range c.Simulation.Years {
		select {
		case <-ctx.Done():
			return
		default:
			e.Next()
			pM.Next(e)
			pI.Next(e)
			fmt.Printf("year %5d,\tMortal: %6d\tdiff: %+5d,\tImmortal: %6d\tfactor: %s\n", year, pM.Size(), pM.Size()-pI.Size(), pI.Size(), e.factorsList())
			if pM.Size() == 0 || pI.Size() == 0 {
				return
			}
		}
	}
}

func InterruptHandler(onFinish func()) {
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	onFinish()
}

func results() {
	metricsM.Store()
	metricsI.Store()
}
