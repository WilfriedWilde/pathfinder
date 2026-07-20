package internal

import (
	"bufio"
	"fmt"
	"os"
	p "pathfinder/pathfinding"
	"slices"
	"strconv"
	"strings"
)

// GetAndValidateUserInput returns a config file with the arguments parsed into
func GetAndValidateUserInput(args []string) (config Config, err error) {
	// extract args
	config.NetworkPath, config.StartStation, config.EndStation = args[0], args[1], args[2]
	numTrainsStr := args[3]

	// check if start station is the same as end station
	if config.StartStation == config.EndStation {
		return Config{}, fmt.Errorf("Error: start and end stations should be different")
	}

	// check if number of train is valid positive integer
	config.NumTrains, err = strconv.Atoi(numTrainsStr)
	if err != nil || config.NumTrains < 1 {
		return Config{}, fmt.Errorf("Error: number of train should be positive integer")
	}

	return
}

// ParseAndValidateNetwork builds and returns the basic (non-split) network
func ParseAndValidateNetwork(config Config) (p.Network, int, error) {
	network := p.Network{
		Stations:    make(map[string]p.Station),
		Connections: make(map[string][]string),
	}

	coordinates := make(map[string]struct{})

	// parse data from network map
	file, err := os.Open(config.NetworkPath)
	if err != nil {
		return p.Network{}, 0, fmt.Errorf("Error opening file: %v", err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lineNumber int
	var readingStations, readingConnections bool
	var hasStationsSection, hasConnectionsSection bool

	for scanner.Scan() {
		lineNumber++
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
				return p.Network{}, 0, fmt.Errorf("Error line %d: station data corrupted", lineNumber)

			}

			stationName, xStr, yStr := strings.TrimSpace(split[0]), strings.TrimSpace(split[1]), strings.TrimSpace(split[2])

			// check if station name is valid
			if strings.ContainsFunc(stationName, func(r rune) bool {
				return !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_')
			}) {
				return p.Network{}, 0, fmt.Errorf("Error line %d: station name not valid", lineNumber)
			}

			x, xErr := strconv.Atoi(xStr)
			if xErr != nil || x < 0 { // Accept 0 and positive integers to match testing examples
				return p.Network{}, 0, fmt.Errorf("Error line %d: x coordinate for station %s is not a valid positive integer", lineNumber, stationName)
			}
			y, yErr := strconv.Atoi(yStr)
			if yErr != nil || y < 0 { // Accept 0 and positive integers to match testing examples
				return p.Network{}, 0, fmt.Errorf("Error line %d: y coordinate for station %s is not a valid positive integer", lineNumber, stationName)
			}

			station := p.Station{
				Name: stationName,
				X:    x,
				Y:    y,
			}

			// check if coordinates already exist
			coords := fmt.Sprintf("%s,%s", xStr, yStr)
			if _, found := coordinates[coords]; found {
				return p.Network{}, 0, fmt.Errorf("Error line %d: coordinates %s are duplicated", lineNumber, coords)
			}

			coordinates[coords] = struct{}{}

			// check if station name is duplicated
			if _, found := network.Stations[stationName]; found {
				return p.Network{}, 0, fmt.Errorf("Error line %d: station name %s is duplicated", lineNumber, stationName)
			}

			network.Stations[stationName] = station
		} else if readingConnections {
			split := strings.Split(line, "-")
			if len(split) != 2 {
				return p.Network{}, 0, fmt.Errorf("Error line %d: connection data corrupted", lineNumber)
			}

			stationNameA, stationNameB := strings.TrimSpace(split[0]), strings.TrimSpace(split[1])

			// check if station name is valid
			if strings.ContainsFunc(stationNameA, func(r rune) bool {
				return !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_')
			}) {
				return p.Network{}, 0, fmt.Errorf("Error line %d: station name not valid", lineNumber)
			}

			// check if station name is valid
			if strings.ContainsFunc(stationNameB, func(r rune) bool {
				return !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_')
			}) {
				return p.Network{}, 0, fmt.Errorf("Error line %d: station name not valid", lineNumber)
			}

			// check if connections are already existing
			if len(network.Connections[stationNameA]) > 0 && slices.Contains(network.Connections[stationNameA], stationNameB) {
				return p.Network{}, 0, fmt.Errorf("Error line %d: a connection is already existing between %s and %s", lineNumber, stationNameA, stationNameB)
			}

			if len(network.Connections[stationNameB]) > 0 && slices.Contains(network.Connections[stationNameB], stationNameA) {
				return p.Network{}, 0, fmt.Errorf("Error line %d: a connection is already existing between %s and %s", lineNumber, stationNameB, stationNameA)
			}

			network.Connections[stationNameA] = append(network.Connections[stationNameA], stationNameB)
			network.Connections[stationNameB] = append(network.Connections[stationNameB], stationNameA)
		}
	}

	if err := scanner.Err(); err != nil {
		return p.Network{}, 0, fmt.Errorf("Scanner error: %v", err.Error())
	}

	if err := ValidateNetwork(network, hasStationsSection, hasConnectionsSection, config); err != nil {
		return p.Network{}, 0, err
	}

	return network, lineNumber, nil
}

// ValidateNetwork returns an error specifying why the network is invalid, if it is.
func ValidateNetwork(network p.Network, hasStationsSection, hasConnectionsSection bool, config Config) error {
	// check if connections are made with valid stations
	for stationName := range network.Connections {
		if _, found := network.Stations[stationName]; !found {
			return fmt.Errorf("Error: connection made with unexisting station %s", stationName)
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

	if _, found := p.BFS(&network, config.StartStation, config.EndStation); !found {
		return fmt.Errorf("Error: No paths found between starting and ending stations")
	}

	return nil
}
