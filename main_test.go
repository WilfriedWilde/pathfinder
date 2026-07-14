package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGetAndValidateUserInput(t *testing.T) {
	tests := map[string]struct {
		input  []string
		output Config
		err    string
	}{
		"Too few command-line arguments": {
			input:  []string{},
			output: Config{},
			err:    "Error: too few arguments in command line\nUsage:\ngo run . [path to file containing network map] [start station] [end station] [number of trains]",
		},
		"Start station and end station are the same": {
			input:  []string{"networks/1.txt", "waterloo", "waterloo", "10"},
			output: Config{},
			err:    "Error: start and end stations should be different",
		},
		"Number of train is 0": {
			input:  []string{"networks/1.txt", "waterloo", "st_pancras", "0"},
			output: Config{},
			err:    "Error: number of train should be positive integer",
		},
		"Number of train is negative": {
			input:  []string{"networks/1.txt", "waterloo", "st_pancras", "0"},
			output: Config{},
			err:    "Error: number of train should be positive integer",
		},
		"Number of train is not numeric": {
			input:  []string{"networks/1.txt", "waterloo", "st_pancras", "a"},
			output: Config{},
			err:    "Error: number of train should be positive integer",
		},
		"Valid args": {
			input: []string{"networks/1.txt", "waterloo", "st_pancras", "10"},
			output: Config{
				networkPath:  "networks/1.txt",
				startStation: "waterloo",
				endStation:   "st_pancras",
				numTrains:    10,
			},
			err: "",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			config, err := GetAndValidateUserInput(test.input)
			var resultError string
			if err != nil {
				resultError = err.Error()
			}

			if config != test.output || resultError != test.err {
				t.Errorf("got: %v %v, expected: %v %v", config, resultError, test.output, test.err)
			}
		})
	}
}

func TestParseAndValidateNetwork(t *testing.T) {
	tests := map[string]struct {
		input  Config
		output Network
		err    string
	}{
		"Too few fields in station data": {
			input: Config{
				networkPath:  "networks/testing/too_few_fields_station_data.txt",
				startStation: "waterloo",
				endStation:   "euston",
				numTrains:    20,
			},
			output: Network{},
			err:    "Error: station data corrupted",
		},
		"Too many fields in station data": {
			input: Config{
				networkPath:  "networks/testing/too_many_fields_station_data.txt",
				startStation: "waterloo",
				endStation:   "euston",
				numTrains:    20,
			},
			output: Network{},
			err:    "Error: station data corrupted",
		},
		"Invalid station name": {
			input: Config{
				networkPath:  "networks/testing/invalid_station_name.txt",
				startStation: "waterloo",
				endStation:   "euston",
				numTrains:    20,
			},
			output: Network{},
			err:    "Error: station name not valid",
		},
		"X coordinate is non numeric": {
			input: Config{
				networkPath:  "networks/testing/char_x_coordinate.txt",
				startStation: "waterloo",
				endStation:   "euston",
				numTrains:    20,
			},
			output: Network{},
			err:    "Error: x coordinate for station waterloo is not a valid positive integer",
		},
		"X coordinate is float": {
			input: Config{
				networkPath:  "networks/testing/float_x_coordinate.txt",
				startStation: "waterloo",
				endStation:   "euston",
				numTrains:    20,
			},
			output: Network{},
			err:    "Error: x coordinate for station waterloo is not a valid positive integer",
		},
		"X coordinate is negative": {
			input: Config{
				networkPath:  "networks/testing/negative_x_coordinate.txt",
				startStation: "waterloo",
				endStation:   "euston",
				numTrains:    20,
			},
			output: Network{},
			err:    "Error: x coordinate for station waterloo is not a valid positive integer",
		},
		"Y coordinate is non numeric": {
			input: Config{
				networkPath:  "networks/testing/char_y_coordinate.txt",
				startStation: "waterloo",
				endStation:   "euston",
				numTrains:    20,
			},
			output: Network{},
			err:    "Error: y coordinate for station waterloo is not a valid positive integer",
		},
		"Y coordinate is float": {
			input: Config{
				networkPath:  "networks/testing/float_y_coordinate.txt",
				startStation: "waterloo",
				endStation:   "euston",
				numTrains:    20,
			},
			output: Network{},
			err:    "Error: y coordinate for station waterloo is not a valid positive integer",
		},
		"Y coordinate is negative": {
			input: Config{
				networkPath:  "networks/testing/negative_y_coordinate.txt",
				startStation: "waterloo",
				endStation:   "euston",
				numTrains:    20,
			},
			output: Network{},
			err:    "Error: y coordinate for station waterloo is not a valid positive integer",
		},
		"Duplicated station name": {
			input: Config{
				networkPath:  "networks/testing/duplicate_stations.txt",
				startStation: "waterloo",
				endStation:   "euston",
				numTrains:    20,
			},
			output: Network{},
			err:    "Error: station name waterloo is duplicated",
		},
		"Duplicated coordinates": {
			input: Config{
				networkPath:  "networks/testing/duplicate_coordinates.txt",
				startStation: "waterloo",
				endStation:   "euston",
				numTrains:    20,
			},
			output: Network{},
			err:    "Error: coordinates 3,1 are duplicated",
		},
		"Too few fields in connection data": {
			input: Config{
				networkPath:  "networks/testing/too_few_fields_connection_data.txt",
				startStation: "waterloo",
				endStation:   "euston",
				numTrains:    20,
			},
			output: Network{},
			err:    "Error: connection data corrupted",
		},
		"Too many fields in connection data": {
			input: Config{
				networkPath:  "networks/testing/too_many_fields_connection_data.txt",
				startStation: "waterloo",
				endStation:   "euston",
				numTrains:    20,
			},
			output: Network{},
			err:    "Error: connection data corrupted",
		},
		"Invalid station name in connection": {
			input: Config{
				networkPath:  "networks/testing/invalid_station_name_connection.txt",
				startStation: "waterloo",
				endStation:   "euston",
				numTrains:    20,
			},
			output: Network{},
			err:    "Error: station name not valid",
		},
		"Duplicated connection AB": {
			input: Config{
				networkPath:  "networks/testing/duplicate_connection-AB.txt",
				startStation: "waterloo",
				endStation:   "euston",
				numTrains:    20,
			},
			output: Network{},
			err:    "Error: a connection is already existing between waterloo and victoria",
		},
		"Duplicated connection BA": {
			input: Config{
				networkPath:  "networks/testing/duplicate_connection-BA.txt",
				startStation: "waterloo",
				endStation:   "euston",
				numTrains:    20,
			},
			output: Network{},
			err:    "Error: a connection is already existing between victoria and waterloo",
		},
		"Valid output": {
			input: Config{
				networkPath:  "networks/1.txt",
				startStation: "waterloo",
				endStation:   "euston",
				numTrains:    20,
			},
			output: Network{
				Stations: map[string]Station{
					"euston":     {Name: "euston", X: 11, Y: 23},
					"st_pancras": {Name: "st_pancras", X: 5, Y: 15},
					"victoria":   {Name: "victoria", X: 6, Y: 7},
					"waterloo":   {Name: "waterloo", X: 3, Y: 1},
				},
				Connections: map[string][]string{
					"euston":     {"waterloo", "st_pancras"},
					"st_pancras": {"euston", "victoria"},
					"victoria":   {"waterloo", "st_pancras"},
					"waterloo":   {"victoria", "euston"},
				},
			},
			err: "",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			network, err := ParseAndValidateNetwork(test.input)
			var resultError string
			if err != nil {
				resultError = err.Error()
			}

			if !reflect.DeepEqual(network, test.output) || resultError != test.err {
				t.Errorf("got: %v %q, expected: %v %q", network, resultError, test.output, test.err)
			}
		})
	}
}

func TestValidateNetwork(t *testing.T) {
	tests := map[string]struct {
		network               Network
		hasStationsSection    bool
		hasConnectionsSection bool
		config                Config
		output                string
	}{
		"Connection to unexisting station": {
			network: Network{
				Stations: map[string]Station{
					"st_pancras": {Name: "st_pancras", X: 5, Y: 15},
					"victoria":   {Name: "victoria", X: 6, Y: 7},
					"waterloo":   {Name: "waterloo", X: 3, Y: 1},
				},
				Connections: map[string][]string{
					"euston":     {"waterloo", "st_pancras"},
					"st_pancras": {"euston", "victoria"},
					"victoria":   {"waterloo", "st_pancras"},
					"waterloo":   {"victoria", "euston"},
				},
			},
			hasStationsSection:    true,
			hasConnectionsSection: true,
			config: Config{
				networkPath:  "networks/1.txt",
				startStation: "waterloo",
				endStation:   "victoria",
				numTrains:    20,
			},
			output: "Error: connection made with unexisting station euston",
		},
		"No stations section": {
			network: Network{
				Stations: map[string]Station{
					"waterloo": {"waterloo", 3, 1},
					"victoria": {"victoria", 6, 7},
				},
				Connections: map[string][]string{
					"waterloo": {"victoria"},
					"victoria": {"waterloo"},
				},
			},
			hasStationsSection:    false,
			hasConnectionsSection: true,
			config: Config{
				networkPath:  "networks/testing/whatever.txt",
				startStation: "waterloo",
				endStation:   "victoria",
				numTrains:    20,
			},
			output: "Error: network map doesn't include a 'stations' section",
		},
		"No connections section": {
			network: Network{
				Stations: map[string]Station{
					"waterloo": {"waterloo", 3, 1},
					"victoria": {"victoria", 6, 7},
				},
				Connections: map[string][]string{
					"waterloo": {"victoria"},
					"victoria": {"waterloo"},
				},
			},
			hasStationsSection:    true,
			hasConnectionsSection: false,
			config: Config{
				networkPath:  "networks/testing/whatever.txt",
				startStation: "waterloo",
				endStation:   "victoria",
				numTrains:    20,
			},
			output: "Error: network map doesn't include a 'connections' section",
		},
		"Start section does not exist": {
			network: Network{
				Stations: map[string]Station{
					"waterloo": {"waterloo", 3, 1},
					"victoria": {"victoria", 6, 7},
				},
				Connections: map[string][]string{
					"waterloo": {"victoria"},
					"victoria": {"waterloo"},
				},
			},
			hasStationsSection:    true,
			hasConnectionsSection: true,
			config: Config{
				networkPath:  "networks/testing/whatever.txt",
				startStation: "euston",
				endStation:   "victoria",
				numTrains:    20,
			},
			output: "Error: start station does not exist",
		},
		"End section does not exist": {
			network: Network{
				Stations: map[string]Station{
					"waterloo": {"waterloo", 3, 1},
					"victoria": {"victoria", 6, 7},
				},
				Connections: map[string][]string{
					"waterloo": {"victoria"},
					"victoria": {"waterloo"},
				},
			},
			hasStationsSection:    true,
			hasConnectionsSection: true,
			config: Config{
				networkPath:  "networks/testing/whatever.txt",
				startStation: "waterloo",
				endStation:   "euston",
				numTrains:    20,
			},
			output: "Error: end station does not exist",
		},
		"More than 10000 stations": {
			network: Network{
				Stations: func() map[string]Station {
					stations := make(map[string]Station, 10003)
					stations["victoria"] = Station{}
					stations["waterloo"] = Station{}
					for i := range 10001 {
						stations[fmt.Sprintf("%d", i)] = Station{}
					}
					return stations
				}(),
				Connections: map[string][]string{
					"waterloo": {"victoria"},
					"victoria": {"waterloo"},
				},
			},
			hasStationsSection:    true,
			hasConnectionsSection: true,
			config: Config{
				networkPath:  "networks/testing/whatever.txt",
				startStation: "waterloo",
				endStation:   "victoria",
				numTrains:    20,
			},
			output: "Error: the network map can only contains up to 10000 stations",
		},
		"Valid": {
			network: Network{
				Stations: map[string]Station{
					"waterloo": {"waterloo", 3, 1},
					"victoria": {"victoria", 6, 7},
				},
				Connections: map[string][]string{
					"waterloo": {"victoria"},
					"victoria": {"waterloo"},
				},
			},
			hasStationsSection:    true,
			hasConnectionsSection: true,
			config: Config{
				networkPath:  "networks/testing/whatever.txt",
				startStation: "waterloo",
				endStation:   "victoria",
				numTrains:    20,
			},
			output: "",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := ValidateNetwork(test.network, test.hasStationsSection, test.hasConnectionsSection, test.config)
			var resultError string
			if err != nil {
				resultError = err.Error()
			}

			if resultError != test.output {
				t.Errorf("got: %s, expected: %s", resultError, test.output)
			}
		})
	}
}

func TestBFS(t *testing.T) {
	tests := map[string]struct {
		network *Network
		start   string
		end     string
		output  []string
		found   bool
	}{
		"Not found": {
			network: &Network{
				Stations: map[string]Station{
					"euston":     {Name: "euston", X: 11, Y: 23},
					"st_pancras": {Name: "st_pancras", X: 5, Y: 15},
					"victoria":   {Name: "victoria", X: 6, Y: 7},
					"waterloo":   {Name: "waterloo", X: 3, Y: 1},
					"whatever":   {},
				},
				Connections: map[string][]string{
					"euston":     {"waterloo", "st_pancras"},
					"st_pancras": {"euston", "victoria"},
					"victoria":   {"waterloo", "st_pancras"},
					"waterloo":   {"victoria", "euston"},
				},
			},
			start:  "waterloo",
			end:    "whatever",
			output: nil,
			found:  false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, found := BFS(test.network, test.start, test.end)
			if !reflect.DeepEqual(result, test.output) || found != test.found {
				t.Errorf("got: %+v %v, expected: %+v %v", result, found, test.output, test.found)
			}
		})
	}
}
