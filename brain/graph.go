package brain

// Graph represents a weighted, directed, cyclic/acyclic graph.
// Zero value is ready for use. Not safe for concurrent use.
type Graph struct {
	edges map[int]map[int]float64
	size  int
}

// Neighbors returns the neighbors of the node 'n'.
func (g *Graph) Neighbors(n int) []int {
	var result []int
	for neighbor := range g.edges[n] {
		result = append(result, neighbor)
	}
	return result
}

// Weight returns the weight of the edge from n1 to n2. Returns 0
// if the edge does not exist.
func (g *Graph) Weight(n1, n2 int) float64 {
	return g.edges[n1][n2]
}

// ReWeight updates the weight of edge from n1 to n2 and returns
// the updated value.
func (g *Graph) ReWeight(n1, n2 int, delta float64, replace bool) float64 {
	if g.edges == nil {
		g.edges = map[int]map[int]float64{
			n1: {n2: delta},
		}
		g.size = 1
		return delta
	} else if g.edges[n1] == nil {
		g.edges[n1] = map[int]float64{n2: delta}
		return delta
	}

	if replace {
		g.edges[n1][n2] = delta
		return delta
	}

	g.edges[n1][n2] += delta
	return g.edges[n1][n2]
}

// Iter iterates through all the edges of the graph and calls 'fn'
// for each edge. fn can terminate the iterations by returning true.
// Order of iteration may vary across different invocations.
func (g *Graph) Iter(fn EdgeFn) (stopped bool) {
	for from, edges := range g.edges {
		for to, w := range edges {
			if fn(from, to, w) {
				return true
			}
		}
	}
	return false
}

// EdgeFn is invoked for every edge in a graph or until an invocation
// returns stop=true.
type EdgeFn func(from, to int, w float64) (stop bool)
