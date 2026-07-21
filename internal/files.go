package internal

import (
	"bufio"
	"os"
	"strings"
)

func GetStationNames(fileName string) []string {
	ret := []string{}
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

			ret = append(ret, stationName)

		}

	}
	return ret
}
