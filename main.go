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
	Levels      map[int]int
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

	LevelNetwork := edgeNetwork.DeepCopy()
	LevelNetwork.Dinic(LevelNetwork.NameToID[config.startStation], LevelNetwork.NameToID[config.endStation])

	copyFlowNetwork := edgeNetwork.DeepCopy()

	mcmfRouteSets, _ := copyFlowNetwork.GetRoutesSets(edgeNetwork.NameToID[config.startStation], edgeNetwork.NameToID[config.endStation], copyFlowNetwork.BFSFlow)
	ekRouteSets, _ := edgeNetwork.GetRoutesSets(edgeNetwork.NameToID[config.startStation], edgeNetwork.NameToID[config.endStation], edgeNetwork.BellmanFord)

	fastestMCMFRoutes := GetFastestRouteSet(mcmfRouteSets, config.numTrains)
	fastestEKRoutes := GetFastestRouteSet(ekRouteSets, config.numTrains)

	_, maxTurnsBFS := divert(repeatedBFSRoutes, config.numTrains)
	_, maxTurnsEK := divert(fastestEKRoutes, config.numTrains)
	trainRoutes, maxTurnsMCMF := divert(fastestMCMFRoutes, config.numTrains)

	fmt.Println("------------------------------------------")
	PrintTurns(trainRoutes, maxTurnsMCMF)
	fmt.Println("------------------------------------------")

	//createMap(network, RoutesUsed{config.startStation, config.endStation, fastestRoutes}, turnData, config.numTrains, true)

	timeTaken := time.Since(now)
	fmt.Println(timeTaken)
	fmt.Println("BFS Remove: ", maxTurnsBFS)
	fmt.Println("Edmonds-Karp: ", maxTurnsEK)
	fmt.Println("MCMF: ", maxTurnsMCMF)
}
