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
	NetworkPath  string
	StartStation string
	EndStation   string
	NumTrains    int
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

	network.TrimNetwork(config.StartStation, config.EndStation)

	edgeNetwork := network.ConvertIntoFlow(config.StartStation, config.EndStation)

	repeatedBFSRoutes, found := network.RemoveShortest(config.StartStation, config.EndStation)
	if !found {
		log.Fatal(fmt.Errorf("No paths found between starting and ending stations."))
	}

	copyFlowNetwork := edgeNetwork.DeepCopy()

	mcmfRouteSets, _ := copyFlowNetwork.GetRoutesSets(copyFlowNetwork.NameToID[config.StartStation], copyFlowNetwork.NameToID[config.EndStation], copyFlowNetwork.BellmanFord)
	ekRouteSets, _ := edgeNetwork.GetRoutesSets(edgeNetwork.NameToID[config.StartStation], edgeNetwork.NameToID[config.EndStation], edgeNetwork.BFSFlow)

	fastestMCMFRoutes := GetFastestRouteSet(mcmfRouteSets, config.NumTrains)
	fastestEKRoutes := GetFastestRouteSet(ekRouteSets, config.NumTrains)

	_, maxTurnsBFS := divert(repeatedBFSRoutes, config.NumTrains)
	_, maxTurnsEK := divert(fastestEKRoutes, config.NumTrains)
	trainRoutes, maxTurnsMCMF := divert(fastestMCMFRoutes, config.NumTrains)

	fmt.Println("------------------------------------------")
	PrintTurns(trainRoutes, maxTurnsMCMF)
	fmt.Println("------------------------------------------")

	//createMap(network, RoutesUsed{config.StartStation, config.EndStation, fastestMCMFRoutes}, turnData, config.NumTrains, true)

	timeTaken := time.Since(now)
	fmt.Println(timeTaken)
	fmt.Println("BFS Remove: ", maxTurnsBFS)
	fmt.Println("Edmonds-Karp: ", maxTurnsEK)
	fmt.Println("MCMF: ", maxTurnsMCMF)
}
