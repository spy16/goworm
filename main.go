package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spy16/goworm/brain"
)

var (
	addr     = flag.String("addr", ":8081", "HTTP Address for control API")
	model    = flag.String("model", "c_elegans.json", "Model file path")
	interval = flag.Duration("tick", 50*time.Millisecond, "Simulation step interval")
)

func main() {
	flag.Parse()

	var r io.Reader = os.Stdin
	if *model != "" {
		f, err := os.Open(*model)
		if err != nil {
			fatalExit("failed to open model file: %v", err)
		}
		defer func() { _ = f.Close() }()
		r = f
	}

	var b brain.PointModel
	if err := json.NewDecoder(r).Decode(&b); err != nil {
		fatalExit("failed to load model: %v\n", err)
	}

	wrm := &Worm{
		Brain:    &b,
		Interval: *interval,
	}

	go func() {
		log.Printf("starting server on %s...", *addr)
		if err := http.ListenAndServe(*addr, &b); err != nil {
			log.Printf("server exited: %v", err)
		}
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
