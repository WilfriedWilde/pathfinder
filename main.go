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

	// TODO: It displays "Error" on stderr when no path exists between the start and end stations.
}
