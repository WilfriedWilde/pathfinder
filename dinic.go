package main

import "fmt"

func (f *FlowNetwork) Dinic(start, end int) ([][][]string, bool) {
	flow := 0

	for f.BFSLevel(start, end) {
		fmt.Println("inside BFS loop")
		for f.DFSBlockingFlow(start, end) {
			flow++
		}
	}

	fmt.Println(flow)
	return nil, false
}

func (f *FlowNetwork) DFSBlockingFlow(start, end int) bool {
	queue := []int{start}
	parents := make([]int, len(f.IDToName))

	for len(queue) > 0 {
		fmt.Println("inside DFS loop")
		curr := queue[len(queue)-1]
		queue = queue[:len(queue)-1]

		if curr == end {
			fmt.Printf("parents: %+v\n", parents)
			return true
		}

		for _, edge := range f.Connections[curr] {
			if edge.Usable && f.Levels[curr]+1 == f.Levels[edge.To] {
				queue = append(queue, edge.To)
				parents[edge.To] = curr
			}
		}
	}

	fmt.Printf("parents: %+v\n", parents)

	return false
}
