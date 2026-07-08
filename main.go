package main

import (
	"fmt"
	"log"
)

type Station struct {
	Name string
	X, Y int
}

type Network struct {
	Stations    map[string]Station
	Connections map[string][]string
}

type TrainRoute struct {
	ID          int
	RouteLength int
	Route       []string
	Trains      int
}

func main() {
	networkPath, startStation, endStation, numTrains, err := getAndValidateUserInput()
	if err != nil {
		log.Fatal(err)
	}

	network, err := parseAndValidateNetwork(networkPath, startStation, endStation)
	if err != nil {
		log.Fatal(err)
	}

	routes, ok := removeShortest(&network, startStation, endStation)

	fmt.Println("routes:", routes, ok)

	fmt.Printf("network: %+v\nnumber of trains: %d\n", network, numTrains)

	trainRoutes := divert(routes, numTrains)
	for _, v := range trainRoutes {
		fmt.Println(v.ID, v.Route, v.RouteLength, v.Trains)
	}

	// TODO: It displays "Error" on stderr when no path exists between the start and end stations.
}
