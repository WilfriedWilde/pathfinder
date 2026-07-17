package main

import (
	"fmt"
	"log"
	"os"
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

	repeatedBFSRoutes, found := network.GreedyBFS(config.StartStation, config.EndStation)
	if !found {
		log.Fatal(fmt.Errorf("No paths found between starting and ending stations."))
	}

	ekNetwork := edgeNetwork.CopyFlowNetwork()
	dinicNetwork := edgeNetwork.CopyFlowNetwork()
	mcmfNetwork := edgeNetwork.CopyFlowNetwork()
	dijkstraNetwork := edgeNetwork.CopyFlowNetwork()

	ekRouteSets, _, ekTime := ekNetwork.GetRoutesSets(ekNetwork.NameToID[config.StartStation], ekNetwork.NameToID[config.EndStation], ekNetwork.BFSFlow)
	//mcmfRouteSets, _, bellFordTime := mcmfNetwork.GetRoutesSets(mcmfNetwork.NameToID[config.StartStation], mcmfNetwork.NameToID[config.EndStation], mcmfNetwork.BellmanFord)
	dinicRouteSets, _, dinicTime := dinicNetwork.Dinic(dinicNetwork.NameToID[config.StartStation], dinicNetwork.NameToID[config.EndStation])
	dijkstraRouteSets, _, dijkstraTime := dijkstraNetwork.GetRoutesSets(mcmfNetwork.NameToID[config.StartStation], mcmfNetwork.NameToID[config.EndStation], dijkstraNetwork.Dijkstra)

	fastestEKRoutes := GetFastestRouteSet(ekRouteSets, config.NumTrains)
	//fastestMCMFRoutes := GetFastestRouteSet(mcmfRouteSets, config.NumTrains)
	fastestDinicRoutes := GetFastestRouteSet(dinicRouteSets, config.NumTrains)
	fastestDijkstraRoutes := GetFastestRouteSet(dijkstraRouteSets, config.NumTrains)

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
	fmt.Println("Edmonds-Karp time: ", ekTime)
	//fmt.Println("Bellman-Ford time: ", bellFordTime)
	fmt.Println("Dinic time: ", dinicTime)
	fmt.Println("Dijkstra time: ", dijkstraTime)
}
