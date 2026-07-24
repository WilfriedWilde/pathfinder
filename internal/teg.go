package internal

import (
	"fmt"
	p "pathfinder/pathfinding"
	"strings"
)

func tegLogic(f *p.FlowNetwork, maxTime, trains int, start string) {
	tegNetwork := f.CreateTEGNetwork(maxTime)
	tegNetwork.AddNeighborEdges(f, maxTime)
	tegNetwork.AddSelfEdges(maxTime, 0, 1)

	for {
		tegNetwork := f.CreateTEGNetwork(maxTime)
		tegNetwork.AddNeighborEdges(f, maxTime)
		tegNetwork.AddSelfEdges(maxTime, 0, 1)
		tegNetwork.Potential = make([]int, len(tegNetwork.IDToName))
		_, possible := tegNetwork.DoableInXTurns(f, 0, 1, trains, maxTime)
		if !possible {
			maxTime++
			break
		}
		maxTime--
	}
	tegNetwork = f.CreateTEGNetwork(maxTime)
	tegNetwork.AddNeighborEdges(f, maxTime)
	tegNetwork.AddSelfEdges(maxTime, 0, 1)
	tegNetwork.Potential = make([]int, len(tegNetwork.IDToName))
	routes, _ := tegNetwork.DoableInXTurns(f, 0, 1, trains, maxTime)

	printTEGOutput(routes, start)
}

func printTEGOutput(routes [][]string, start string) {
	formatRoutes := make([][]string, len(routes))

	for lineInd, line := range routes {
		ph := []string{}
		for _, station := range line {
			phString := fmt.Sprintf("T%d-%s", lineInd+1, station)
			ph = append(ph, phString)

		}
		formatRoutes[lineInd] = ph
	}

	longestRow := 0
	for _, row := range routes {
		if len(row) > longestRow {
			longestRow = len(row)
		}
	}
	printOutput := make([][]string, 0)
	for col := 0; col < longestRow; col++ {
		ph := []string{}
		for row := 0; row < len(formatRoutes); row++ {
			if col < len(formatRoutes[row]) {
				ph = append(ph, formatRoutes[row][col])
			}

		}
		printOutput = append(printOutput, ph)
	}

	printPrintOutput := []string{}

	counter := 0

	for _, line := range printOutput {
		printLine := ""
		for _, trainObj := range line {
			trainObjSplit := strings.Split(trainObj, "-")
			if trainObjSplit[1] == start {
				continue
			}
			counter++
			printLine += trainObj + " "
		}
		printPrintOutput = append(printPrintOutput, printLine)
	}

	for _, line := range printPrintOutput {
		if line == "" {
			continue
		}
		fmt.Println(line)
	}
	fmt.Println(counter, " moves")
}
