package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spy16/goworm/brain"
	"github.com/spy16/goworm/worm"
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

	wrm := &worm.Worm{
		Brain:    b,
		Interval: *interval,
		Sensory: map[string][]string{
			"nose": {
				"FLPR", "FLPL",
				"ASHL", "ASHR",
				"IL1VL", "IL1VR",
				"OLQDL", "OLQDR",
				"OLQVR", "OLQVL",
			},
			"food": {
				"ADFL", "ADFR",
				"ASGR", "ASGL",
				"ASIL", "ASIR",
				"ASJR", "ASJL",
				"AWCL", "AWCR",
				"AWAL", "AWAR",
			},
			"anterior":  {"FLPL", "FLPR", "BDUL", "BDUR", "SDQR"},
			"posterior": {"PVDL", "PVDR", "PVCL", "PVCR"},
		},
	}

	go func() {
		log.Printf("server exited: %v", http.ListenAndServe(":8081", stimulusHandler(wrm)))
	}()

	if err := wrm.Live(context.Background()); err != nil {
		fatalExit("worm died due to illness: %v", err)
	}
	fmt.Println("worm lived and died happy")
}

func fatalExit(msg string, args ...interface{}) {
	fmt.Printf(msg, args...)
	os.Exit(1)
}

func stimulusHandler(w *worm.Worm) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		group := req.URL.Query().Get("group")
		if err := w.Stimulate(group); err != nil {
			http.Error(wr, err.Error(), http.StatusBadRequest)
			wr.WriteHeader(http.StatusBadRequest)
			return
		}
	})
}
