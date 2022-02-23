package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spy16/goworm/brain"
)

var (
	addr     = flag.String("addr", ":8081", "HTTP Address for control API")
	model    = flag.String("model", "c_elegans.json", "Model file path")
	save     = flag.String("save", "", "Save the model at exit to the file (no save if empty)")
	interval = flag.Duration("tick", 50*time.Millisecond, "Simulation step interval")
)

func main() {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	b, err := loadModel(*model)
	if err != nil {
		fatalExit("failed to load model: %v", err)
	}

	wrm := &Worm{
		Brain:    b,
		Interval: *interval,
	}

	if err := wrm.Live(ctx); err != nil {
		fatalExit("worm died due to illness: %v", err)
	}

	if *save != "" {
		if err := saveModel(*save, b); err != nil {
			log.Printf("failed to save model: %v", err)
		}
	}
	log.Println("worm lived and died happy")
}

func loadModel(path string) (*brain.PointModel, error) {
	var r io.Reader = os.Stdin
	if path != "" && path != "-" {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer func() { _ = f.Close() }()
		r = f
	}

	var b brain.PointModel
	if err := json.NewDecoder(r).Decode(&b); err != nil {
		return nil, err
	}
	return &b, nil
}

func saveModel(outPath string, model *brain.PointModel) error {
	var w io.Writer = os.Stdout
	if outPath != "" && outPath != "-" {
		f, err := os.OpenFile(outPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
		if err != nil {
			return err
		}
		defer func() { _ = f.Close() }()
		w = f
	}
	return json.NewEncoder(w).Encode(model)
}

func fatalExit(msg string, args ...interface{}) {
	fmt.Printf(msg, args...)
	os.Exit(1)
}
