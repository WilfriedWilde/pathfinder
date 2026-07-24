package pathfinding

import (
	"container/heap"
	"fmt"
	"math"
	"slices"
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
	teg.FlowNetwork.Potential = make([]int, len(teg.NodeToID))
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
					to := TEGNode{edge.To, time + edge.Cost}
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

func (t *TEGNetwork) GetStationID(station, time int) int {
	return t.NodeToID[TEGNode{station, time}]
}

func (t *TEGNetwork) FindRoute(start, end, maxTime int) ([]int, int, bool) {
	startID := t.NodeToID[TEGNode{start, 0}]
	parents, distance, _ := TEGDijkstra(&t.FlowNetwork, startID)
	bestID := -1
	bestDist := math.MaxInt32
	for i := 0; i <= maxTime; i++ {
		endID := t.NodeToID[TEGNode{end, i}]
		if distance[endID] < bestDist {
			bestDist = distance[endID]
			bestID = endID
		}
	}
	if bestID == -1 {
		return nil, -1, false
	}
	return parents, bestID, true
}

func (t *TEGNetwork) ExtractPath(parents []int, start, end int) []TEGNode {
	res := []TEGNode{t.IDToNode[end]}

	for end != -1 {
		fmt.Println(t.IDToNode[end])
		if t.IDToNode[end].StationID%2 == 0 {
			res = append(res, t.IDToNode[end])
		}

		end = parents[end]
	}
	slices.Reverse(res)
	return res
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

func TEGDijkstra(f *FlowNetwork, start int) ([]int, []int, bool) {
	distance := make([]int, len(f.IDToName))
	parents := make([]int, len(f.IDToName))
	for i := range distance {
		distance[i] = math.MaxInt32
		parents[i] = -1
	}

	distance[start] = 0

	prioQueue := &PriorityQueue{}
	heap.Init(prioQueue)
	heap.Push(prioQueue, Node{start, 0})

	for prioQueue.Len() != 0 {
		curr := heap.Pop(prioQueue).(Node)

		if curr.Dist != distance[curr.ID] {
			continue
		}

		for _, edge := range f.Connections[curr.ID] {
			if !edge.Usable {
				continue
			}

			cost := edge.Cost + f.Potential[curr.ID] - f.Potential[edge.To]
			if cost < 0 {
				panic("negative")
			}

			if distance[curr.ID]+cost < distance[edge.To] {
				distance[edge.To] = distance[curr.ID] + cost
				parents[edge.To] = curr.ID
				heap.Push(prioQueue, Node{edge.To, distance[edge.To]})
			}
		}
	}

	for i := range distance {
		if distance[i] != math.MaxInt32 {
			f.Potential[i] += distance[i]
		}
	}

	return parents, distance, true
}
