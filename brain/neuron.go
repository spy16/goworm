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

// Fire forces the cell to spike.
func (cell *Neuron) Fire() { cell.state = cell.threshold }

// IsFiring returns true if the cell spiked in the last step.
func (cell *Neuron) IsFiring() bool { return cell.spiked }

// Apply applies the given stimulus to the cell.
func (cell *Neuron) Apply(st float64) { cell.nextState += st }

// Step updates the state of the cell and returns true if a spike occurred.
func (cell *Neuron) Step(onSpike func()) bool {
	cell.spiked = false // reset spike status

	if cell.state >= cell.threshold {
		if onSpike != nil {
			onSpike()
		}
		cell.nextState = 0
		cell.state = 0
		cell.spiked = true
		return true
	}

	cell.state += cell.nextState
	cell.nextState = 0
	return false
}

func (cell Neuron) String() string {
	return fmt.Sprintf("Neuron{id='%d', spiking=%t}", cell.id, cell.IsFiring())
}
