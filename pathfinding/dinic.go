package pathfinding

import (
	"time"
)

// Dinic returns a set of distinct, non-overlapping routes for each flow increment
//
// It runs BFS once to set levels for each station and run DFS as many times as it finds a valid path
func (f *FlowNetwork) Dinic(start, end int) ([][][]string, int, time.Duration) {
	fCopy := f.CopyFlowNetwork()
	now := time.Now()

	routesSets := make([][][]string, 0)
	flow := 0

	for fCopy.BFSLevel(start, end) {
		clear(fCopy.StartEdges)

		for fCopy.DFSBlockingFlow(start, end) {
			flow++
			routesSets = append(routesSets, fCopy.ExtractPaths(flow, start, end))
		}
	}

	if flow == 0 {
		return nil, 0, 0
	}

	return routesSets, flow, time.Since(now)
}

// BFSLevel returns true if a path was found (i.e. if end station level was ever updated)
func (f *FlowNetwork) BFSLevel(start, end int) bool {
	queue := []int{start}
	clear(f.Levels)
	f.Levels[start] = 0

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		for _, edge := range f.Connections[curr] {
			if _, found := f.Levels[edge.To]; !found && edge.Usable {
				f.Levels[edge.To] = f.Levels[curr] + 1
				queue = append(queue, edge.To)
			}
		}
	}

	_, found := f.Levels[end]
	return found
}

// DFSBlockingFlow returns true if a path was found
//
// It checks recursively from the start station if it's possible to access a new station that satisfies: nextLevel = startLevel + 1
func (f *FlowNetwork) DFSBlockingFlow(start, end int) bool {
	if start == end {
		return true
	}

	startEdge := f.StartEdges[start]

	for edgeID, edge := range f.Connections[start][startEdge:] {
		f.StartEdges[start] = edgeID + startEdge

		if edge.Usable && f.Levels[start]+1 == f.Levels[edge.To] {
			if f.DFSBlockingFlow(edge.To, end) {
				edge.Usable = false
				f.Connections[edge.To][edge.Rev].Usable = true

				return true
			}
		}
	}

	f.Levels[start] = -1

	return false
}
