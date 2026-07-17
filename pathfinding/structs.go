package pathfinding

type Station struct {
	Name string
	X, Y int
}

type Network struct {
	Stations    map[string]Station
	Connections map[string][]string
}

type FlowNetwork struct {
	IDToName    map[int]string
	NameToID    map[string]int
	Connections map[int][]*Edge
	Levels      map[int]int // For each station, stores its level (distance from the start to this station in number of stations) - Used in BFS level
	StartEdges  map[int]int // For each station, stores the start edge index - Used in the DFS blocking flow
	Potential   []int       // Used for Dijkstra implementation
}

type Edge struct {
	To       int
	Usable   bool
	Original bool
	Rev      int
	Cost     int // For Bellman-Ford and Dijkstra
	Level    int // Used for Dinic
}

type TrainRoute struct {
	RouteLength int
	Route       []string
	Trains      int
}
