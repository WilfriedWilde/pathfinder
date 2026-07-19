package server

import (
	"fmt"
	"math"
	"os"
	p "pathfinder/pathfinding"
	"sort"
)

type RoutesUsed struct {
	Start  string
	End    string
	Routes [][]string
}

func sx(x, minX int, s float64) float64 {
	return (float64(x - minX)) * s
}

func sy(y, minY int, s float64) float64 {
	return (float64(y - minY)) * s
}
func CreateMap(n p.Network, r RoutesUsed, turnData []string, trains int, labels bool) error {
	err := createSVG(n, r, trains)
	if err != nil {
		return err
	}
	err = createJSTurns(turnData)
	if err != nil {
		return err
	}
	startServer()
	return nil
}

func createSVG(n p.Network, r RoutesUsed, trains int) error {
	//create a [stationname] map that saves the color-to-be-used, empty if not part of found path
	colorMap := make(map[string]string)
	routeEdges := make(map[string]struct{}) // track used edges
	colorMap[r.Start] = "green"
	colorMap[r.End] = "red"
	for _, route := range r.Routes {
		for ind, routeStation := range route {
			if ind < len(route)-1 {
				routeEdges[route[ind]+"|"+route[ind+1]] = struct{}{}
				routeEdges[route[ind+1]+"|"+route[ind]] = struct{}{}
			}
			if ind == 0 || ind == len(route)-1 { // skip the already colored start/end stations
				continue
			}
			colorMap[routeStation] = "orange"
		}
	}

	//get min/max values of input coordinates -> get max height/width -> calculate density/scale
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

	density := math.Sqrt(float64(len(n.Stations)))
	width := maxX - minX
	height := maxY - minY
	canvas := 1000.0 //baseline that actual coordinates get scaled against
	scale := canvas / float64(max(width, height))
	radius := 6.0 / density * 10
	radius = math.Max(radius, 0.15)
	radius = math.Min(radius, 6.0)
	stroke := 2.0 / density * 10
	stroke = math.Max(stroke, 0.15)
	stroke = math.Min(stroke, 2.0)
	fontSize := 20.0
	labelXOffset := 12.0 //offsets for showing text/trains without overlapping with actual station
	labelYOffset := -7.0

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
			_, ok := routeEdges[startSStation.Name+"|"+endSStation.Name]
			if ok {
				fmt.Fprintf(f, `<line x1="%f" y1="%f" x2="%f" y2="%f" stroke="blue" stroke-width="%.2f"/>`+"\n", sx(startSStation.X, minX, scale), sy(startSStation.Y, minY, scale), sx(endSStation.X, minX, scale), sy(endSStation.Y, minY, scale), stroke)
			} else {
				fmt.Fprintf(f, `<line x1="%f" y1="%f" x2="%f" y2="%f" stroke="grey" stroke-width="%.2f"/>`+"\n", sx(startSStation.X, minX, scale), sy(startSStation.Y, minY, scale), sx(endSStation.X, minX, scale), sy(endSStation.Y, minY, scale), stroke)

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
		_, ok := colorMap[s.Name]
		if !ok {
			fmt.Fprintf(f, `<circle id="%s" cx="%f" cy="%f" r="%.2f" stroke="black"/>`+"\n", s.Name, sx(s.X, minX, scale), sy(s.Y, minY, scale), radius)
		} else {
			fmt.Fprintf(f, `<circle id="%s" cx="%f" cy="%f" r="%.2f" stroke="%s" fill="%s"/>`+"\n", s.Name, sx(s.X, minX, scale), sy(s.Y, minY, scale), radius, colorMap[s.Name], colorMap[s.Name])
		}
		fmt.Fprintf(f, `<text x="%.2f" y="%.2f" text-anchor="start" font-family="Arial" font-size="%.2f" fill="black" stroke="white" stroke-width="2.8" paint-order="stroke">%s</text>`+"\n", sx(s.X, minX, scale)+labelXOffset, sy(s.Y, minY, scale)+labelYOffset, fontSize, s.Name)
		fmt.Fprintln(f, `</g>`)
	}

	//initialize all required trains so they can be easily moved in js, also hide by default
	for i := 1; i <= trains; i++ {
		fmt.Fprintf(f, `<g class="train" id="T%d" visibility="hidden">`+"\n", i)
		fmt.Fprintf(f, `<circle cx="%f" cy="%f" r="%.2f" fill="brown" stroke="brown"/>`+"\n", 0., 0., radius*0.8)
		fmt.Fprintf(f, `<text x="%f" y="%f" font-family="Arial" font-size="%.2f" fill="black" stroke="white" stroke-width="2.5" paint-order="stroke">T%d</text>`+"\n", 0., 0., fontSize, i)
		fmt.Fprintln(f, `</g>`)

	}

	fmt.Fprintln(f, `</svg>`)
	return nil
}

func createJSTurns(turnData []string) error {
	f, err := os.Create("output/turns.js")
	if err != nil {
		return err
	}
	defer f.Close()
	fmt.Fprintln(f, "const turns = [")
	fmt.Fprintln(f, "	")
	for _, turn := range turnData {
		fmt.Fprintf(f, "	%q,\n", turn)
	}
	fmt.Fprintln(f, "];")
	return nil
}
