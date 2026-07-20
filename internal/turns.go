package internal

import (
	"fmt"
	p "pathfinder/pathfinding"
	"strconv"
	"strings"
)

// PrintTurns prints and returns each turn, informing about which train are moving to which station.
//
// Build the output one train at a time, not one turn at a time.
func PrintTurns(trainRoutes []p.TrainRoute, maxTurns int) ([]string, error) {
	if maxTurns == -1 {
		return nil, fmt.Errorf("Something went wrong with pathfinding, no valid routes found.")
	}

	turns := make([]strings.Builder, maxTurns)
	turnCAP := 50 // caping fo visualization
	ret := make([]string, 0)

	var maxTrains int
	for _, route := range trainRoutes {
		if route.Trains > maxTrains {
			maxTrains = route.Trains
		}
	}

	trainLabel := 1
	for trainsSent := range maxTrains { // increase offset for turnNum calculation - trainsSent is the number of trains sent per route
		for _, route := range trainRoutes { // try to send a train on every route
			if trainsSent >= route.Trains {
				break
			}

			for stationID, station := range route.Route {
				turnNum := stationID + trainsSent                                                     //
				turns[turnNum].WriteString(fmt.Sprintf("T%s-%s ", strconv.Itoa(trainLabel), station)) // append train to the corresponding line
			}

			trainLabel++
		}
	}

	// store output for visualization and print it to terminal
	for _, turn := range turns {
		if maxTurns <= turnCAP {
			ret = append(ret, strings.TrimSpace(turn.String()))
		}

		fmt.Println(turn.String())
	}

	return ret, nil
}

// GetFastestRouteSet returns the set of distinct, non-overlapping routes that results in the fastest turn completion
func GetFastestRouteSet(routes [][][]string, trains int) [][]string {
	fastestID := 0
	_, solutionTurns := divert(routes[0], trains)
	for i := 1; i < len(routes); i++ {
		_, ph := divert(routes[i], trains)
		if ph < solutionTurns {
			fastestID = i
			solutionTurns = ph
		}
	}
	return routes[fastestID]
}
