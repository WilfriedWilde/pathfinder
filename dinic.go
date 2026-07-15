package main

func (f *FlowNetwork) Dinic(start, end int) ([][][]string, bool) {
	// Initialize levels
	for _, edges := range f.Connections {
		for _, edge := range edges {
			edge.Level = -1
		}
	}

	return nil, false
}
