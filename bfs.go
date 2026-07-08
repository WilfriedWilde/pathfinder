package main

import (
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
