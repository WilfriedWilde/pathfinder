package main

import (
	"fmt"
	"log"
	"os"
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

	routes, found := removeShortest(&network, config.startStation, config.endStation)
	if !found {
		log.Fatal("Error: no path exists between the start and end stations")
	}

	fmt.Println("routes:", routes, found)
	fmt.Printf("network: %+v\nnumber of trains: %d\n", network, config.numTrains)

	trainRoutes := divert(routes, config.numTrains)
	for _, v := range trainRoutes {
		fmt.Println(v.Route, v.RouteLength, v.Trains)
	}

	turns := getTurns(trainRoutes, config.numTrains)
	for _, t := range turns {
		fmt.Println(t)
	}
}
