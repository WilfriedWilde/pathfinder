package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

type networks interface {
	Network | FlowNetwork
}

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
}

type Edge struct {
	To       int
	Usable   bool
	Original bool
	Rev      int
	Cost     int
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
	config, err := GetAndValidateUserInput(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	network, err := ParseAndValidateNetwork(config)
	if err != nil {
		log.Fatal(err)
	}
	edgeNetwork := convertIntoFlow(&network, config.startStation, config.endStation)

	repeatedBFSRoutes, found := removeShortest(&network, config.startStation, config.endStation)
	if !found {
		log.Fatal(fmt.Errorf("No paths found between starting and ending stations."))
	}

	copyFlowNetwork := edgeNetwork.DeepCopy()
	mcmfRouteSets, _ := MCMF(copyFlowNetwork, edgeNetwork.NameToID[config.startStation], edgeNetwork.NameToID[config.endStation])

	ekRouteSets, _ := EK(&edgeNetwork, edgeNetwork.NameToID[config.startStation], edgeNetwork.NameToID[config.endStation])

	fastestMCMFRoutes := getFastestRouteSet(mcmfRouteSets, config.numTrains)
	fastestEKRoutes := getFastestRouteSet(ekRouteSets, config.numTrains)

	_, maxTurnsBFS := divert(repeatedBFSRoutes, config.numTrains)
	_, maxTurnsEK := divert(fastestEKRoutes, config.numTrains)
	trainRoutes, maxTurnsMCMF := divert(fastestMCMFRoutes, config.numTrains)

	fmt.Println("------------------------------------------")
	printTurns(trainRoutes, maxTurnsMCMF)
	fmt.Println("------------------------------------------")

	//createMap(network, RoutesUsed{config.startStation, config.endStation, fastestRoutes}, turnData, config.numTrains, true)

	timeTaken := time.Since(now)
	fmt.Println(timeTaken)
	fmt.Println("BFS Remove: ", maxTurnsBFS)
	fmt.Println("Edmonds-Karp: ", maxTurnsEK)
	fmt.Println("MCMF: ", maxTurnsMCMF)
}
