package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	network := Network{
		Stations: make(map[string]Station),
		Connections: make(map[string][]string),
	}
	// check number command line arguments
	args := os.Args
	if len(args) < 5 {
		fmt.Println("Error: too few arguments in command line\nUsage:\ngo run . [path to file containing network map] [start station] [end station] [number of trains]")
		return
	}

	// extract args
	networkMapPath, startStation, endStation, numTrainsStr := args[1], args[2], args[3], args[4]

	// check if number of train is valid positive integer
	numTrains, err := strconv.Atoi(numTrainsStr)
	if err != nil || numTrains < 1 {
		fmt.Println("Error: number of train should be positive integer")
		return
	}

	// parse data from network map
	file, err := os.Open(networkMapPath)
	if err != nil {
		fmt.Println("Error opening file:",  err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var readingStations, readingConnections bool 

	for scanner.Scan(){
		line := strings.TrimSpace(scanner.Text())
		hashtagIndex := strings.Index(line, "#")

		if hashtagIndex != -1 {
			line = strings.TrimSpace(line[:hashtagIndex])
			continue
		}

		if line == "" {
			continue
		}

		if line == "stations:" {
			readingStations = true
			continue
		} else if line == "connections:" {
			readingStations = false
			readingConnections = true
			continue
		}

		if readingStations {
			split := strings.Split(line, ",")
			if len(split) != 3 {
				fmt.Println("Error: station data corrupted")
			}

			stationName, xStr, yStr := split[0], split[1], split[2]
			
			if strings.ContainsFunc(stationName, func(r rune) bool {
				return !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_')
			}) {
				fmt.Println("Error: station name not valid")
				return
			}

			x, xErr := strconv.Atoi(xStr)
			if xErr != nil || x < 0 {
				fmt.Printf("Error: x coordinate for station %s is not a valid integer\n", stationName)
			}
			y, yErr := strconv.Atoi(yStr)
			if yErr != nil || y < 0 {
				fmt.Printf("Error: y coordinate for station %s is not a valid integer\n", stationName)
			}

			if xErr != nil || yErr != nil {
				return
			}

			station := Station{
				Name: stationName,
				X: x,
				Y: y,
			}

			network.Stations[stationName] = station
		} else if readingConnections {
			split := strings.Split(line, "-")
			if len(split) != 2 {
				fmt.Println("Error: connection data corrupted")
			}

			stationNameA, stationNameB := split[0], split[1]
			
			if strings.ContainsFunc(stationNameA, func(r rune) bool {
				return !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_')
			}) {
				fmt.Println("Error: station name not valid")
				return
			}

			if strings.ContainsFunc(stationNameB, func(r rune) bool {
				return !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_')
			}) {
				fmt.Println("Error: station name not valid")
				return
			}

			network.Connections[stationNameA] = append(network.Connections[stationNameA], stationNameB)
		}
	}
	
	if err := scanner.Err(); err != nil {
		fmt.Println("Scanner error:", err)
		return
	}

	fmt.Printf("%+v\n", network)

	// check if start and end stations exist
	_, foundStart := network.Stations[startStation]
	_, foundEnd := network.Stations[endStation]

	if !foundStart {
		fmt.Println("Error: start session does not exist")
	}

	if !foundEnd {
		fmt.Println("Error: end session does not exist")
	}

	if !foundStart || !foundEnd {
		return
	}

	// check if start station is the same as end station
	if startStation == endStation {
		fmt.Println("Error: start and end stations should be different")
		return
	}



}
