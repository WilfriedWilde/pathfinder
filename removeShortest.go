package main

func removeShortest(n *Network, start, end string) ([][]string, bool) {
	ret := make([][]string, 0)
	copy := copyNetwork(n) //copy so original network doesn't get modified, in case it's useful elsewhere
	for {                  //keep running BFS, removing middle stations (and start-end connection specifically if that's found) until no more paths exist
		route, found := BFS(copy, start, end)
		if found == false {
			break
		}
		ret = append(ret, route)
		if len(route) == 2 { //has to be start-end
			ind := 0
			for i, v := range copy.Connections[start] {
				if v == end {
					ind = i
					break
				}
			}
			copy.Connections[start] = append(copy.Connections[start][:ind], copy.Connections[start][ind+1:]...)
		}
		removeMap := make(map[string]struct{})
		for _, stationName := range route {
			if stationName != start && stationName != end {
				removeMap[stationName] = struct{}{}
			}

		}
		// loop over all station[connections], compare against ones used in route, build a new []connection without the used conns
		for stationName, connections := range copy.Connections {
			filtered := make([]string, 0)
			for _, conn := range connections {
				if _, remove := removeMap[conn]; !remove {
					filtered = append(filtered, conn)
				}
			}
			copy.Connections[stationName] = filtered
		}
	}
	if len(ret) == 0 {
		return nil, false
	}
	return ret, true
}

func copyNetwork(n *Network) *Network {
	copy := &Network{make(map[string]Station), make(map[string][]string)}
	for name, station := range n.Stations {
		copy.Stations[name] = station
	}
	for station, connections := range n.Connections {
		copy.Connections[station] = append([]string(nil), connections...)
	}
	return copy
}
