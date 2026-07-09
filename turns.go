package main

import (
	"fmt"
	"strings"
)

func printTurns(trainRoutes []TrainRoute) {
	totalRoutes := len(trainRoutes)

	for i, route := range trainRoutes {
		fmt.Printf("route %d: RouteLength=%d Trains=%d\n", i, route.RouteLength, route.Trains)
	}

	// Total turns = length of smallest route (turns needed for one train to complete it)
	// plus (Trains - 1) extra turns for staggered starts, since each train departs
	// one turn after the previous one.
	totalTurns := trainRoutes[0].RouteLength + (trainRoutes[0].Trains - 1)

	var maxTrains int
	for _, route := range trainRoutes { // get the max number of trains across all routes to avoid out of range looping through it
		if route.Trains > maxTrains {
			maxTrains = route.Trains
		}
	}

	for turnNum := range totalTurns {
		var turn strings.Builder
		skipedTrains := 0

		for trainInd := range maxTrains {
			for routeInd, route := range trainRoutes {
				if trainInd >= route.Trains { // check if train index exceeds the number of train in current route
					skipedTrains++
					break
				}

				stationInd := turnNum - trainInd // calculate station indexs

				if !(stationInd < 0 || stationInd > route.RouteLength-1) { // check if train has already started or arrived, in other words if it's in movement
					label := trainInd*totalRoutes + routeInd + 1 - skipedTrains
					fmt.Fprintf(&turn, "T%d-%s ", label, route.Route[stationInd])
				}
			}
		}

		fmt.Println(turn.String())
	}
}
