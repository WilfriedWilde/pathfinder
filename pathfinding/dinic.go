package pathfinding

import (
	"time"
)

func (f *FlowNetwork) Dinic(start, end int) ([][][]string, int, time.Duration) {
	now := time.Now()

	routesSets := make([][][]string, 0)
	flow := 0

	for f.BFSLevel(start, end) {
		clear(f.StartEdges)

		for f.DFSBlockingFlow(start, end) {
			flow++
			routesSets = append(routesSets, f.ExtractPaths(flow, start, end))
		}
	}

	if flow == 0 {
		return nil, 0, 0
	}

	return routesSets, flow, time.Since(now)
}

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
