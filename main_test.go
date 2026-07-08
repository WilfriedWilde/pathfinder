package main

import (
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
			err:    "Error: too few arguments in command line\nUsage:\ngo run . [path to file containing network map] [start station] [end station] [number of trains]",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			config, err := GetAndValidateUserInput(test.input)

			if config != test.output || err.Error() != test.err {
				t.Errorf("got: %v %v, expected: %v %v", config, err.Error(), test.output, test.err)
			}
		})
	}
}

func TestParseAndValidateNetwork(t *testing.T) {

}
