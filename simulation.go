package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	if len(args) != 2 {
		usage(args[0])
		return
	}
	c, err := readConfig()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	var e *Environment
	switch args[1] {
	case "store":
		e = NewRandEnvironment(&c.Environment)
		err = e.MakeAndStore(envFile, c.Simulation.Years)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		e.SaveHistograms()
		return
	case "random":
		e = NewRandEnvironment(&c.Environment)
		defer e.SaveHistograms()
	case "stored":
		e, err = NewStoredEnvironment(envFile, c.Environment.Capacity, c.Environment.OverCapFactor, c.Simulation.Years)
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
	defer pM.stop()
	pI := NewPopulation(&c.Population, 0, metricsI)
	defer pI.stop()
	for i, c := range pM.Creatures {
		pI.Creatures[i] = c.Copy()
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go InterruptHandler(ctx, cancel)
	diffCum := 0
	started := time.Now()
	year := 0
	defer func() {
		elapsed := time.Since(started)
		fmt.Println()
		results()
		fmt.Printf("Simulation for %d years finished in %v (%.3f years per second)\n", year+1, elapsed, float64(year)/elapsed.Seconds())
		fmt.Printf("Cumulative relative difference %.3f\n", float64(diffCum)/float64(year))
	}()
	for year = range c.Simulation.Years {
		select {
		case <-ctx.Done():
			return
		default:
			e.Next()
			pM.Next(e)
			pI.Next(e)
			diff := pI.Size() - pM.Size()
			fmt.Printf("year %6d\tMortal: %6d\tdiff: %+5d,\tImmortal: %6d\tenvironment: %s\r", year, pM.Size(), diff, pI.Size(), e.factorsList())
			diffCum += diff
			if pM.Size() == 0 || pI.Size() == 0 {
				return
			}
		}
	}
}

func InterruptHandler(ctx context.Context, onFinish func()) {
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
	case <-sig:
		onFinish()
	}
}

func results() {
	metricsM.Store()
	metricsI.Store()
}
