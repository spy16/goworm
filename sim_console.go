// +build !gui

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spy16/goworm/brain"
	"github.com/spy16/goworm/worm"
)

type simulator struct {
	Interval  time.Duration
	Brain     *brain.PointModel
	ModelName string
}

func (sim *simulator) Run() {
	go func() {
		addr := strings.TrimSpace(os.Getenv("GOWORM_ADDR"))
		if addr == "" {
			addr = ":8081"
		}
		log.Fatalf("server exited : %v", http.ListenAndServe(addr, sim))
	}()

	wrm := worm.Worm{
		Brain:    sim.Brain,
		Interval: sim.Interval,
	}

	if err := wrm.Live(context.Background()); err != nil {
		log.Fatalf("worm died due to illness: %v", err)
	}
	log.Printf("worm died")
}

func (sim *simulator) poke(cellGroup string) bool {
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
		return false
	}

	for _, cellName := range cells {
		sim.Brain.Cell(cellName).Fire()
	}
	return true
}

func (sim *simulator) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	cellGroup := strings.TrimSpace(req.URL.Query().Get("cell_group"))
	if !sim.poke(cellGroup) {
		wr.WriteHeader(http.StatusBadRequest)
		_, _ = wr.Write([]byte(fmt.Sprintf("invalid cell group: %s", cellGroup)))
	}
}
