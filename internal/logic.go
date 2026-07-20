package internal

import (
	"os"
	p "pathfinder/pathfinding"
	"slices"
	"time"
)

type Result struct {
	Name        string // algorithm name
	TrainRoutes []p.TrainRoute
	MaxTurns    int
	TimeTaken   time.Duration
}

type PageContent struct {
	GraphVisibility bool
	Results         []Result
	Routes          []string
}

type Config struct {
	NetworkPath  string
	StartStation string
	EndStation   string
	NumTrains    int
}

func RunLogic(args []string, algorithms []string) (PageContent, error) {
	config, err := GetAndValidateUserInput(os.Args[1:])
	if err != nil {
		return PageContent{}, err
	}

	network, NumLines, err := ParseAndValidateNetwork(config)
	if err != nil {
		return PageContent{}, err
	}

	content := PageContent{
		GraphVisibility: true,
		Results:         make([]Result, 0),
		Routes:          make([]string, 0),
	}

	if NumLines > 25000 {
		content.GraphVisibility = false
	}

	network.TrimNetwork(config.StartStation, config.EndStation)

	edgeNetwork := network.ConvertIntoFlow(config.StartStation, config.EndStation)

	if slices.Contains(algorithms, "greedy") {
		repeatedBFSRoutes, _, bfsTime := network.GreedyBFS(config.StartStation, config.EndStation)
		trainRoutes, maxTurnsBFS := divert(repeatedBFSRoutes, config.NumTrains)

		result := Result{
			Name:        "greedy",
			TrainRoutes: trainRoutes,
			MaxTurns:    maxTurnsBFS,
			TimeTaken:   bfsTime,
		}

		content.Results = append(content.Results, result)
	}

	if slices.Contains(algorithms, "edmonds") {
		ekRouteSets, _, ekTime := edgeNetwork.GetRoutesSets(edgeNetwork.NameToID[config.StartStation], edgeNetwork.NameToID[config.EndStation], p.BFSFlow)
		fastestEKRoutes := GetFastestRouteSet(ekRouteSets, config.NumTrains)
		trainRoutes, maxTurnsEK := divert(fastestEKRoutes, config.NumTrains)

		result := Result{
			Name:        "edmonds",
			TrainRoutes: trainRoutes,
			MaxTurns:    maxTurnsEK,
			TimeTaken:   ekTime,
		}

		content.Results = append(content.Results, result)
	}

	if slices.Contains(algorithms, "dinic") {
		dinicRouteSets, _, dinicTime := edgeNetwork.Dinic(edgeNetwork.NameToID[config.StartStation], edgeNetwork.NameToID[config.EndStation])
		fastestDinicRoutes := GetFastestRouteSet(dinicRouteSets, config.NumTrains)
		trainRoutes, maxTurnsDinic := divert(fastestDinicRoutes, config.NumTrains)

		result := Result{
			Name:        "dinic",
			TrainRoutes: trainRoutes,
			MaxTurns:    maxTurnsDinic,
			TimeTaken:   dinicTime,
		}

		content.Results = append(content.Results, result)
	}

	if slices.Contains(algorithms, "bellman") {
		mcmfRouteSets, _, bellFordTime := edgeNetwork.GetRoutesSets(edgeNetwork.NameToID[config.StartStation], edgeNetwork.NameToID[config.EndStation], p.BellmanFord)
		fastestMCMFRoutes := GetFastestRouteSet(mcmfRouteSets, config.NumTrains)
		trainRoutes, maxTurnsMCMF := divert(fastestMCMFRoutes, config.NumTrains)

		result := Result{
			Name:        "bellman",
			TrainRoutes: trainRoutes,
			MaxTurns:    maxTurnsMCMF,
			TimeTaken:   bellFordTime,
		}

		content.Results = append(content.Results, result)
	}

	if slices.Contains(algorithms, "dijkstra") {
		dijkstraRouteSets, _, dijkstraTime := edgeNetwork.GetRoutesSets(edgeNetwork.NameToID[config.StartStation], edgeNetwork.NameToID[config.EndStation], p.Dijkstra)
		fastestDijkstraRoutes := GetFastestRouteSet(dijkstraRouteSets, config.NumTrains)
		trainRoutes, maxTurnsDijkstra := divert(fastestDijkstraRoutes, config.NumTrains)

		result := Result{
			Name:        "dijkstra",
			TrainRoutes: trainRoutes,
			MaxTurns:    maxTurnsDijkstra,
			TimeTaken:   dijkstraTime,
		}

		content.Results = append(content.Results, result)
	}

	return content, nil
}
