package brain

import "fmt"

// Neuron represents a neuron in the network.
type Neuron struct {
	// neuron attributes.
	id        int
	name      string
	threshold float64

	// runtime states.
	spiked    bool
	state     float64
	nextState float64
}

// Step updates the state of the cell and returns true if a spike occurred.
func (nu *Neuron) Step(onSpike func()) bool {
	nu.spiked = false // reset spike status

	if nu.state >= nu.threshold {
		onSpike()
		nu.nextState = 0
		nu.state = 0
		nu.spiked = true
		return true
	}

	nu.state += nu.nextState
	nu.nextState = 0
	return false

}

// Fire forces the cell to spike.
func (nu *Neuron) Fire() { nu.state = nu.threshold }

func (nu Neuron) String() string {
	return fmt.Sprintf("Neuron{id='%d'}", nu.id)
}
