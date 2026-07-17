package pathfinding

func (f *FlowNetwork) BFSFlow(start, end int) ([]int, bool) {
	queue := []int{start}
	visited := map[int]struct{}{start: {}}
	parents := make([]int, len(f.IDToName))

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
