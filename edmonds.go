package main

func (f *FlowNetwork) GetRoutesSets(start, end int, pathfinder func(start, end int) ([]int, bool)) ([][][]string, int) {
	routesSets := make([][][]string, 0)
	flow := 0
	for {
		usedEdges, found := pathfinder(start, end)
		if !found {
			break
		}

		for station := end; station != start; station = usedEdges[station] {
			prev := usedEdges[station]
			for _, edge := range f.Connections[prev] {
				if edge.To == station {
					edge.Usable = false
					f.Connections[station][edge.Rev].Usable = true
					break
				}
			}
		}

		flow++
		routesSets = append(routesSets, f.ExtractPaths(flow, start, end))
	}

	if flow == 0 {
		return nil, 0
	}

	return routesSets, flow
}

func (f *FlowNetwork) ExtractPaths(flow, start, end int) [][]string {

	paths := make([][]string, 0)
	path := make([]string, 0)

	modifiedIndexes := make([]int, 0)
	for range flow {
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
