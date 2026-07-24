package internal

import (
	"fmt"
	p "pathfinder/pathfinding"
)

func tegLogic(f *p.FlowNetwork, maxTime, trains int) {
	tegNetwork := f.CreateTEGNetwork(maxTime)
	tegNetwork.AddNeighborEdges(f, maxTime)
	tegNetwork.AddSelfEdges(maxTime, 0, 1)

	routes, possible := tegNetwork.DoableInXTurns(f, 0, 1, trains, maxTime)
	if !possible {
		return
	}
	for _, route := range routes {
		fmt.Println(route)
	}

	// parents, endID, found := tegNetwork.FindRoute(0, 1, maxTime)
	// if !found {
	// 	return
	// }

	// path := tegNetwork.ExtractPath(parents, 0, endID)
	// for _, station := range path {
	// 	fmt.Println(f.IDToName[station.StationID], station.Time)
	// }

}
