package internal

import (
	"fmt"
	p "pathfinder/pathfinding"
)

func tegLogic(f *p.FlowNetwork, maxTime, trains int) {
	tegNetwork := f.CreateTEGNetwork(maxTime)
	tegNetwork.AddSelfEdges(maxTime)
	tegNetwork.AddNeighborEdges(f, maxTime)
	fmt.Println("TEG nodes:", len(tegNetwork.NodeToID))
	fmt.Println("TEG connections:", len(tegNetwork.Connections))

	fmt.Println("start node:", tegNetwork.GetStationID(0, 0))
	start := tegNetwork.GetStationID(0, 0)

	fmt.Println("Outgoing from start@0")

	for _, e := range tegNetwork.Connections[start] {
		fmt.Println(
			tegNetwork.IDToNode[e.To],
			e.Usable,
		)
	}

	for t := 0; t <= maxTime; t++ {
		fmt.Println("end@", t, "=", tegNetwork.GetStationID(1, t))
	}
	parents, endID, found := tegNetwork.FindRoute(0, 1, maxTime)
	fmt.Println("FOUND END:")
	fmt.Println(endID)
	fmt.Println(tegNetwork.IDToNode[endID])
	if !found {
		return
	}

	path := tegNetwork.ExtractPath(parents, 0, endID)
	for _, station := range path {
		fmt.Println(tegNetwork.IDToName[station.StationID], station.Time)
	}
	for id, parent := range parents {
		if parent != -1 {
			fmt.Printf("%d <- %d\n", id, parent)
		}
	}

}
