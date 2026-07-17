package pathfinding

import (
	"maps"
	"slices"
)

// TrimNetwork reduces network size by removing all stations, and their corresponding connections, that are not part of any path from start to end.
func (n *Network) TrimNetwork(start, end string) {
	forward := n.BFSFrom(start, start, end)
	backward := n.BFSFrom(end, start, end)

	// nothing to trim, return
	if len(forward) == len(backward) && len(forward) == len(n.Stations) {
		return
	}

	relevantStations := make(map[string]struct{})
	for station := range forward {
		if _, found := backward[station]; found {
			relevantStations[station] = struct{}{}
		}
	}

	// modify original stations and connections by keeping only the relevant stations
	for stationName := range n.Stations {
		if _, found := relevantStations[stationName]; !found {
			delete(n.Stations, stationName)
		}
	}

	for stationName := range n.Connections {
		if _, found := relevantStations[stationName]; !found {
			delete(n.Connections, stationName)
			continue
		}

		n.Connections[stationName] = slices.DeleteFunc(n.Connections[stationName], func(name string) bool {
			_, found := relevantStations[name]
			return !found
		})
	}
}

// BFSFrom returns a map of all visited stations from a specific start station.
func (n *Network) BFSFrom(station, start, end string) map[string]struct{} {
	queue := []string{station}
	visited := map[string]struct{}{station: {}}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if station == start && curr == end {
			visited[end] = struct{}{}
			continue
		} else if station == end && curr == start {
			visited[start] = struct{}{}
			continue
		}

		for _, conn := range n.Connections[curr] {
			if _, found := visited[conn]; !found {
				visited[conn] = struct{}{}
				queue = append(queue, conn)
			}
		}
	}

	return visited
}

func (n *Network) ConvertIntoFlow(start, end string) FlowNetwork {
	f := FlowNetwork{make(map[int]string), make(map[string]int), make(map[int][]*Edge), make(map[int]int), make(map[int]int), make([]int, len(n.Connections))}
	index := 2
	for station := range n.Connections {
		if station == start { // don't split start or end
			f.IDToName[0] = station
			f.NameToID[station] = 0
			continue
		}
		if station == end { // don't split start or end
			f.IDToName[1] = station
			f.NameToID[station] = 1
			continue
		}
		f.IDToName[index], f.IDToName[index+1] = station, station // in, out
		addEdge(f.Connections, index, index+1, 0)                 // connect station_in into station_out
		f.NameToID[station] = index
		index += 2
	}

	for station, connections := range n.Connections {
		from := f.NameToID[station] // in, or start/end
		if station != start && station != end {
			from++ // out
		}
		for _, conn := range connections {
			to := f.NameToID[conn]              // in
			addEdge(f.Connections, from, to, 1) // connect station_out (or start/end) into conn_in (or start/end)
		}
	}
	return f
}

// add forward edge and reverse edge at the same time, so the indices are correct
// not splitting the original connection, but redirecting it and then adding a reverse edge
// reverse edge only used for potential backtracking? never for main pathfinding?
func addEdge(edgeMap map[int][]*Edge, from, to, cost int) {
	edgeMap[from] = append(edgeMap[from], &Edge{
		To:       to,
		Usable:   true, // all forward edges are traversable at start, none of the reverse edges are - swapping booleans upon using an edge should enable the reverse and disable the used edge
		Original: true,
		Rev:      len(edgeMap[to]), // does not exist yet, will be created below
		Cost:     cost,
	})

	edgeMap[to] = append(edgeMap[to], &Edge{
		To:       from,
		Usable:   false,
		Original: false,
		Rev:      len(edgeMap[from]) - 1, // was just created above
		Cost:     -cost,
	})
}

func (f *FlowNetwork) CopyFlowNetwork() *FlowNetwork {
	copy := &FlowNetwork{
		IDToName:    make(map[int]string, len(f.IDToName)),
		NameToID:    make(map[string]int, len(f.NameToID)),
		Connections: make(map[int][]*Edge, len(f.Connections)),
		Levels:      make(map[int]int, len(f.NameToID)),
		StartEdges:  make(map[int]int, len(f.NameToID)),
		Potential:   make([]int, len(f.Connections)),
	}

	maps.Copy(copy.IDToName, f.IDToName)
	maps.Copy(copy.NameToID, f.NameToID)

	for stationID, edges := range f.Connections {
		newEdges := make([]*Edge, len(edges))

		for i, edge := range edges {
			edgeCopy := *edge
			newEdges[i] = &edgeCopy
		}

		copy.Connections[stationID] = newEdges
	}

	return copy
}

func copyNetwork(n *Network) *Network {
	network := &Network{make(map[string]Station), make(map[string][]string)}

	maps.Copy(network.Stations, n.Stations)
	maps.Copy(network.Connections, n.Connections)

	return network
}
