package internal

import (
	"fmt"
	p "pathfinder/pathfinding"
	"slices"
)

// RunLogic handles one full run of the program with given arguments, returning a PageContent struct containing the results (routes, max turns, etc)
func RunLogic(args []string, algorithms []string) (PageContent, error) {
	if algorithms == nil {
		return PageContent{}, fmt.Errorf("No pathfinding algorithm selected.")
	}
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
		Summaries:       make([]Summary, 0),
	}

	if NumLines > 25000 { // don't do graph if it's too big for svg to handle
		content.GraphVisibility = false
	}

	network.TrimNetwork(config.StartStation, config.EndStation)

	edgeNetwork := network.ConvertIntoFlow(config.StartStation, config.EndStation)

	// run all selected algorithms, appending each result to the []Summaries
	if slices.Contains(algorithms, "dijkstra") {
		dijkstraRouteSets, _, dijkstraTime := edgeNetwork.GetRoutesSets(edgeNetwork.NameToID[config.StartStation], edgeNetwork.NameToID[config.EndStation], p.Dijkstra)
		fastestDijkstraRoutes := GetFastestRouteSet(dijkstraRouteSets, config.NumTrains)
		trainRoutes, maxTurnsDijkstra := divert(fastestDijkstraRoutes, config.NumTrains)
		summary := Summary{
			Name:        "dijkstra",
			TrainRoutes: trainRoutes,
			MaxTurns:    maxTurnsDijkstra,
			TimeTaken:   dijkstraTime,
		}

		content.Summaries = append(content.Summaries, summary)
	}

	if slices.Contains(algorithms, "bellman") {
		mcmfRouteSets, _, bellFordTime := edgeNetwork.GetRoutesSets(edgeNetwork.NameToID[config.StartStation], edgeNetwork.NameToID[config.EndStation], p.BellmanFord)
		fastestMCMFRoutes := GetFastestRouteSet(mcmfRouteSets, config.NumTrains)
		trainRoutes, maxTurnsMCMF := divert(fastestMCMFRoutes, config.NumTrains)

		summary := Summary{
			Name:        "bellman",
			TrainRoutes: trainRoutes,
			MaxTurns:    maxTurnsMCMF,
			TimeTaken:   bellFordTime,
		}

		content.Summaries = append(content.Summaries, summary)
	}

	if slices.Contains(algorithms, "dinic") {
		dinicRouteSets, _, dinicTime := edgeNetwork.Dinic(edgeNetwork.NameToID[config.StartStation], edgeNetwork.NameToID[config.EndStation])
		fastestDinicRoutes := GetFastestRouteSet(dinicRouteSets, config.NumTrains)
		trainRoutes, maxTurnsDinic := divert(fastestDinicRoutes, config.NumTrains)

		summary := Summary{
			Name:        "dinic",
			TrainRoutes: trainRoutes,
			MaxTurns:    maxTurnsDinic,
			TimeTaken:   dinicTime,
		}

		content.Summaries = append(content.Summaries, summary)
	}

	if slices.Contains(algorithms, "edmonds") {
		ekRouteSets, _, ekTime := edgeNetwork.GetRoutesSets(edgeNetwork.NameToID[config.StartStation], edgeNetwork.NameToID[config.EndStation], p.BFSFlow)
		fastestEKRoutes := GetFastestRouteSet(ekRouteSets, config.NumTrains)
		trainRoutes, maxTurnsEK := divert(fastestEKRoutes, config.NumTrains)

		summary := Summary{
			Name:        "edmonds",
			TrainRoutes: trainRoutes,
			MaxTurns:    maxTurnsEK,
			TimeTaken:   ekTime,
		}

		content.Summaries = append(content.Summaries, summary)
	}

	if slices.Contains(algorithms, "greedy") {
		repeatedBFSRoutes, _, bfsTime := network.GreedyBFS(config.StartStation, config.EndStation)
		trainRoutes, maxTurnsBFS := divert(repeatedBFSRoutes, config.NumTrains)

		summary := Summary{
			Name:        "greedy",
			TrainRoutes: trainRoutes,
			MaxTurns:    maxTurnsBFS,
			TimeTaken:   bfsTime,
		}

		content.Summaries = append(content.Summaries, summary)
	}

	routes := [][]string{}
	for _, v := range content.Summaries[0].TrainRoutes {
		routes = append(routes, v.Route)
	}

	if content.GraphVisibility && config.Server {
		turnData, err := PrintTurns(content.Summaries[0].TrainRoutes, content.Summaries[0].MaxTurns)
		if err != nil {
			return PageContent{}, err
		}
		CreateMap(network, RoutesUsed{config.StartStation, config.EndStation, routes}, turnData, config.NumTrains, false)
	} else {
		content.Routes = routes
	}

	return content, nil
}
