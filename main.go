package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

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
}

type Edge struct {
	To       int
	Usable   bool
	Original bool
	Rev      int
	Cost     int
	Level    int
}

type TrainRoute struct {
	RouteLength int
	Route       []string
	Trains      int
}

type Config struct {
	networkPath  string
	startStation string
	endStation   string
	numTrains    int
}

func main() {
	now := time.Now()
	config, err := GetAndValidateUserInput(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	network, err := ParseAndValidateNetwork(config)
	if err != nil {
		log.Fatal(err)
	}

	network.TrimNetwork(config.startStation, config.endStation)

	edgeNetwork := network.ConvertIntoFlow(config.startStation, config.endStation)

	repeatedBFSRoutes, found := network.RemoveShortest(config.startStation, config.endStation)
	if !found {
		log.Fatal(fmt.Errorf("No paths found between starting and ending stations."))
	}

	ekNetwork := edgeNetwork.DeepCopy()
	dinicNetwork := edgeNetwork.DeepCopy()
	//mcmfNetwork := edgeNetwork.DeepCopy()

	ekRouteSets, _, ekTime := ekNetwork.GetRoutesSets(ekNetwork.NameToID[config.startStation], ekNetwork.NameToID[config.endStation], ekNetwork.BFSFlow)
	//mcmfRouteSets, _, bellFordTime := mcmfNetwork.GetRoutesSets(mcmfNetwork.NameToID[config.startStation], mcmfNetwork.NameToID[config.endStation], mcmfNetwork.BellmanFord)
	dinicRouteSets, _, dinicTime := dinicNetwork.Dinic(dinicNetwork.NameToID[config.startStation], dinicNetwork.NameToID[config.endStation])

	fastestEKRoutes := GetFastestRouteSet(ekRouteSets, config.numTrains)
	//fastestMCMFRoutes := GetFastestRouteSet(mcmfRouteSets, config.numTrains)
	fastestDinicRoutes := GetFastestRouteSet(dinicRouteSets, config.numTrains)

	_, maxTurnsBFS := divert(repeatedBFSRoutes, config.numTrains)
	_, maxTurnsEK := divert(fastestEKRoutes, config.numTrains)
	//trainRoutes, maxTurnsMCMF := divert(fastestMCMFRoutes, config.numTrains)
	_, maxTurnsDinic := divert(fastestDinicRoutes, config.numTrains)

	fmt.Println("------------------------------------------")
	//PrintTurns(trainRoutes, maxTurnsMCMF)
	fmt.Println("------------------------------------------")

	//createMap(network, RoutesUsed{config.startStation, config.endStation, fastestRoutes}, turnData, config.numTrains, true)

	timeTaken := time.Since(now)
	fmt.Println(timeTaken)
	fmt.Println("BFS Remove: ", maxTurnsBFS)
	fmt.Println("Edmonds-Karp: ", maxTurnsEK)
	//fmt.Println("MCMF: ", maxTurnsMCMF)
	fmt.Println("Dinic: ", maxTurnsDinic)
	fmt.Println("Edmonds-Karp time: ", ekTime)
	//fmt.Println("Bellman-Ford time: ", bellFordTime)
	fmt.Println("Dinic time: ", dinicTime)
}
