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
	Coordinates map[string]struct{}
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

	fmt.Println("network:", network, "number of trains:", numTrains)

	// TODO: It displays "Error" on stderr when no path exists between the start and end stations.
}
