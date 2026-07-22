package internal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func GetStationNames(fileName string) ([]string, error) {
	stationNames := []string{}
	file, _ := os.Open(fileName)
	defer file.Close()

	readingStations := false

	scanner := bufio.NewScanner(file)
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
			readingStations = true
			continue
		}

		if line == "connections:" && readingStations {
			break
		}

		if readingStations {
			split := strings.Split(line, ",")

			stationName := strings.TrimSpace(split[0])

			stationNames = append(stationNames, stationName)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Scanner error reading %s: %s", fileName, err.Error())
	}

	return stationNames, nil
}
