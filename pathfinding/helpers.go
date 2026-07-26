package pathfinding

import (
	"fmt"
	"slices"
	"time"
)

// GetRoutesSets returns a set of distinct, non-overlapping routes for each flow increment
//
// Gets passed a pathfinding algorithm as an argument (Bellman-Ford, Edmonds-Karp or Dijkstra+Johnson potentials)
func (f *FlowNetwork) GetRoutesSets(start, end, trains int, pathfinder func(f *FlowNetwork, start, end int) ([]int, bool)) ([][][]string, int, time.Duration) {
	fCopy := f.CopyFlowNetwork()
	now := time.Now()

	routesSets := make([][][]string, 0)
	routeSet := make([][]string, 0) // append to routesSets and create a new one whenever reverse edges have been used
	flow := 0

	for {
		reverseEdgeSeen := false
		route := make([]string, 0) // build singular route and append it if no reverse edges
		usedEdges, found := pathfinder(fCopy, start, end)
		if !found {
			break
		}

		for station := end; station != start; station = usedEdges[station] {
			prev := usedEdges[station]
			for _, edge := range fCopy.Connections[prev] {
				if edge.To == station && edge.Usable {
					if !edge.Original { // reverse edge
						reverseEdgeSeen = true
					}
					edge.Usable = false
					fCopy.Connections[station][edge.Rev].Usable = true

					if edge.To%2 == 0 {
						route = append(route, f.IDToName[edge.To])
					}

					break
				}
			}
		}

		flow++ // has to be here due to ExtractPaths functionality being tied to flow

		if reverseEdgeSeen { // existing routes have been modified, append what we have to [][][], then extract a new [][]
			routesSets = append(routesSets, routeSet)
			routeSet = fCopy.ExtractPaths(flow, start, end)
		} else { // no modifications, just append [] to current [][]
			route = append(route, f.IDToName[start])
			slices.Reverse(route)
			routeSet = append(routeSet, route)
		}

		if flow >= trains { // early exit, no need to get flow above the traincount
			break
		}
	}

	if flow == 0 {
		return nil, 0, 0
	}

	routesSets = append(routesSets, routeSet) // add the "leftover" [][] from !found
	fmt.Println(routesSets)

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
