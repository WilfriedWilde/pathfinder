package main

import (
	"maps"
)

func removeShortest(n *Network, start, end string) ([][]string, bool) {
	routes := make([][]string, 0)
	network := copyNetwork(n) //copy so original network doesn't get modified, in case it's useful elsewhere

	for { //keep running BFS, removing middle stations (and start-end connection specifically if that's found) until no more paths exist
		route, found := BFS(network, start, end)
		if found == false {
			break
		}

		routes = append(routes, route)

		if len(route) == 2 { //has to be start-end
			ind := 0

			for i, v := range network.Connections[start] {
				if v == end {
					ind = i
					break
				}
			}

			network.Connections[start] = append(network.Connections[start][:ind], network.Connections[start][ind+1:]...)

		}

		removeMap := make(map[string]struct{})

		for _, stationName := range route {
			if stationName != start && stationName != end {
				removeMap[stationName] = struct{}{}
			}
		}

		// loop over all station[connections], compare against ones used in route, build a new []connection without the used conns
		for stationName, connections := range network.Connections {
			filtered := make([]string, 0)

			for _, conn := range connections {
				if _, remove := removeMap[conn]; !remove {
					filtered = append(filtered, conn)
				}
			}

			network.Connections[stationName] = filtered
		}
	}

	if len(routes) == 0 {
		return nil, false
	}

	return routes, true
}

func copyNetwork(n *Network) *Network {
	network := &Network{make(map[string]Station), make(map[string][]string)}

	maps.Copy(network.Stations, n.Stations)
	maps.Copy(network.Connections, n.Connections)

	return network
}
