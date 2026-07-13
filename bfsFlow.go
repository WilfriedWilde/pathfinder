package main

import (
	"fmt"
	"slices"
)

func EK(f *FlowNetwork, start, end int) ([][]string, int) {
	ret := make([][]string, 0)
	flow := 0
	for {
		usedEdges, found := BFSFlow(f, start, end)
		if !found {
			break
		}

		for node := end; node != start; node = usedEdges[node] {
			next := usedEdges[node]
			for _, edge := range f.Connections[next] {
				if edge.To == node {
					edge.Usable = false
					f.Connections[node][edge.Rev].Usable = true
					break
				}
			}
		}
		flow++
		if flow == 3 {
			ret = append(ret, ExtractPaths(f, flow, start, end)...)
		}

	}
	if flow == 0 {
		return nil, 0
	}
	return ret, flow
}

func ExtractPaths(f *FlowNetwork, flow, start, end int) [][]string {

	paths := make([][]string, 0)
	path := make([]string, 0)
	for i := 0; i < flow; i++ {
		fmt.Println(i)

		for station, edges := range f.Connections {
			for _, edge := range edges {
				fmt.Printf("%+v, %+v, %+v\n", f.IDToName[station], edge, f.IDToName[edge.To])
			}
		}
		curr := start
		path = []string{}
		for curr != end {
			for _, edge := range f.Connections[curr] {
				if edge.Original && !edge.Usable {
					if curr == start {
						edge.Usable = true
					}
					if curr%2 == 0 {
						path = append(path, f.IDToName[curr])
					}

					curr = edge.To
					break
				}
			}
		}
		path = append(path, f.IDToName[end])
		paths = append(paths, path)
		for _, v := range paths {
			fmt.Println(v)
		}
	}
	return paths
}

func BFSFlow(f *FlowNetwork, start, end int) (map[int]int, bool) {
	queue := []int{start}
	visited := map[int]struct{}{start: {}}
	parentMap := make(map[int]int)
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr == end {
			route := make([]string, 0)

			for node := end; ; node = parentMap[node] {
				if !slices.Contains(route, f.IDToName[node]) {
					route = append(route, f.IDToName[node])
				}
				if node == start {
					break
				}
			}

			slices.Reverse(route)
			return parentMap, true
		}

		for _, edge := range f.Connections[curr] {
			if _, found := visited[edge.To]; !found {
				if edge.Usable {
					visited[edge.To] = struct{}{}
					parentMap[edge.To] = curr
					queue = append(queue, edge.To)
				}

			}
		}
	}

	return nil, false
}
