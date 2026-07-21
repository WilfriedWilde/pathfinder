package internal

import (
	p "pathfinder/pathfinding"
	"slices"
)

func RunLogic(args []string, algorithms []string) (PageContent, error) {
	config, err := GetAndValidateUserInput(args)
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
	}

	if NumLines > 25000 {
		content.GraphVisibility = false
	}

	network.TrimNetwork(config.StartStation, config.EndStation)

	edgeNetwork := network.ConvertIntoFlow(config.StartStation, config.EndStation)

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

	if content.GraphVisibility {
		routes := [][]string{}
		for _, v := range content.Results[0].TrainRoutes {
			routes = append(routes, v.Route)
		}
		turnData, err := PrintTurns(content.Results[0].TrainRoutes, content.Results[0].MaxTurns)
		if err != nil {
			return PageContent{}, err
		}
		CreateMap(network, RoutesUsed{config.StartStation, config.EndStation, routes}, turnData, config.NumTrains, false)
	}

	return content, nil
}
