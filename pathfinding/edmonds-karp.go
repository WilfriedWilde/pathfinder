package pathfinding

// BFSFlow returns a slice where the index is the stationID and the value is its parent (station where the station of the given index can be reached in the least amount of edges)
func BFSFlow(f *FlowNetwork, start, end int) ([]int, bool) {
	queue := []int{start}
	visited := map[int]struct{}{start: {}}
	parents := make([]int, len(f.IDToName))
	for i := range parents {
		parents[i] = -1
	}

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
