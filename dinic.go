package main

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

/*
func (f *FlowNetwork) DFSBlockingFlow(start, end int) bool {
	queue := []int{start}
	parents := make([]int, len(f.IDToName))

	for len(queue) > 0 {
		curr := queue[len(queue)-1]
		queue = queue[:len(queue)-1]

		if curr == end {
			// swap edges usability only if path was successfully found
			for station := end; station != start; station = parents[station] {
				prev := parents[station]

				for _, edge := range f.Connections[prev] {
					if edge.To == station && edge.Original {
						edge.Usable = false
						f.Connections[station][edge.Rev].Usable = true
						break
					}
				}
			}

			return true
		}

		for _, edge := range f.Connections[curr] {
			if edge.Usable && f.Levels[curr]+1 == f.Levels[edge.To] {
				queue = append(queue, edge.To)
				parents[edge.To] = curr
			}
		}
	}

	fmt.Printf("parents: %+v\n", parents)

	return false
}
*/
