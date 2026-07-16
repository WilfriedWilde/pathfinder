package main

import (
	"fmt"
	"maps"
	"slices"
)

func BFS(n *Network, start, end string) ([]string, bool) {
	queue := []string{start}
	visited := map[string]struct{}{start: {}}
	parentMap := make(map[string]string)

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr == end {
			route := make([]string, 0)

			for node := end; ; node = parentMap[node] {
				route = append(route, node)

				if node == start {
					break
				}
			}

			slices.Reverse(route)
			return route, true
		}

		for _, connection := range n.Connections[curr] {
			if _, found := visited[connection]; !found {
				visited[connection] = struct{}{}
				parentMap[connection] = curr
				queue = append(queue, connection)
			}
		}
	}

	return nil, false
}

// BFSFrom returns a map of all visited stations from a specific start station.
func (n *Network) BFSFrom(station string) map[string]struct{} {
	queue := []string{station}
	visited := map[string]struct{}{station: {}}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		for _, conn := range n.Connections[curr] {
			if _, found := visited[conn]; !found {
				visited[conn] = struct{}{}
				queue = append(queue, conn)
			}
		}
	}

	return visited
}

func (f *FlowNetwork) BFSFlow(start, end int) ([]int, bool) {
	queue := []int{start}
	visited := map[int]struct{}{start: {}}
	parents := make([]int, len(f.IDToName))

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr == end {
			return parents, true
		}

		for _, edge := range f.Connections[curr] {
			if _, found := visited[edge.To]; !found {
				if edge.Usable {
					visited[edge.To] = struct{}{}
					parents[edge.To] = curr
					queue = append(queue, edge.To)
				}

			}
		}
	}

	return nil, false
}

func (f *FlowNetwork) BFSLevel(start, end int) bool {
	queue := []int{start}
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

	for station, level := range f.Levels {
		fmt.Printf("station %d has level %d\n", station, level)
	}

	return false
}

func (n *Network) RemoveShortest(start, end string) ([][]string, bool) {
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
