package main

import (
	"fmt"
	p "pathfinder/pathfinding"
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
				NetworkPath:  "networks/1.txt",
				StartStation: "waterloo",
				EndStation:   "st_pancras",
				NumTrains:    10,
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
		output p.Network
		err    string
	}{
		"Too few fields in station data": {
			input: Config{
				NetworkPath:  "networks/testing/too_few_fields_station_data.txt",
				StartStation: "waterloo",
				EndStation:   "euston",
				NumTrains:    20,
			},
			output: p.Network{},
			err:    "Error line 4: station data corrupted",
		},
		"Too many fields in station data": {
			input: Config{
				NetworkPath:  "networks/testing/too_many_fields_station_data.txt",
				StartStation: "waterloo",
				EndStation:   "euston",
				NumTrains:    20,
			},
			output: p.Network{},
			err:    "Error line 4: station data corrupted",
		},
		"Invalid station name": {
			input: Config{
				NetworkPath:  "networks/testing/invalid_station_name.txt",
				StartStation: "waterloo",
				EndStation:   "euston",
				NumTrains:    20,
			},
			output: p.Network{},
			err:    "Error line 6: station name not valid",
		},
		"X coordinate is non numeric": {
			input: Config{
				NetworkPath:  "networks/testing/char_x_coordinate.txt",
				StartStation: "waterloo",
				EndStation:   "euston",
				NumTrains:    20,
			},
			output: p.Network{},
			err:    "Error line 4: x coordinate for station waterloo is not a valid positive integer",
		},
		"X coordinate is float": {
			input: Config{
				NetworkPath:  "networks/testing/float_x_coordinate.txt",
				StartStation: "waterloo",
				EndStation:   "euston",
				NumTrains:    20,
			},
			output: p.Network{},
			err:    "Error line 4: x coordinate for station waterloo is not a valid positive integer",
		},
		"X coordinate is negative": {
			input: Config{
				NetworkPath:  "networks/testing/negative_x_coordinate.txt",
				StartStation: "waterloo",
				EndStation:   "euston",
				NumTrains:    20,
			},
			output: p.Network{},
			err:    "Error line 4: x coordinate for station waterloo is not a valid positive integer",
		},
		"Y coordinate is non numeric": {
			input: Config{
				NetworkPath:  "networks/testing/char_y_coordinate.txt",
				StartStation: "waterloo",
				EndStation:   "euston",
				NumTrains:    20,
			},
			output: p.Network{},
			err:    "Error line 4: y coordinate for station waterloo is not a valid positive integer",
		},
		"Y coordinate is float": {
			input: Config{
				NetworkPath:  "networks/testing/float_y_coordinate.txt",
				StartStation: "waterloo",
				EndStation:   "euston",
				NumTrains:    20,
			},
			output: p.Network{},
			err:    "Error line 4: y coordinate for station waterloo is not a valid positive integer",
		},
		"Y coordinate is negative": {
			input: Config{
				NetworkPath:  "networks/testing/negative_y_coordinate.txt",
				StartStation: "waterloo",
				EndStation:   "euston",
				NumTrains:    20,
			},
			output: p.Network{},
			err:    "Error line 4: y coordinate for station waterloo is not a valid positive integer",
		},
		"Duplicated coordinates": {
			input: Config{
				NetworkPath:  "networks/testing/duplicate_coordinates.txt",
				StartStation: "waterloo",
				EndStation:   "euston",
				NumTrains:    20,
			},
			output: p.Network{},
			err:    "Error line 5: coordinates 3,1 are duplicated",
		},
		"Duplicated station name": {
			input: Config{
				NetworkPath:  "networks/testing/duplicate_stations.txt",
				StartStation: "waterloo",
				EndStation:   "euston",
				NumTrains:    20,
			},
			output: p.Network{},
			err:    "Error line 5: station name waterloo is duplicated",
		},
		"Too few fields in connection data": {
			input: Config{
				NetworkPath:  "networks/testing/too_few_fields_connection_data.txt",
				StartStation: "waterloo",
				EndStation:   "euston",
				NumTrains:    20,
			},
			output: p.Network{},
			err:    "Error line 10: connection data corrupted",
		},
		"Too many fields in connection data": {
			input: Config{
				NetworkPath:  "networks/testing/too_many_fields_connection_data.txt",
				StartStation: "waterloo",
				EndStation:   "euston",
				NumTrains:    20,
			},
			output: p.Network{},
			err:    "Error line 10: connection data corrupted",
		},
		"Invalid station name in connection": {
			input: Config{
				NetworkPath:  "networks/testing/invalid_station_name_connection.txt",
				StartStation: "waterloo",
				EndStation:   "euston",
				NumTrains:    20,
			},
			output: p.Network{},
			err:    "Error line 11: station name not valid",
		},
		"Duplicated connection AB": {
			input: Config{
				NetworkPath:  "networks/testing/duplicate_connection-AB.txt",
				StartStation: "waterloo",
				EndStation:   "euston",
				NumTrains:    20,
			},
			output: p.Network{},
			err:    "Error line 11: a connection is already existing between waterloo and victoria",
		},
		"Duplicated connection BA": {
			input: Config{
				NetworkPath:  "networks/testing/duplicate_connection-BA.txt",
				StartStation: "waterloo",
				EndStation:   "euston",
				NumTrains:    20,
			},
			output: p.Network{},
			err:    "Error line 11: a connection is already existing between victoria and waterloo",
		},
		"Valid output": {
			input: Config{
				NetworkPath:  "networks/1.txt",
				StartStation: "waterloo",
				EndStation:   "euston",
				NumTrains:    20,
			},
			output: p.Network{
				Stations: map[string]p.Station{
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
		network               p.Network
		hasStationsSection    bool
		hasConnectionsSection bool
		config                Config
		output                string
	}{
		"Connection to unexisting station": {
			network: p.Network{
				Stations: map[string]p.Station{
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
				NetworkPath:  "networks/1.txt",
				StartStation: "waterloo",
				EndStation:   "victoria",
				NumTrains:    20,
			},
			output: "Error: connection made with unexisting station euston",
		},
		"No stations section": {
			network: p.Network{
				Stations: map[string]p.Station{
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
				NetworkPath:  "networks/testing/whatever.txt",
				StartStation: "waterloo",
				EndStation:   "victoria",
				NumTrains:    20,
			},
			output: "Error: network map doesn't include a 'stations' section",
		},
		"No connections section": {
			network: p.Network{
				Stations: map[string]p.Station{
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
				NetworkPath:  "networks/testing/whatever.txt",
				StartStation: "waterloo",
				EndStation:   "victoria",
				NumTrains:    20,
			},
			output: "Error: network map doesn't include a 'connections' section",
		},
		"Start section does not exist": {
			network: p.Network{
				Stations: map[string]p.Station{
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
				NetworkPath:  "networks/testing/whatever.txt",
				StartStation: "euston",
				EndStation:   "victoria",
				NumTrains:    20,
			},
			output: "Error: start station does not exist",
		},
		"End section does not exist": {
			network: p.Network{
				Stations: map[string]p.Station{
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
				NetworkPath:  "networks/testing/whatever.txt",
				StartStation: "waterloo",
				EndStation:   "euston",
				NumTrains:    20,
			},
			output: "Error: end station does not exist",
		},
		"More than 10000 stations": {
			network: p.Network{
				Stations: func() map[string]p.Station {
					stations := make(map[string]p.Station, 10003)
					stations["victoria"] = p.Station{}
					stations["waterloo"] = p.Station{}
					for i := range 10001 {
						stations[fmt.Sprintf("%d", i)] = p.Station{}
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
				NetworkPath:  "networks/testing/whatever.txt",
				StartStation: "waterloo",
				EndStation:   "victoria",
				NumTrains:    20,
			},
			output: "Error: the network map can only contains up to 10000 stations",
		},
		"Valid": {
			network: p.Network{
				Stations: map[string]p.Station{
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
				NetworkPath:  "networks/testing/whatever.txt",
				StartStation: "waterloo",
				EndStation:   "victoria",
				NumTrains:    20,
			},
			output: "",
		},
		"Path not found": {
			network: p.Network{
				Stations: map[string]p.Station{
					"euston":     {Name: "euston", X: 11, Y: 23},
					"st_pancras": {Name: "st_pancras", X: 5, Y: 15},
					"victoria":   {Name: "victoria", X: 6, Y: 7},
					"waterloo":   {Name: "waterloo", X: 3, Y: 1},
					"whatever":   {Name: "whatever", X: 3, Y: 2},
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
				NetworkPath:  "networks/testing/whatever.txt",
				StartStation: "waterloo",
				EndStation:   "whatever",
				NumTrains:    10,
			},
			output: "Error: No paths found between starting and ending stations",
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
		network *p.Network
		start   string
		end     string
		output  []string
		found   bool
	}{
		"Not found": {
			network: &p.Network{
				Stations: map[string]p.Station{
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
			result, found := p.BFS(test.network, test.start, test.end)
			if !reflect.DeepEqual(result, test.output) || found != test.found {
				t.Errorf("got: %+v %v, expected: %+v %v", result, found, test.output, test.found)
			}
		})
	}
}
