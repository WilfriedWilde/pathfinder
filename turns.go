package main

import (
	"fmt"
	"strconv"
	"strings"
)

func PrintTurns(trainRoutes []TrainRoute, maxTurns int) ([]string, error) {
	if maxTurns == -1 {
		return nil, fmt.Errorf("Something went wrong with pathfinding, no valid routes found.")
	}
	turns := make([]strings.Builder, maxTurns)
	turnCAP := 50
	ret := make([]string, 0)

	var maxTrains int
	for _, route := range trainRoutes {
		if route.Trains > maxTrains {
			maxTrains = route.Trains
		}
	}

	label := 1
	for trainNum := range maxTrains {
		for _, route := range trainRoutes {
			if trainNum >= route.Trains {
				break
			}
			for stationInd, station := range route.Route {
				turnNum := stationInd + trainNum
				turns[turnNum].WriteString(fmt.Sprintf("T%s-%s ", strconv.Itoa(label), station))
			}
			label++
		}
	}

	for _, turn := range turns {
		if maxTurns <= turnCAP {
			ret = append(ret, strings.TrimSpace(turn.String()))
		}
		fmt.Println(turn.String())
	}
	return ret, nil
}

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
