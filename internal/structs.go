package internal

import (
	p "pathfinder/pathfinding"
	"time"
)

type Result struct {
	Name        string // algorithm name
	TrainRoutes []p.TrainRoute
	MaxTurns    int
	TimeTaken   time.Duration
}

type PageContent struct {
	GraphVisibility bool
	Results         []Result
	Routes          []string
	ErrorMessage    string
}

type Config struct {
	NetworkPath  string
	StartStation string
	EndStation   string
	NumTrains    int
	Server       bool
}
