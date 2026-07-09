package main

import (
	"fmt"
	"slices"
)

func BFSFlow(f *FlowNetwork, start, end int) ([]string, bool) {
	queue := []int{start}
	visited := map[int]struct{}{start: {}}
	parentMap := make(map[int]int)

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr == end {
			route := make([]string, 0)

			for node := end; ; node = parentMap[node] {
				route = append(route, f.IDToName[node])
				fmt.Println(node)

				if node == start {
					break
				}
			}

			slices.Reverse(route)
			return route, true
		}

		for _, edge := range f.Connections[curr] {
			if _, found := visited[edge.To]; !found {
				visited[edge.To] = struct{}{}
				parentMap[edge.To] = curr
				queue = append(queue, edge.To)
			}
		}
	}

	return nil, false
}
