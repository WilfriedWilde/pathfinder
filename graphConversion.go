package main

func convertIntoFlow(n *Network, start, end string) FlowNetwork {
	f := FlowNetwork{make(map[int]string), make(map[int][]Edge)}
	index := 0
	indLookUpMap := make(map[string]int) // two separate maps for index <--> name switcheroos, only need to keep one alive outside of this function
	for station := range n.Connections {
		if station == start || station == end { // don't split start or end
			f.Name[index] = station
			indLookUpMap[station] = index
			index++
			continue
		}
		f.Name[index], f.Name[index+1] = station, station // in, out
		addEdge(f.Connections, index, index+1)            // connect station_in into station_out
		indLookUpMap[station] = index
		index += 2
	}

	for station, connections := range n.Connections {
		from := indLookUpMap[station] // out
		if station != start && station != end {
			from++
		}
		for _, conn := range connections {
			to := indLookUpMap[conn]         // in
			addEdge(f.Connections, from, to) // connect station_out into conn_in
		}
	}
	return f
}

// always create edge and its reverse (the old single connection, split into two) at the same time to keep indices correct
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
