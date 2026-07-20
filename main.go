package main

import (
	"log"
	"os"
	i "pathfinder/internal"
	s "pathfinder/server"
)

func main() {
	argsLength := len(os.Args[1:])
	content := i.PageContent{}
	var err error

	if argsLength == 1 && os.Args[1] == "-s" {
		// Server mode
		config := i.Config{
			Server: true,
			Algorithms: []string{
				"greedy",
				"edmonds",
				"dinic",
				"bellman",
				"dijkstra",
			},
		}

		config.GetAllFilesData()
		s.StartServer(config)
		return
	} else if argsLength < 4 {
		log.Fatal("Error: too few arguments in command line\nUsage:\ngo run . [path to file containing network map] [start station] [end station] [number of trains]")
	} else if argsLength == 4 {
		// Default mode
		content, err = i.RunLogic(os.Args[1:], []string{"dijkstra"})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("Error: too many arguments in command line\nUsage:\ngo run . [path to file containing network map] [start station] [end station] [number of trains]")
	}

	_, err = i.PrintTurns(content.Results[0].TrainRoutes, content.Results[0].MaxTurns)
	if err != nil {
		log.Fatal(err)
	}
}
