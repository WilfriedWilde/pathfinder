package main

func (f *FlowNetwork) EK(start, end int) ([][][]string, int) {
	ret := make([][][]string, 0)
	flow := 0
	for {
		usedEdges, found := f.BFSFlow(start, end)
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
		ret = append(ret, f.ExtractPaths(flow, start, end))

	}
	if flow == 0 {
		return nil, 0
	}
	return ret, flow
}

func (f *FlowNetwork) ExtractPaths(flow, start, end int) [][]string {

	paths := make([][]string, 0)
	path := make([]string, 0)

	modifiedIndexes := make([]int, 0)
	for i := 0; i < flow; i++ {
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
	}
	// undo graph modifications
	for _, edgeInd := range modifiedIndexes {
		f.Connections[start][edgeInd].Usable = false
	}
	return paths
}

func (f *FlowNetwork) BFSFlow(start, end int) ([]int, bool) {
	queue := []int{start}
	visited := map[int]struct{}{start: {}}
	parentMap := make([]int, len(f.IDToName))
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr == end {
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
