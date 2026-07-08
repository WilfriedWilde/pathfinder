package main

import (
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
		"Number of train is a not numeric": {
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
		"Invalid station name": {
			input: Config{
				networkPath:  "networks/invalid_station_name.txt",
				startStation: "waterloo",
				endStation:   "euston",
				numTrains:    20,
			},
			output: Network{},
			err:    "Error: station name not valid",
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
