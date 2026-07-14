package main

import (
	"fmt"
	"math"
)

func MCMF(f *FlowNetwork, start, end int) ([][][]string, int) {
	ret := make([][][]string, 0)
	flow := 0
	for {
		distance, parents := bellmanFord(f, start, end)
		if distance[end] == math.MaxInt32 { // no path found
			break
		}
		for node := end; node != start; node = parents[node] {
			prev := parents[node]
			for _, edge := range f.Connections[prev] {
				if edge.To == node {
					edge.Usable = false
					f.Connections[node][edge.Rev].Usable = true
					break
				}
			}
		}
		flow++
		ret = append(ret, ExtractPaths(f, flow, start, end))

	}
	if flow == 0 {
		return nil, 0
	}
	return ret, flow
}

func ExtractPathsMCMF(f *FlowNetwork, flow, start, end int) [][]string {

	paths := make([][]string, 0)
	path := make([]string, 0)

	modifiedIndexes := make([]int, 0)
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
			for ind, edge := range f.Connections[curr] {
				if edge.Original && !edge.Usable {
					if curr == start {
						edge.Usable = true
						modifiedIndexes = append(modifiedIndexes, ind)
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
	// undo graph modifications
	for _, edgeInd := range modifiedIndexes {
		f.Connections[start][edgeInd].Usable = false
	}
	return paths
}

func bellmanFord(f *FlowNetwork, start, end int) ([]int, []int) {
	distance := make([]int, len(f.IDToName))
	parents := make([]int, len(f.IDToName))
	for i := range distance {
		distance[i] = math.MaxInt32
		parents[i] = -1
	}
	distance[start] = 0
	for i := 0; i < len(distance)-1; i++ {
		diff := false
		for station, edges := range f.Connections {
			for _, edge := range edges {
				if distance[station] != math.MaxInt32 &&
					distance[station]+edge.Cost < distance[edge.To] &&
					edge.Usable {

					distance[edge.To] = distance[station] + edge.Cost
					parents[edge.To] = station
					diff = true
				}
			}
		}
		if !diff {
			break
		}
	}

	return distance, parents
}
