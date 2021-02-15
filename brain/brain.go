package brain

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

// Load reads PointModel from a CSV file and returns an instance.
func Load(file string, threshold float64) (*PointModel, error) {
	f, err := os.OpenFile(file, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	pm := &PointModel{threshold: threshold}
	if err := pm.fromCSVStream(f); err != nil {
		return nil, err
	}
	return pm, nil
}

// PointModel represents a network of point-neurons.
type PointModel struct {
	threshold float64
	cells     []Neuron
	cellIds   map[string]int
	conns     Graph
	synapses  int
}

// Step re-computes state of all the cells and propagates spikes if any.
func (model *PointModel) Step() []string {
	var spikedCells []string
	for i := range model.cells {
		cell := &model.cells[i]
		cell.Step(func() {
			for _, target := range model.conns.Neighbors(cell.id) {
				model.cells[target].nextState += model.conns.Weight(cell.id, target)
			}
			spikedCells = append(spikedCells, cell.name)
		})
	}
	return spikedCells
}

// Cell allocates a neuron in the model with given name and returns a pointer
// to the cell. If the name already exists, returns the pointer to the existing
// cell.
func (model *PointModel) Cell(name string) *Neuron {
	if idx, found := model.cellIds[name]; found {
		return &model.cells[idx]
	}

	if model.cellIds == nil {
		model.cellIds = map[string]int{}
	}

	cell := Neuron{
		id:        len(model.cells),
		name:      name,
		threshold: model.threshold,
	}
	model.cells = append(model.cells, cell)
	model.cellIds[name] = cell.id
	return &model.cells[cell.id]
}

func (model *PointModel) fromCSVStream(rwc io.Reader) error {
	model.cellIds = map[string]int{}
	model.cells = nil
	model.synapses = 0

	cr := csv.NewReader(rwc)
	cr.Comment = ';'
	cr.FieldsPerRecord = 3
	cr.ReuseRecord = true

	for {
		rec, err := cr.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		pre, post := rec[0], rec[1]
		weight, err := strconv.ParseFloat(rec[2], 64)
		if err != nil {
			return err
		}

		model.synapses++
		model.conns.ReWeight(
			model.Cell(pre).id,
			model.Cell(post).id,
			weight,
			true,
		)
	}

	return nil
}

func (model *PointModel) String() string {
	return fmt.Sprintf("PointModel{size=%d, synapses=%d}", len(model.cells), model.synapses)
}
