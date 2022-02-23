// Package brain provides spiking neural network models to be used by the
// GoWorm simulations.
package brain

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// PointModel represents a network of point-neurons. Zero value is ready
// for use. Needs external synchronisation for concurrent use.
type PointModel struct {
	Groups map[string][]string

	name      string
	cells     []Neuron
	cellIds   map[string]int
	conns     Graph
	synapses  int
	threshold float64
}

// Step re-computes state of all the cells and propagates spikes if any.
func (pm *PointModel) Step() []string {
	var spikedCells []string
	for i := range pm.cells {
		cell := &pm.cells[i]
		cell.Step(func() {
			for _, target := range pm.conns.Neighbors(cell.id) {
				pm.cells[target].nextState += pm.conns.Weight(cell.id, target)
			}
			spikedCells = append(spikedCells, cell.name)
		})
	}
	return spikedCells
}

// Cell allocates a neuron in the model with given name and returns a pointer
// to the cell. If the name already exists, returns the pointer to the existing
// cell.
func (pm *PointModel) Cell(name string) *Neuron {
	if idx, found := pm.cellIds[name]; found {
		return &pm.cells[idx]
	}

	if pm.cellIds == nil {
		pm.cellIds = map[string]int{}
	}

	cell := Neuron{id: len(pm.cells), name: name, threshold: pm.threshold}
	pm.cells = append(pm.cells, cell)
	pm.cellIds[name] = cell.id
	return &pm.cells[cell.id]
}

// Poke forces a pre-defined group of cells to fire.
func (pm *PointModel) Poke(group string) []string {
	cells, found := pm.Groups[group]
	if !found {
		return nil
	}
	for _, cell := range cells {
		pm.Cell(cell).Fire()
	}
	return cells
}

func (pm *PointModel) String() string {
	return fmt.Sprintf("PointModel{size=%d, synapses=%d}", len(pm.cells), pm.synapses)
}

func (pm *PointModel) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/poke":
		q := req.URL.Query().Get("group")
		_ = json.NewEncoder(wr).Encode(pm.Poke(q))

	default:
		wr.WriteHeader(http.StatusNotFound)
	}
}

func (pm *PointModel) MarshalJSON() ([]byte, error) {
	var m jsonPointModel

	m.Name = pm.name
	m.Groups = pm.Groups
	m.Cells = make([]jsonCellModel, len(pm.cells), len(pm.cells))

	for i, cell := range pm.cells {
		targets := pm.conns.Neighbors(cell.id)

		jm := jsonCellModel{
			Name:     cell.name,
			Synapses: make([]string, len(targets), len(targets)),
		}

		for synapseIdx, postID := range targets {
			jm.Synapses[synapseIdx] = fmt.Sprintf("%s,%f", pm.cells[postID].name, pm.conns.Weight(cell.id, postID))
		}

		m.Cells[i] = jm
	}

	return json.Marshal(m)
}

func (pm *PointModel) UnmarshalJSON(bytes []byte) error {
	var m jsonPointModel
	if err := json.Unmarshal(bytes, &m); err != nil {
		return err
	}

	if m.Threshold == 0 {
		m.Threshold = 1
	}

	pm.name = strings.TrimSpace(m.Name)
	pm.Groups = m.Groups
	pm.threshold = m.Threshold

	for _, cm := range m.Cells {
		cell := pm.Cell(cm.Name)
		for _, syn := range cm.Synapses {
			parts := strings.Split(syn, ",")
			if len(parts) != 2 {
				return fmt.Errorf("invalid synapse '%s': must be in 'post,weight' format", syn)
			}

			weight, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return fmt.Errorf("invalid weight in synapse '%s': %v", syn, err)
			}

			pm.synapses++
			pm.conns.ReWeight(
				cell.id,
				pm.Cell(parts[0]).id,
				weight,
				true,
			)
		}
	}

	return nil
}

type jsonPointModel struct {
	Name      string              `json:"name"`
	Cells     []jsonCellModel     `json:"cells"`
	Groups    map[string][]string `json:"groups"`
	Threshold float64             `json:"threshold"`
}

type jsonCellModel struct {
	Name     string   `json:"name"`
	Synapses []string `json:"synapses"`
}
