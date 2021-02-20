package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/spy16/goworm/brain"
)

// Dead is returned when the worm dies unexpectedly.
var Dead = errors.New("dead")

// Worm represents an entire worm lifetime.
type Worm struct {
	Brain    *brain.PointModel
	Action   func(spikes []string)
	Interval time.Duration
}

// Live starts the lifetime of a worm and runs until context is cancelled.
func (w *Worm) Live(ctx context.Context) error {
	if w.Brain == nil {
		return fmt.Errorf("%w: died without a brain", Dead)
	} else if w.Action == nil {
		w.Action = func(spikes []string) {
			log.Printf("spikes: %v", spikes)
		}
	}

	tick := time.NewTicker(w.Interval)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil

		case <-tick.C:
			w.Action(w.Brain.Step())
		}
	}
}
