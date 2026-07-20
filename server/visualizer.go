package server

import (
	"fmt"
	"math"
	"os"
	p "pathfinder/pathfinding"
	"sort"
	"strings"
)

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

func CreateMap(n p.Network, r RoutesUsed, turnData []string, trains int, labels bool) error {
	err := createSVG(n, r, trains)
	if err != nil {
		return err
	}

	turnData = addStartStationTurn(turnData, r.Start)

	err = createJSTurns(turnData)
	if err != nil {
		return err
	}

	return nil
}

// addStartStationTurn adds an artificial "turn 0" for all trains, to properly initialize it at start station for animation purposes
func addStartStationTurn(turnData []string, start string) []string {
	seen := make(map[string]struct{})
	initTurn := make(map[int][]string)

	// go through turns, when seeing new trainid, appending that id (string) to the turn's index in initTurn
	for ind, turn := range turnData {
		for _, trainPos := range strings.Split(turn, " ") {
			trainID := strings.Split(trainPos, "-")[0][1:] // "T15-stationName" so get the "T15" and skip the first index
			if _, ok := seen[trainID]; !ok {
				seen[trainID] = struct{}{}
				initTurn[ind] = append(initTurn[ind], trainID)
			}
		}
	}
	// go through turns again, appending an extra start station turn for all the trains on the turn before their first appearance
	for ind := range turnData {
		for _, v := range initTurn[ind+1] {
			turnData[ind] += " T" + v + "-" + start
		}
	}
	// add "turn zero", so start station turn for the trains that show up on the first turn
	turnZero := ""
	for _, v := range initTurn[0] {
		turnZero += "T" + v + "-" + start + " "
	}
	turnZero = strings.TrimSpace(turnZero)
	turnData = append([]string{turnZero}, turnData...)
	return turnData
}

func createSVG(n p.Network, r RoutesUsed, trains int) error {

	stationColors, edgeColors := buildRouteColors(r)

	//get min/max values of input coordinates -> get max height/width -> calculate density/scale
	minX, maxX, minY, maxY := getBounds(n)

	config := createConfig(len(n.Stations), maxX-minX, maxY-minY)

	f, err := os.Create("output/map.svg")
	if err != nil {
		return err
	}
	defer f.Close()

	//generate viewbox with a buffer of 100 on every side compared to the svg map
	fmt.Fprintf(f, `<svg id="map" xmlns="http://www.w3.org/2000/svg" viewBox="%d %d %d %d">`+"\n", -100, -100, 1200, 1200)

	//loop through all stations, loop through connected stations, check if connection has already been drawn, draw connection if not
	usedConnection := make(map[string]struct{}, 0)
	for startS, endS := range n.Connections {
		startSStation := n.Stations[startS]
		for _, toS := range endS {
			endSStation := n.Stations[toS]
			_, ok1 := usedConnection[startSStation.Name+"|"+endSStation.Name]
			_, ok2 := usedConnection[endSStation.Name+"|"+startSStation.Name]
			if ok1 || ok2 {
				continue
			}
			usedConnection[startSStation.Name+"|"+endSStation.Name] = struct{}{}
			edgeColor, ok := edgeColors[startSStation.Name+"|"+endSStation.Name]
			if ok {
				fmt.Fprintf(f, `<line x1="%f" y1="%f" x2="%f" y2="%f" stroke="%s" stroke-width="%.2f"/>`+"\n", scaleX(startSStation.X, minX, config.Scale), scaleY(startSStation.Y, minY, config.Scale), scaleX(endSStation.X, minX, config.Scale), scaleY(endSStation.Y, minY, config.Scale), edgeColor, config.Stroke)
			} else {
				fmt.Fprintf(f, `<line x1="%f" y1="%f" x2="%f" y2="%f" stroke="grey" stroke-width="%.2f"/>`+"\n", scaleX(startSStation.X, minX, config.Scale), scaleY(startSStation.Y, minY, config.Scale), scaleX(endSStation.X, minX, config.Scale), scaleY(endSStation.Y, minY, config.Scale), config.Stroke)

			}
		}
	}

	//create and sort a slice by x position, so that rightmost elements get drawn first (so that labels can be drawn on top of them)
	stationSlice := make([]p.Station, 0, len(n.Stations))
	for _, s := range n.Stations {
		stationSlice = append(stationSlice, s)
	}
	sort.Slice(stationSlice, func(i, j int) bool {
		return stationSlice[i].X > stationSlice[j].X
	})

	//loop through sorted stations, check if it's part of a route and color the drawn circle based on that
	//also add text label at offset
	for _, s := range stationSlice {
		fmt.Fprintln(f, `<g class="station">`)
		_, ok := stationColors[s.Name]
		if !ok {
			fmt.Fprintf(f, `<circle id="%s" cx="%f" cy="%f" r="%.2f" stroke="black"/>`+"\n", s.Name, scaleX(s.X, minX, config.Scale), scaleY(s.Y, minY, config.Scale), config.Radius)
		} else {
			fmt.Fprintf(f, `<circle id="%s" cx="%f" cy="%f" r="%.2f" stroke="%s" fill="%s"/>`+"\n", s.Name, scaleX(s.X, minX, config.Scale), scaleY(s.Y, minY, config.Scale), config.Radius, stationColors[s.Name], stationColors[s.Name])
		}
		fmt.Fprintf(f, `<text x="%.2f" y="%.2f" text-anchor="start" font-family="Arial" font-size="%.2f" fill="black" stroke="white" stroke-width="2.8" paint-order="stroke">%s</text>`+"\n", scaleX(s.X, minX, config.Scale)+config.LabelXOffset, scaleY(s.Y, minY, config.Scale)+config.LabelYOffset, config.FontSize, s.Name)
		fmt.Fprintln(f, `</g>`)
	}

	//initialize all required trains so they can be easily moved in js, also hide by default
	for i := 1; i <= trains; i++ {
		fmt.Fprintf(f, `<g class="train" id="T%d" visibility="hidden">`+"\n", i)
		fmt.Fprintf(f, `<circle cx="%f" cy="%f" r="%.2f" fill="brown" stroke="brown"/>`+"\n", 0., 0., config.Radius*0.8)
		fmt.Fprintf(f, `<text x="%f" y="%f" font-family="Arial" font-size="%.2f" fill="black" stroke="white" stroke-width="2.5" paint-order="stroke">T%d</text>`+"\n", 0., 0., config.FontSize, i)
		fmt.Fprintln(f, `</g>`)

	}

	fmt.Fprintln(f, `</svg>`)
	return nil
}

// buildRouteColors returns [stationname]stationcolor and [edge]color maps
func buildRouteColors(r RoutesUsed) (map[string]string, map[string]string) {
	stationColors := make(map[string]string)
	edgeColors := make(map[string]string)
	routeColors := []string{
		"blue",
		"green",
		"orange",
		"purple",
		"teal",
		"red",
	}

	stationColors[r.Start] = "green"
	stationColors[r.End] = "red"

	for i, route := range r.Routes {
		for ind, routeStation := range route {
			if ind < len(route)-1 { // color the edges of the used route
				edgeColors[route[ind]+"|"+route[ind+1]] = routeColors[i%6] // add both ways, because at this point I can't remember if this had been normalized earlier or not
				edgeColors[route[ind+1]+"|"+route[ind]] = routeColors[i%6]
			}
			if ind == 0 || ind == len(route)-1 { // skip the already colored start/end stations
				continue
			}
			stationColors[routeStation] = "orange"
		}
	}
	return stationColors, edgeColors
}

// getBounds returns the min/max x/y coordinates of the given graph
func getBounds(n p.Network) (int, int, int, int) {
	minX, maxX := math.MaxInt, math.MinInt
	minY, maxY := math.MaxInt, math.MinInt
	for _, s := range n.Stations {
		if s.X < minX {
			minX = s.X
		}
		if s.X > maxX {
			maxX = s.X
		}
		if s.Y < minY {
			minY = s.Y
		}
		if s.Y > maxY {
			maxY = s.Y
		}
	}
	return minX, maxX, minY, maxY
}

// createConfig returns a config struct with scalable values based on the input graph density/scale
func createConfig(stationCount, width, height int) RenderConfig {
	const canvas = 1000.0
	const fontSize = 22.0
	const labelXOffset = 10.0
	const labelYOffset = -7.0

	density := math.Sqrt(float64(stationCount))
	scale := canvas / float64(max(width, height))
	radius := math.Min(math.Max(6.0/density*10, 0.15), 6.0)
	stroke := math.Min(math.Max(2.0/density*10, 0.15), 2.0)

	return RenderConfig{
		Scale:        scale,
		Radius:       radius,
		Stroke:       stroke,
		FontSize:     fontSize,
		LabelXOffset: labelXOffset,
		LabelYOffset: labelYOffset,
	}
}

// normalize the x coordinate to svg's scale
func scaleX(x, minX int, s float64) float64 {
	return (float64(x - minX)) * s
}

// normalize the y coordinate to svg's scale
func scaleY(y, minY int, s float64) float64 {
	return (float64(y - minY)) * s
}

// createJSTurns creates const turns[] in turns.js so javascript receives the turn info
func createJSTurns(turnData []string) error {
	f, err := os.Create("output/turns.js")
	if err != nil {
		return err
	}
	defer f.Close()
	fmt.Fprintln(f, "const turns = [")
	for _, turn := range turnData {
		fmt.Fprintf(f, "	%q,\n", turn)
	}
	fmt.Fprintln(f, "];")
	return nil
}
