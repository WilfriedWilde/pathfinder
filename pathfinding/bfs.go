package main

import "slices"

func bfs(n *Network, start, end string) ([]string, bool) {
	route := make([]string, 0)
	queue := make([]string, 0)
	visited := make(map[string]struct{})
	queue = append(queue, start)
	parentMap := make(map[string]string)
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		if curr == end {
			route = append(route, end)
			break
		}
		for _, connection := range n.Connections[curr] {
			_, ok := visited[connection]
			if ok {
				continue
			}
			visited[connection] = struct{}{}
			parentMap[connection] = curr
			queue = append(queue, connection)
		}
	}
	for {
		parent := parentMap[end]
		route = append(route, parent)
		end = parent
		if parent == start {
			break
		}
	}
	if len(route) != 0 {
		slices.Reverse(route)
		return route, true
	}
	return nil, false
}
