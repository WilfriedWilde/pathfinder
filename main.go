package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
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
	Coordinates map[string]struct{}
}

func main() {
	network := Network{
		Stations:    make(map[string]Station),
		Connections: make(map[string][]string),
		Coordinates: make(map[string]struct{}),
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
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var readingStations, readingConnections bool
	var hasStationsSection, hasConnectionsSection bool

	for scanner.Scan() {
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
			hasStationsSection = true
			readingStations = true
			readingConnections = false
			continue
		} else if line == "connections:" {
			hasConnectionsSection = true
			readingStations = false
			readingConnections = true
			continue
		}

		if readingStations {
			split := strings.Split(line, ",")
			if len(split) != 3 {
				fmt.Println("Error: station data corrupted")
			}

			stationName, xStr, yStr := strings.TrimSpace(split[0]), strings.TrimSpace(split[1]), strings.TrimSpace(split[2])

			// check if station name is valid
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
				X:    x,
				Y:    y,
			}

			// check if coordinates already exist
			coords := fmt.Sprintf("%s,%s", xStr, yStr)
			if _, found := network.Coordinates[coords]; found {
				fmt.Printf("Error: coordinates %s are duplicated", coords)
				return
			}

			network.Coordinates[coords] = struct{}{}

			// check if station name is duplicated
			if _, found := network.Stations[stationName]; found {
				fmt.Printf("Error: station name %s is duplicated\n", stationName)
				return
			}

			network.Stations[stationName] = station
		} else if readingConnections {
			split := strings.Split(line, "-")
			if len(split) != 2 {
				fmt.Println("Error: connection data corrupted")
				return
			}

			stationNameA, stationNameB := strings.TrimSpace(split[0]), strings.TrimSpace(split[1])

			// check if station name is valid
			if strings.ContainsFunc(stationNameA, func(r rune) bool {
				return !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_')
			}) {
				fmt.Println("Error: station name not valid")
				return
			}

			// check if station name is valid
			if strings.ContainsFunc(stationNameB, func(r rune) bool {
				return !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_')
			}) {
				fmt.Println("Error: station name not valid")
				return
			}

			// check if connections are already existing
			if len(network.Connections[stationNameA]) > 0 && slices.Contains(network.Connections[stationNameA], stationNameB) {
				fmt.Printf("Error: a connection is already existing between %s and %s\n", stationNameA, stationNameB)
				return
			}

			if len(network.Connections[stationNameB]) > 0 && slices.Contains(network.Connections[stationNameB], stationNameA) {
				fmt.Printf("Error: a connection is already existing between %s and %s\n", stationNameA, stationNameB)
				return
			}

			network.Connections[stationNameA] = append(network.Connections[stationNameA], stationNameB)
			network.Connections[stationNameB] = append(network.Connections[stationNameB], stationNameA)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Scanner error:", err)
		return
	}

	// check if connections are made with valid stations
	for stationName := range network.Connections {
		if _, found := network.Stations[stationName]; !found {
			fmt.Printf("Error: connection made with unexisting station %s\n", stationName)
			return
		}
	}

	// check if network map includes both stations and connections sections
	if !hasStationsSection {
		fmt.Println("Error: network map doesn't include a 'stations' section")
	}

	if !hasConnectionsSection {
		fmt.Println("Error: network map doesn't include a 'connections' section")
	}

	if !hasStationsSection || !hasConnectionsSection {
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

	// check max number stations
	if len(network.Stations) > 10000 {
		fmt.Println("Error: the network map can only contains up to 10000 stations")
		return
	}

	// TODO: It displays "Error" on stderr when no path exists between the start and end stations.
}
