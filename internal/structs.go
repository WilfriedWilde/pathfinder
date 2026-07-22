package internal

import (
	p "pathfinder/pathfinding"
	"time"
)

type Summary struct {
	Name        string // algorithm name
	TrainRoutes []p.TrainRoute
	MaxTurns    int
	TimeTaken   time.Duration
}

type PageContent struct {
	GraphVisibility bool
	Summaries       []Summary
	Routes          [][]string
	ErrorMessage    string
}

type Config struct {
	NetworkPath  string
	StartStation string
	EndStation   string
	NumTrains    int
	Server       bool
}

type RoutesUsed struct {
	Start  string
	End    string
	Routes [][]string
}

type RenderConfig struct {
	Scale        float64
	Radius       float64
	Stroke       float64
	FontSize     float64
	LabelXOffset float64
	LabelYOffset float64
}
