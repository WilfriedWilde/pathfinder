package main

import (
	"testing"
)

func TestGetAndValidateUserInput(t *testing.T) {
	tests := map[string]struct {
		input  []string
		output Config
		err    error
	}{
		"Vkhhliih": {
			input: []string{"", "", "", "", ""},
			output: Config{
				networkPath:  "",
				startStation: "",
				endStation:   "",
				numTrains:    0,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			config, err := GetAndValidateUserInput(test.input)

			if config != test.output || err != test.err {
				t.Errorf("got: %v %v, expected: %v %v", config, err, test.output, test.err)
			}
		})
	}
}

func TestParseAndValidateNetwork(t *testing.T) {

}
