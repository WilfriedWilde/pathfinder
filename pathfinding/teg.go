package pathfinding

import "slices"

type TEGNode struct {
	StationID int
	Time      int
}

type TEGNetwork struct {
	FlowNetwork
	NodeToID map[TEGNode]int
	IDToNode map[int]TEGNode
}

func (f *FlowNetwork) CreateTEGNetwork(maxTime int) *TEGNetwork {
	teg := &TEGNetwork{
		FlowNetwork: FlowNetwork{
			IDToName:    make(map[int]string),
			NameToID:    make(map[string]int),
			Connections: make(map[int][]*Edge),
		},
		NodeToID: make(map[TEGNode]int),
		IDToNode: make(map[int]TEGNode),
	}

	id := 0

	for station := range f.IDToName {
		for i := 0; i <= maxTime; i++ {
			phNode := TEGNode{station, i}
			teg.IDToNode[id] = phNode
			teg.NodeToID[phNode] = id

			name := f.IDToName[station]
			teg.IDToName[id] = name
			teg.NameToID[name] = id

			id++
		}
	}
	return teg
}

func (t *TEGNetwork) AddNeighborEdges(f *FlowNetwork, maxTime int) {
	for station, edges := range f.Connections {
		for _, edge := range edges {
			if !edge.Original {
				continue
			}
			for i := 0; i < maxTime; i++ {
				from := TEGNode{station, i}
				to := TEGNode{edge.To, i + 1}
				addEdge(t.Connections, t.NodeToID[from], t.NodeToID[to], 1)
			}
		}
	}
}

func (t *TEGNetwork) AddSelfEdges(maxTime int) {
	for node, id := range t.NodeToID {
		if node.Time == maxTime {
			continue
		}
		next := TEGNode{node.StationID, node.Time + 1}
		addEdge(t.Connections, id, t.NodeToID[next], 1)
	}
}

func (f *FlowNetwork) EstimateTEGViability(trains int) (int, bool) {
	shortestPath, found := BFSFlow(f, 0, 1) // start, end
	if !found {
		return 0, false
	}
	maxTime := len(shortestPath) + trains
	stationCount := len(f.IDToName) * maxTime

	edgeCount := 0

	for _, edges := range f.Connections {
		edgeCount += len(edges)
	}

	edgeCount *= maxTime

	if stationCount < 30000 && edgeCount < 150000 {
		return maxTime, true
	}
	return 0, false
}

func (t *TEGNetwork) GetStationID(station, time int) int {
	return t.NodeToID[TEGNode{station, time}]
}

func (t *TEGNetwork) FindRoute(start, end, maxTime int) ([]int, int, bool) {
	for i := 0; i <= maxTime; i++ {
		startID := t.GetStationID(start, 0)
		endID := t.GetStationID(end, i)

		route, found := BFSFlow(&t.FlowNetwork, startID, endID)
		if found {
			return route, endID, true
		}
	}
	return nil, -1, false
}

func (t *TEGNetwork) ExtractPath(parents []int, start, end int) []TEGNode {
	res := []TEGNode{}

	for end != -1 {
		res = append(res, t.IDToNode[end])
		end = parents[end]
	}
	slices.Reverse(res)
	return res
}
