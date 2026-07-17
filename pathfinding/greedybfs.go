package pathfinding

import (
	"slices"
	"time"
)

// GreedyBFS returns a set of distinct, non-overlapping routes
//
// It repeatedly runs BFS and removes the non-terminal stations from the graph
func (n *Network) GreedyBFS(start, end string) ([][]string, bool, time.Duration) {
	nCopy := n.CopyNetwork()
	now := time.Now()
	routes := make([][]string, 0)

	for { //keep running BFS, removing middle stations (and start-end connection specifically if that's found) until no more paths exist
		route, found := BFS(nCopy, start, end)
		if found == false {
			break
		}

		routes = append(routes, route)

		if len(route) == 2 { //has to be start-end
			ind := 0

			for i, v := range nCopy.Connections[start] {
				if v == end {
					ind = i
					break
				}
			}

			nCopy.Connections[start] = append(nCopy.Connections[start][:ind], nCopy.Connections[start][ind+1:]...)

		}

		removeMap := make(map[string]struct{})

		for _, stationName := range route {
			if stationName != start && stationName != end {
				removeMap[stationName] = struct{}{}
			}
		}

		// loop over all station[connections], compare against ones used in route, build a new []connection without the used conns
		for stationName, connections := range nCopy.Connections {
			filtered := make([]string, 0)

			for _, conn := range connections {
				if _, remove := removeMap[conn]; !remove {
					filtered = append(filtered, conn)
				}
			}

			nCopy.Connections[stationName] = filtered
		}
	}

	if len(routes) == 0 {
		return nil, false, 0
	}

	return routes, true, time.Since(now)
}

// BFS returns the route used, as well as if it found a route
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
