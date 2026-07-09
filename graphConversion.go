package main

func convertIntoFlow(n *Network, start, end string) FlowNetwork {
	f := FlowNetwork{make(map[int]string), make(map[int][]Edge)}
	index := 0
	indLookUpMap := make(map[string]int)
	for station := range n.Connections {
		if station == start || station == end { //don't split start or end
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
		from := indLookUpMap[station]
		if station != start && station != end {
			from++
		}
		for _, conn := range connections {
			to := indLookUpMap[conn]
			addEdge(f.Connections, from, to)
		}
	}
	return f
}

func addEdge(edgeMap map[int][]Edge, from, to int) {
	edgeMap[from] = append(edgeMap[from], Edge{
		To:   to,
		Flow: false,
		Rev:  len(edgeMap[to]),
	})

	edgeMap[to] = append(edgeMap[to], Edge{
		To:   from,
		Flow: false,
		Rev:  len(edgeMap[from]) - 1,
	})
}
