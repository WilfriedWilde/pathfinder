package pathfinding

import (
	"time"
)

// GetRoutesSets returns a set of distinct, non-overlapping routes for each flow increment
//
// Gets passed a pathfinding algorithm as an argument (Bellman-Ford, Edmonds-Karp or Dijkstra+Johnson potentials)
func (f *FlowNetwork) GetRoutesSets(start, end int, pathfinder func(f *FlowNetwork, start, end int) ([]int, bool)) ([][][]string, int, time.Duration) {
	fCopy := f.CopyFlowNetwork()
	now := time.Now()

	routesSets := make([][][]string, 0)
	flow := 0

	for {
		usedEdges, found := pathfinder(fCopy, start, end)
		if !found {
			break
		}

		for station := end; station != start; station = usedEdges[station] {
			prev := usedEdges[station]
			for _, edge := range fCopy.Connections[prev] {
				if edge.To == station && edge.Usable {
					edge.Usable = false
					fCopy.Connections[station][edge.Rev].Usable = true

					break
				}
			}
		}

		flow++
		routesSets = append(routesSets, fCopy.ExtractPaths(flow, start, end))
	}

	if flow == 0 {
		return nil, 0, 0
	}

	return routesSets, flow, time.Since(now)
}

// ExtractPaths extracts the paths that have been used in the residual graph it gets as argument
//
// It does not permanently modify the graph
func (f *FlowNetwork) ExtractPaths(flow, start, end int) [][]string {
	paths := make([][]string, 0)
	path := make([]string, 0)
	modifiedIndexes := make([]int, 0)

	for range flow {

		curr := start
		path = []string{}

		for curr != end {
			for ind, edge := range f.Connections[curr] {

				if edge.Original && !edge.Usable {

					if curr == start {
						edge.Usable = true
						modifiedIndexes = append(modifiedIndexes, ind)
					}
					if curr%2 == 0 {
						path = append(path, f.IDToName[curr])
					}

					curr = edge.To
					break
				}
			}
		}

		path = append(path, f.IDToName[end])
		paths = append(paths, path)
	}

	// undo graph modifications
	for _, edgeInd := range modifiedIndexes {
		f.Connections[start][edgeInd].Usable = false
	}

	return paths
}
