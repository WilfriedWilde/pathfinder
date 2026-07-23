package pathfinding

import (
	"container/heap"
	"math"
)

type Node struct {
	ID   int
	Dist int
}

type PriorityQueue []Node

// BellmanFord returns a slice where the index is the stationID and the value is its parent (station where the station of the given index can be reached in the least amount of actual distance)
func BellmanFord(f *FlowNetwork, start, end int) ([]int, bool) {
	distance := make([]int, len(f.IDToName)) // stationID is the index, value is the distance from start station to stationID
	parents := make([]int, len(f.IDToName))  // stationID index, value stationID of the station from which the shortest path has been found

	for i := range distance {
		distance[i] = math.MaxInt32
		parents[i] = -1
	}

	distance[start] = 0
	// loop over all stations, act on the stations that have a known distance to start
	for range distance[:len(distance)-1] {
		diff := false
		// for every station
		for station, edges := range f.Connections {
			// skip stations that have not yet been connected into
			if distance[station] == math.MaxInt32 {
				continue
			}
			// check all edges and if station has known distance, add that distance + edge cost to the stations the edges connect to
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

	if distance[end] == math.MaxInt32 { // no path found
		return nil, false
	}

	return parents, true
}

// Dijkstra returns a slice where the index is the stationID and the value is its parent (station where the station of the given index can be reached in the least amount of actual distance)
//
// # Compared to BFS, uses a priority queue instead, always popping the node which currently has the least distance stored into it
//
// Since Dijkstra cannot handle negative cost paths, Johnson's potentials are used to normalize the path costs to be positive, while retaining their relative cost differences
func Dijkstra(f *FlowNetwork, start, end int) ([]int, bool) {
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

			if distance[curr.ID]+cost < distance[edge.To] {
				distance[edge.To] = distance[curr.ID] + cost
				parents[edge.To] = curr.ID
				heap.Push(prioQueue, Node{edge.To, distance[edge.To]})
			}
		}
	}

	if distance[end] == math.MaxInt32 { // no path found
		return nil, false
	}

	for i := range distance {
		if distance[i] != math.MaxInt32 {
			f.Potential[i] += distance[i]
		}
	}

	return parents, true
}

// required for being able to implement a container/heap interface

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Dist < pq[j].Dist
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x any) {
	*pq = append(*pq, x.(Node))
}

func (pq *PriorityQueue) Pop() any {
	n := len(*pq)
	popped := (*pq)[n-1]
	*pq = (*pq)[:n-1]
	return popped
}
