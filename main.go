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
	config, err := GetAndValidateUserInput(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	network, err := ParseAndValidateNetwork(config)
	if err != nil {
		log.Fatal(err)
	}

	flowNetwork := convertIntoFlow(&network, config.startStation, config.endStation)

	routes, found := removeShortest(&network, config.startStation, config.endStation)
	if !found {
		log.Fatal(fmt.Errorf("No paths found between starting and ending stations."))
	}

	flowRoutes, flow := EK(&flowNetwork, flowNetwork.NameToID[config.startStation], flowNetwork.NameToID[config.endStation])

	fmt.Printf("%+v\n", flowRoutes)
	fmt.Printf("Flow: %d\n", flow)
	fmt.Print("Flowroute: ")
	for i, v := range flowRoutes {
		fmt.Println(i)
		fmt.Println(v)
	}
	fastestID := 0
	_, solutionTurns := divert(flowRoutes[0], config.numTrains)
	for i := 1; i < len(flowRoutes); i++ {
		_, ph := divert(flowRoutes[i], config.numTrains)
		if ph < solutionTurns {
			fastestID = i
			solutionTurns = ph
		}
	}
	fmt.Printf("Fastest solution uses flow of %d\n", fastestID+1)

	fmt.Println("routes:", routes, found)
	fmt.Printf("network: %+v\nnumber of trains: %d\n", network, config.numTrains)

	// fmt.Println("routes:", routes, found)
	// fmt.Printf("network: %+v\nnumber of trains: %d\n", network, config.numTrains)

	trainRoutes, maxTurns := divert(routes, config.numTrains)
	for i, v := range trainRoutes {
		fmt.Printf("train route %d:\nrouteLength: %v | numTrains: %d\n\n", i, len(v.Route), v.Trains)
	}
	fmt.Println("max turns:", maxTurns)

	now := time.Now()
	printTurns(trainRoutes, maxTurns)
	timeTaken := time.Since(now)
	fmt.Println(timeTaken)

	trainRoutes2, maxTurns2 := divert(flowRoutes[fastestID], config.numTrains)
	for i, v := range trainRoutes2 {
		fmt.Printf("train route %d:\nrouteLength: %v | numTrains: %d\n\n", i, len(v.Route), v.Trains)
	}
	turnData := printTurns(trainRoutes2, maxTurns2)
	fmt.Println(trainRoutes[0].RouteLength + (trainRoutes[0].Trains - 1))
	fmt.Println(trainRoutes2[0].RouteLength + (trainRoutes2[0].Trains - 1))
	createMap(network, RoutesUsed{config.startStation, config.endStation, flowRoutes[fastestID]}, turnData, config.numTrains, true)

}
