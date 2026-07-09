package main

import (
	"sort"
)

func divert(routes [][]string, trains int) []TrainRoute {
	trainRoutes := make([]TrainRoute, 0)
	for _, route := range routes {
		trainRoutes = append(trainRoutes, TrainRoute{len(route) - 1, route[1:], 0})
	}

	// sort once, so we only need to loop once going from smallest to largest
	sort.Slice(trainRoutes, func(i, j int) bool {
		return trainRoutes[i].RouteLength < trainRoutes[j].RouteLength
	})

	// go through all routes, from shortest to longest, incrementing the current and any previous to match the next
	for i := 0; i < len(trainRoutes)-1; i++ {
		diff := trainRoutes[i+1].RouteLength - trainRoutes[i].RouteLength
		if diff == 0 {
			continue
		}
		groupSize := i + 1 // number of routes we hanve handled so far, so to increment them by two, we need to have groupsize * 2 trains to use
		needed := diff * groupSize

		// not enough trains to fully increment all routes to match the next longest route length
		if trains < needed {
			distributeTrains(trains, groupSize, trainRoutes)
			trains = 0
			break
		}

		// increment to match next longest route length
		for j := 0; j < groupSize; j++ {
			trainRoutes[j].Trains += diff
		}

		trains -= needed
	}

	if trains > 0 {
		distributeTrains(trains, len(trainRoutes), trainRoutes)
	}
	return trainRoutes
}

func distributeTrains(trains, groupSize int, trainRoutes []TrainRoute) {
	increment := trains / groupSize
	remainder := trains % groupSize

	for j := 0; j < groupSize; j++ {
		trainRoutes[j].Trains += increment
	}

	for j := 0; j < remainder; j++ {
		trainRoutes[j].Trains++
	}
}
