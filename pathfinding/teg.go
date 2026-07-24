package pathfinding

import (
	"math"
)

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
		for time := 0; time <= maxTime; time++ {
			phNode := TEGNode{station, time}
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
			if edge.Cost == 0 {
				for time := 0; time <= maxTime; time++ {
					from := TEGNode{station, time}
					to := TEGNode{edge.To, time}
					addEdge(t.Connections, t.NodeToID[from], t.NodeToID[to], 0)
				}
			} else {
				for time := 0; time < maxTime; time++ {
					from := TEGNode{station, time}
					to := TEGNode{edge.To, time + edge.Cost}
					addEdge(t.Connections, t.NodeToID[from], t.NodeToID[to], 1)
				}
			}

		}
	}
}

func (t *TEGNetwork) AddSelfEdges(maxTime, start, end int) {
	for node, id := range t.NodeToID {
		if node.Time == maxTime {
			continue
		}
		next := TEGNode{node.StationID, node.Time + 1}
		addEdge(t.Connections, id, t.NodeToID[next], 0)
		if node.StationID == start || node.StationID == end {
			edges := t.Connections[id]
			edges[len(edges)-1].Unlimited = true
		}
	}
}

func (n *Network) EstimateTEGViability(start, end string, trains int) (int, bool) {
	shortestPath, found := BFS(n, start, end) // start, end
	if !found {
		return 0, false
	}
	maxTime := len(shortestPath) + trains
	stationCount := len(n.Connections) * 2 // roughly account for splits
	stationCount *= maxTime                // # of layers in teg

	edgeCount := 0

	for _, connections := range n.Connections {
		edgeCount += len(connections)
	}

	edgeCount = edgeCount*2 + len(n.Connections)*2 // split all, account for new edges between station_in and station_out and for waiting edges between layers
	edgeCount *= maxTime                           // roughly account for teg layers

	if stationCount < 30000 && edgeCount < 150000 {
		return maxTime, true
	}
	return 0, false
}

func (t *TEGNetwork) FindRoute(start, end, maxTime int) ([]int, int, bool) {
	startID := t.NodeToID[TEGNode{start, 0}]

	bestID := -1
	bestDistance := math.MaxInt32
	var bestParents []int

	for i := 0; i <= maxTime; i++ {
		endID := t.NodeToID[TEGNode{end, i}]
		parents, found := BFSFlow(&t.FlowNetwork, startID, endID)
		if !found {
			continue
		}

		distance := 0
		for curr := endID; curr != startID; curr = parents[curr] {
			distance++
		}
		if distance < bestDistance {
			bestDistance = distance
			bestID = endID
			bestParents = parents
		}
	}

	if bestID == -1 {
		return nil, -1, false
	}
	return bestParents, bestID, true
}

func (t *TEGNetwork) DoableInXTurns(f *FlowNetwork, start, end, trainNum, maxTurns int) ([][]string, bool) {

	routesSets := make([][]string, 0)
	flow := 0
	for {

		usedEdges, endID, found := t.FindRoute(start, end, maxTurns)
		if !found {
			break
		}

		startID := t.NodeToID[TEGNode{start, 0}]

		for station := endID; station != startID; station = usedEdges[station] {

			prev := usedEdges[station]
			for _, edge := range t.Connections[prev] {

				if edge.Original && edge.To == station && edge.Usable {
					if !edge.Unlimited {
						edge.Usable = false
						t.Connections[station][edge.Rev].Usable = true
					}
					break
				}

			}
		}

		flow++
		if flow == trainNum {
			routesSets = append(routesSets, t.ExtractTEGPaths(f, flow, start, end)...)
		}
	}

	if flow < trainNum {
		return nil, false
	}
	return routesSets, true
}

func (t *TEGNetwork) ExtractTEGPaths(f *FlowNetwork, flow, start, end int) [][]string {
	paths := make([][]string, 0)
	path := make([]string, 0)

	for range flow {

		startID := t.NodeToID[TEGNode{start, 0}]

		curr := startID
		path = []string{}

		for t.IDToNode[curr].StationID != end {

			for _, edge := range t.Connections[curr] {

				if edge.Original && (edge.Unlimited || !edge.Usable) {
					if !edge.Unlimited {
						edge.Usable = true
					}
					if t.IDToNode[curr].StationID%2 == 0 {
						path = append(path, f.IDToName[t.IDToNode[curr].StationID])
					}

					curr = edge.To
					break
				}
			}
		}

		path = append(path, f.IDToName[end])
		paths = append(paths, path)
	}

	return paths
}
