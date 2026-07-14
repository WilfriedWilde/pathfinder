package main

import (
	"math"
)

func MCMF(f *FlowNetwork, start, end int) ([][][]string, int) {
	ret := make([][][]string, 0)
	flow := 0
	for {
		parents, found := bellmanFord(f, start, end)
		if !found {
			break
		}
		for node := end; node != start; node = parents[node] {
			prev := parents[node]
			for _, edge := range f.Connections[prev] {
				if edge.To == node {
					edge.Usable = false
					f.Connections[node][edge.Rev].Usable = true
					break
				}
			}
		}
		flow++
		ret = append(ret, ExtractPaths(f, flow, start, end))

	}
	if flow == 0 {
		return nil, 0
	}
	return ret, flow
}

func bellmanFord(f *FlowNetwork, start, end int) ([]int, bool) {
	distance := make([]int, len(f.IDToName)) // stationID is the index, value is the distance from start station to stationID
	parents := make([]int, len(f.IDToName))  // stationID index, value stationID of the station from which the shortest path has been found
	for i := range distance {
		distance[i] = math.MaxInt32
		parents[i] = -1
	}
	distance[start] = 0
	// loop over all stations, act on the stations that have a known distance to start
	for i := 0; i < len(distance)-1; i++ {
		diff := false
		// for every station
		for station, edges := range f.Connections {
			// check all edges and if station has known distance, add that distance + edge cost to the stations the edges connect to
			for _, edge := range edges {
				if distance[station] != math.MaxInt32 &&
					distance[station]+edge.Cost < distance[edge.To] &&
					edge.Usable {

					distance[edge.To] = distance[station] + edge.Cost
					parents[edge.To] = station
					diff = true
				}
			}
		}
		if !diff {
			break
		}
	}
	if distance[end] == math.MaxInt32 { // no path found
		return nil, false
	}

	return parents, true
}
