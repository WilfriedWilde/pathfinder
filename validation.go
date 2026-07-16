package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func GetAndValidateUserInput(args []string) (config Config, err error) {
	// check number command line arguments
	if len(args) < 5 {
		return config, fmt.Errorf("Error: too few arguments in command line\nUsage:\ngo run . [path to file containing network map] [start station] [end station] [number of trains]")
	}

	// extract args
	config.NetworkPath, config.StartStation, config.EndStation = args[1], args[2], args[3]
	NumTrainsStr := args[4]

	// check if start station is the same as end station
	if config.StartStation == config.EndStation {
		return config, fmt.Errorf("Error: start and end stations should be different")
	}

	// check if number of train is valid positive integer
	config.NumTrains, err = strconv.Atoi(NumTrainsStr)
	if err != nil || config.NumTrains < 1 {
		return config, fmt.Errorf("Error: number of train should be positive integer")
	}

	return
}

func ParseAndValidateNetwork(config Config) (Network, error) {
	network := Network{
		Stations:    make(map[string]Station),
		Connections: make(map[string][]string),
	}

	coordinates := make(map[string]struct{})

	// parse data from network map
	file, err := os.Open(config.NetworkPath)
	if err != nil {
		return Network{}, fmt.Errorf("Error opening file: %v", err.Error())
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
				return Network{}, fmt.Errorf("Error: station data corrupted\n")

			}

			stationName, xStr, yStr := strings.TrimSpace(split[0]), strings.TrimSpace(split[1]), strings.TrimSpace(split[2])

			// check if station name is valid
			if strings.ContainsFunc(stationName, func(r rune) bool {
				return !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_')
			}) {
				return Network{}, fmt.Errorf("Error: station name not valid")
			}

			x, xErr := strconv.Atoi(xStr)
			if xErr != nil || x < 0 { // Accept 0 and positive integers to match testing examples
				return Network{}, fmt.Errorf("Error: x coordinate for station %s is not a valid integer\n", stationName)
			}
			y, yErr := strconv.Atoi(yStr)
			if yErr != nil || y < 0 { // Accept 0 and positive integers to match testing examples
				return Network{}, fmt.Errorf("Error: y coordinate for station %s is not a valid integer\n", stationName)
			}

			station := Station{
				Name: stationName,
				X:    x,
				Y:    y,
			}

			// check if coordinates already exist
			coords := fmt.Sprintf("%s,%s", xStr, yStr)
			if _, found := coordinates[coords]; found {
				return Network{}, fmt.Errorf("Error: coordinates %s are duplicated", coords)
			}

			coordinates[coords] = struct{}{}

			// check if station name is duplicated
			if _, found := network.Stations[stationName]; found {
				return Network{}, fmt.Errorf("Error: station name %s is duplicated\n", stationName)
			}

			network.Stations[stationName] = station
		} else if readingConnections {
			split := strings.Split(line, "-")
			if len(split) != 2 {
				return Network{}, fmt.Errorf("Error: connection data corrupted")
			}

			stationNameA, stationNameB := strings.TrimSpace(split[0]), strings.TrimSpace(split[1])

			// check if station name is valid
			if strings.ContainsFunc(stationNameA, func(r rune) bool {
				return !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_')
			}) {
				return Network{}, fmt.Errorf("Error: station name not valid")
			}

			// check if station name is valid
			if strings.ContainsFunc(stationNameB, func(r rune) bool {
				return !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_')
			}) {
				return Network{}, fmt.Errorf("Error: station name not valid")
			}

			// check if connections are already existing
			if len(network.Connections[stationNameA]) > 0 && slices.Contains(network.Connections[stationNameA], stationNameB) {
				return Network{}, fmt.Errorf("Error: a connection is already existing between %s and %s\n", stationNameA, stationNameB)
			}

			if len(network.Connections[stationNameB]) > 0 && slices.Contains(network.Connections[stationNameB], stationNameA) {
				return Network{}, fmt.Errorf("Error: a connection is already existing between %s and %s\n", stationNameA, stationNameB)
			}

			network.Connections[stationNameA] = append(network.Connections[stationNameA], stationNameB)
			network.Connections[stationNameB] = append(network.Connections[stationNameB], stationNameA)
		}
	}

	if err := scanner.Err(); err != nil {
		return Network{}, fmt.Errorf("Scanner error: %v", err.Error())
	}

	if err := validateNetwork(network, hasStationsSection, hasConnectionsSection, config); err != nil {
		return Network{}, err
	}

	return network, nil
}

func validateNetwork(network Network, hasStationsSection, hasConnectionsSection bool, config Config) error {
	// check if connections are made with valid stations
	for stationName := range network.Connections {
		if _, found := network.Stations[stationName]; !found {
			return fmt.Errorf("Error: connection made with unexisting station %s\n", stationName)
		}
	}

	// check if network map includes both stations and connections sections
	if !hasStationsSection {
		return fmt.Errorf("Error: network map doesn't include a 'stations' section")
	}

	if !hasConnectionsSection {
		return fmt.Errorf("Error: network map doesn't include a 'connections' section")
	}

	// check if start and end stations exist
	_, foundStart := network.Stations[config.StartStation]
	_, foundEnd := network.Stations[config.EndStation]

	if !foundStart {
		return fmt.Errorf("Error: start station does not exist")
	}

	if !foundEnd {
		return fmt.Errorf("Error: end station does not exist")
	}

	// check max number stations
	if len(network.Stations) > 10000 {
		return fmt.Errorf("Error: the network map can only contains up to 10000 stations")
	}

	return nil
}
