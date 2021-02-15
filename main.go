package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/faiface/pixel/pixelgl"

	"github.com/spy16/goworm/brain"
)

var (
	model     = flag.String("model", "c_elegans.csv", "Model CSV file path")
	interval  = flag.Duration("tick", 50*time.Millisecond, "Simulation step interval")
	threshold = flag.Float64("threshold", 30, "Firing threshold for all the cells")
)

func main() {
	flag.Parse()

	b, err := brain.Load(*model, *threshold)
	if err != nil {
		fatalExit("failed to load model: %v\n", err)
	}

	sim := &simulator{
		Brain:     b,
		Interval:  *interval,
		ModelName: *model,
	}
	pixelgl.Run(sim.Run)
}

func fatalExit(msg string, args ...interface{}) {
	fmt.Printf(msg, args...)
	os.Exit(1)
}
