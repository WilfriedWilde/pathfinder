package main

func convertIntoFlow(n *Network, start, end string) FlowNetwork {
	f := FlowNetwork{make(map[int]string), make(map[string]int), make(map[int][]Edge)}
	index := 0
	for station := range n.Connections {
		if station == start || station == end { // don't split start or end
			f.IDToName[index] = station
			f.NameToID[station] = index
			index++
			continue
		}
		f.IDToName[index], f.IDToName[index+1] = station, station // in, out
		addEdge(f.Connections, index, index+1)                    // connect station_in into station_out
		f.NameToID[station] = index
		index += 2
	}

	for station, connections := range n.Connections {
		from := f.NameToID[station] // in, or start/end
		if station != start && station != end {
			from++ // out
		}
		for _, conn := range connections {
			to := f.NameToID[conn]           // in
			addEdge(f.Connections, from, to) // connect station_out (or start/end) into conn_in (or start/end)
		}
	}
	return f
}

// add forward edge and reverse edge at the same time, so the indices are correct
// not splitting the original connection, but redirecting it and then adding a reverse edge
func addEdge(edgeMap map[int][]Edge, from, to int) {
	edgeMap[from] = append(edgeMap[from], Edge{
		To:   to,
		Flow: false,
		Rev:  len(edgeMap[to]), // does not exist yet, will be created below
	})

	edgeMap[to] = append(edgeMap[to], Edge{
		To:   from,
		Flow: false,
		Rev:  len(edgeMap[from]) - 1, // was just created above
	})
}
