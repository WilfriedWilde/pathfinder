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
	Connections map[int][]Edge
}

type Edge struct {
	To   int
	Flow bool
	Rev  int
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

	routes, ok := removeShortest(&network, config.startStation, config.endStation)
	if !ok {
		log.Fatal(fmt.Errorf("No paths found between starting and ending stations."))
	}

	flowRoute, ok := BFSFlow(&flowNetwork, flowNetwork.NameToID[config.startStation], flowNetwork.NameToID[config.endStation])
	fmt.Print("Flowroute: ")
	fmt.Println(flowRoute)

	fmt.Println("routes:", routes, ok)
	fmt.Printf("network: %+v\nnumber of trains: %d\n", network, config.numTrains)

	// fmt.Println("routes:", routes, found)
	// fmt.Printf("network: %+v\nnumber of trains: %d\n", network, config.numTrains)

	trainRoutes := divert(routes, config.numTrains)
	for i, v := range trainRoutes {
		fmt.Printf("train route %d:\nrouteLength: %v | numTrains: %d\n\n", i, len(v.Route), v.Trains)
	}
	now := time.Now()
	printTurns(trainRoutes)
	timeTaken := time.Since(now)
	fmt.Println(timeTaken)

}
