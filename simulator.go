package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/spy16/goworm/brain"
)

type simulator struct {
	Interval  time.Duration
	Brain     *brain.PointModel
	ModelName string
}

func (sim *simulator) Run() {
	cfg := pixelgl.WindowConfig{
		Title:  fmt.Sprintf("GoWorm : %s", sim.ModelName),
		Bounds: pixel.R(0, 0, 600, 400),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	tick := time.NewTicker(sim.Interval)
	defer tick.Stop()

	for range tick.C {
		if win.Closed() {
			break
		}
		win.Clear(colornames.Azure)

		switch {
		case win.Pressed(pixelgl.KeyN):
			sim.poke("nose")

		case win.Pressed(pixelgl.KeyF):
			sim.poke("food")

		case win.Pressed(pixelgl.KeyA):
			sim.poke("anterior")

		case win.Pressed(pixelgl.KeyP):
			sim.poke("posterior")
		}

		spiked := sim.Brain.Step()
		log.Printf("%s", strings.Join(spiked, ","))

		win.Update()
	}
}

func (sim *simulator) poke(cellGroup string) {
	var cells []string

	switch cellGroup {
	case "nose":
		cells = []string{
			"FLPR", "FLPL",
			"ASHL", "ASHR",
			"IL1VL", "IL1VR",
			"OLQDL", "OLQDR",
			"OLQVR", "OLQVL",
		}

	case "food":
		cells = []string{
			"ADFL", "ADFR",
			"ASGR", "ASGL",
			"ASIL", "ASIR",
			"ASJR", "ASJL",
			"AWCL", "AWCR",
			"AWAL", "AWAR",
		}

	case "anterior":
		cells = []string{"FLPL", "FLPR", "BDUL", "BDUR", "SDQR"}

	case "posterior":
		cells = []string{"PVDL", "PVDR", "PVCL", "PVCR"}

	default:
		log.Printf("warn: invalid cell group '%s'", cellGroup)
	}

	for _, cellName := range cells {
		sim.Brain.Cell(cellName).Fire()
	}
}
