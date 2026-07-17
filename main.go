package main

import (
	"fmt"
	"log"
	"os"
	p "pathfinder/pathfinding"
	"time"
)

type Config struct {
	NetworkPath  string
	StartStation string
	EndStation   string
	NumTrains    int
}

func main() {
	now := time.Now()
	config, err := GetAndValidateUserInput(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	network, err := ParseAndValidateNetwork(config)
	if err != nil {
		log.Fatal(err)
	}

	network.TrimNetwork(config.StartStation, config.EndStation)

	edgeNetwork := network.ConvertIntoFlow(config.StartStation, config.EndStation)

	repeatedBFSRoutes, _, bfsTime := network.GreedyBFS(config.StartStation, config.EndStation)

	ekRouteSets, _, ekTime := edgeNetwork.GetRoutesSets(edgeNetwork.NameToID[config.StartStation], edgeNetwork.NameToID[config.EndStation], p.BFSFlow)
	//mcmfRouteSets, _, bellFordTime := edgeNetwork.GetRoutesSets(edgeNetwork.NameToID[config.StartStation], edgeNetwork.NameToID[config.EndStation], p.BellmanFord)
	dinicRouteSets, _, dinicTime := edgeNetwork.Dinic(edgeNetwork.NameToID[config.StartStation], edgeNetwork.NameToID[config.EndStation])
	dijkstraRouteSets, _, dijkstraTime := edgeNetwork.GetRoutesSets(edgeNetwork.NameToID[config.StartStation], edgeNetwork.NameToID[config.EndStation], p.Dijkstra)

	fastestEKRoutes := GetFastestRouteSet(ekRouteSets, config.NumTrains)
	//fastestMCMFRoutes := GetFastestRouteSet(mcmfRouteSets, config.NumTrains)
	fastestDinicRoutes := GetFastestRouteSet(dinicRouteSets, config.NumTrains)
	fastestDijkstraRoutes := GetFastestRouteSet(dijkstraRouteSets, config.NumTrains)
	for i, v := range fastestDinicRoutes {
		fmt.Printf("flow %d: paths: %+v", i, v)
	}
	fmt.Println("dijkstra:")
	for i, v := range fastestDijkstraRoutes {
		fmt.Printf("flow %d: paths: %+v", i, v)
	}

	_, maxTurnsBFS := divert(repeatedBFSRoutes, config.NumTrains)
	_, maxTurnsEK := divert(fastestEKRoutes, config.NumTrains)
	//trainRoutes, maxTurnsMCMF := divert(fastestMCMFRoutes, config.NumTrains)
	_, maxTurnsDinic := divert(fastestDinicRoutes, config.NumTrains)
	trainRoutes, maxTurnsDijkstra := divert(fastestDijkstraRoutes, config.NumTrains)

	fmt.Println("------------------------------------------")
	_, err = PrintTurns(trainRoutes, maxTurnsDijkstra)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("------------------------------------------")

	//createMap(network, RoutesUsed{config.StartStation, config.EndStation, fastestMCMFRoutes}, turnData, config.NumTrains, true)

	timeTaken := time.Since(now)
	fmt.Println("Total time: ", timeTaken)

	fmt.Println("BFS Remove: ", maxTurnsBFS)
	fmt.Println("Edmonds-Karp: ", maxTurnsEK)
	//fmt.Println("Bellford: ", maxTurnsMCMF)
	fmt.Println("Dijkstra: ", maxTurnsDijkstra)
	fmt.Println("Dinic: ", maxTurnsDinic)
	fmt.Println("BFS Remove: ", bfsTime)
	fmt.Println("Edmonds-Karp time: ", ekTime)
	//fmt.Println("Bellman-Ford time: ", bellFordTime)
	fmt.Println("Dinic time: ", dinicTime)
	fmt.Println("Dijkstra time: ", dijkstraTime)
}
