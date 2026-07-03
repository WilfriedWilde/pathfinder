package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
)

type Station struct {
	Name string
	X, Y int
}

type Graph struct {
	Stations    map[string]Station
	Connections map[string][]string
}

type Graph2 struct {
	Stations    map[uint16]Station
	Connections map[uint16][]uint16
}

type RoutesUsed struct {
	Start  string
	End    string
	Routes [][]string
	Exists bool
}

type Turns struct {
	TurnData []string
	Exists   bool
}

func main() {
	g := Graph{nil, nil}
	r := RoutesUsed{"", "", nil, false}
	t := Turns{nil, false}
	trains := 0
	labels := true
	createMap(g, r, t, trains, labels)
}

func createMap(g Graph, r RoutesUsed, t Turns, trains int, labels bool) {
	g = Graph{make(map[string]Station), make(map[string][]string)}
	g.Stations["beethoven"] = Station{"beethoven", 1, 6}
	g.Stations["verdi"] = Station{"verdi", 7, 1}
	g.Stations["albinoni"] = Station{"albinoni", 1, 1}
	g.Stations["handel"] = Station{"handel", 3, 14}
	g.Stations["mozart"] = Station{"mozart", 14, 9}
	g.Stations["part"] = Station{"part", 10, 0}
	g.Connections["beethoven"] = append(g.Connections["beethoven"], "handel", "albinoni", "verdi")
	g.Connections["verdi"] = append(g.Connections["verdi"], "beethoven", "part", "albinoni")
	g.Connections["albinoni"] = append(g.Connections["albinoni"], "verdi", "beethoven", "mozart")
	g.Connections["handel"] = append(g.Connections["handel"], "mozart", "beethoven")
	g.Connections["mozart"] = append(g.Connections["mozart"], "handel", "albinoni", "part")
	g.Connections["part"] = append(g.Connections["part"], "verdi", "mozart")
	//g = generateGrid(10, 10)
	r = RoutesUsed{"beethoven", "part", [][]string{{"verdi"}, {"handel", "mozart"}}, true}
	//t = Turns{[]string{"T1-station_011", "T1-station_021 T2-station_011", "T1-station_022 T2-station_021", "T1-station_032 T2-station_022", "T1-station_033 T2-station_032", "T1-station_034 T2-station_033", "T2-station_034"}, true}
	svgMap(g, r, t, 3, labels)
}

func sx(x, minX int, s float64) float64 {
	return (float64(x - minX)) * s
}

func sy(y, minY int, s float64) float64 {
	return (float64(y - minY)) * s
}

func svgMap(g Graph, r RoutesUsed, t Turns, trains int, labels bool) {
	//create a [stationname] map that saves the color-to-be-used, empty if not part of found path
	colorMap := make(map[string]string, 0)
	if r.Exists {
		colorMap[r.Start] = "green"
		colorMap[r.End] = "red"
		for _, v := range r.Routes {
			for _, v2 := range v {
				colorMap[v2] = "orange"
			}
		}
	}

	//get min/max values of input coordinates -> get max height/width -> calculate density/scale
	minX, maxX := math.MaxInt, math.MinInt
	minY, maxY := math.MaxInt, math.MinInt

	for _, s := range g.Stations {
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

	density := math.Sqrt(float64(len(g.Stations)))
	width := maxX - minX
	height := maxY - minY
	canvas := 1000.0 //baseline that actual coordinates get scaled against
	scale := canvas / float64(max(width, height))
	//	extent := max(width, height)
	radius := 6.0 / density * 10
	radius = math.Max(radius, 0.15)
	radius = math.Min(radius, 6.0)
	stroke := 2.0 / density * 10
	stroke = math.Max(stroke, 0.15)
	stroke = math.Min(stroke, 2.0)
	fontSize := 20.0
	labelXOffset := 12.0 //offsets for showing text/trains without overlapping with actual station
	labelYOffset := -7.0

	f, err := os.Create("map2.html")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	fmt.Fprintln(f, `<!DOCTYPE html>`)
	fmt.Fprintln(f, `<html>`)
	fmt.Fprintln(f, `<body>`)
	//generate viewbox with a buffer of 100 on every side compared to the svg map
	fmt.Fprintf(f, `<svg id="map" xmlns="http://www.w3.org/2000/svg" viewBox="%d %d %d %d">`+"\n", -100, -100, 1200, 1200)

	fmt.Fprintln(f, `<style>`)
	//hide/show labels by default
	if !labels {
		fmt.Fprintln(f, `.station text {	display: none; }
.station:hover text {display: block; }
.train text {	display: none; }
.train:hover text {display: block; }
svg {
	width: 100vw;
	height: 100vh;
	display: block;
}`)
	}
	//if turns exist, css for turn buttons
	if t.Exists {
		fmt.Fprintln(f, `#controls {
	position: fixed;
	bottom: 20px;
	left: 50%;
	transform: translateX(-50%);
	display: flex;
	align-items: center;
	gap: 16px;
	font-family: Arial, sans-serif;
}
#controls button {
	font-size: 14px;
	cursor: center;
}
#turnLabel {
	min-width: 100px;
	text-align: center;
}`)
	}
	fmt.Fprintln(f, `</style>`)

	//loop through all stations, loop through connected stations, check if connection has already been drawn, draw connection if not
	usedConnection := make(map[string]struct{}, 0)
	for startS, endS := range g.Connections {
		startSStation := g.Stations[startS]
		for _, toS := range endS {
			endSStation := g.Stations[toS]
			_, ok1 := usedConnection[startSStation.Name+"|"+endSStation.Name]
			_, ok2 := usedConnection[endSStation.Name+"|"+startSStation.Name]
			if ok1 || ok2 {
				continue
			}
			usedConnection[startSStation.Name+"|"+endSStation.Name] = struct{}{}
			_, ok1 = colorMap[startSStation.Name]
			_, ok2 = colorMap[endSStation.Name]
			if ok1 && ok2 {
				fmt.Fprintf(f, `<line x1="%f" y1="%f" x2="%f" y2="%f" stroke="blue" stroke-width="%.2f"/>`+"\n", sx(startSStation.X, minX, scale), sy(startSStation.Y, minY, scale), sx(endSStation.X, minX, scale), sy(endSStation.Y, minY, scale), stroke)
			} else {
				fmt.Fprintf(f, `<line x1="%f" y1="%f" x2="%f" y2="%f" stroke="grey" stroke-width="%.2f"/>`+"\n", sx(startSStation.X, minX, scale), sy(startSStation.Y, minY, scale), sx(endSStation.X, minX, scale), sy(endSStation.Y, minY, scale), stroke)

			}
		}
	}

	//create and sort a slice by x position, so that rightmost elements get drawn first (so that labels can be drawn on top of them)
	stationSlice := make([]Station, 0, len(g.Stations))
	for _, s := range g.Stations {
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
		fmt.Fprintf(f, `<circle cx="%f" cy="%f" r="%.2f" fill="brown" stroke="brown"/>`+"\n", 0, 0, radius*0.8)
		fmt.Fprintf(f, `<text x="%f" y="%f" font-family="Arial" font-size="%.2f" fill="black" stroke="white" stroke-width="2.5" paint-order="stroke">T%d</text>`+"\n", 0, 0, fontSize, i)
		fmt.Fprintln(f, `</g>`)

	}

	fmt.Fprintln(f, `</svg>`)
	//controls for turns
	fmt.Fprintf(f, `<div id="controls">
	<button id="prev">Previous</button>
	<span id="turnLabel"></span>
	<button id="next">Next</button>
</div>`)
	//javascript
	fmt.Fprintln(f, `
<script>


//viewbox panning/zooming
const svg = document.getElementById("map");
const viewBox = svg.viewBox.baseVal;

let dragging = false;
let lastX = 0;
let lastY = 0;

svg.addEventListener("wheel", e => {
	e.preventDefault();
	const rect = svg.getBoundingClientRect();
	const mouseX = (e.clientX - rect.left) / rect.width;
	const mouseY = (e.clientY - rect.top) / rect.height;
	const oldWidth = viewBox.width;
	const oldHeight = viewBox.height;
	const factor = 1.2;

	if (e.deltaY < 0) {
		viewBox.width /= factor;
		viewBox.height /= factor;
	} else {
		viewBox.width *= factor;
		viewBox.height *= factor;
	}

	const widthDiff = oldWidth - viewBox.width;
	const heightDiff = oldHeight - viewBox.height;

	viewBox.x += widthDiff * mouseX;
	viewBox.y += heightDiff * mouseY;
});

svg.addEventListener("mousedown", e => {
	dragging = true;
	lastX = e.clientX;
	lastY = e.clientY;
});

window.addEventListener("mouseup", () => {
	dragging = false;
});

window.addEventListener("mousemove", e => {

	if (!dragging) {
		return;
	}

	const dx = e.clientX - lastX;
	const dy = e.clientY - lastY;

	viewBox.x -= dx * (viewBox.width / svg.clientWidth);
	viewBox.y -= dy * (viewBox.height / svg.clientHeight);

	lastX = e.clientX;
	lastY = e.clientY;
});`)

	//arrow key hotkeys for turns, convert turndata into js consts
	if t.Exists {
		fmt.Fprintf(f, `window.addEventListener("keydown", e => {
	if (e.key === "ArrowLeft") {
		e.preventDefault();
		document.getElementById("prev").click();
	} else if (e.key === "ArrowRight") {
		e.preventDefault();
		document.getElementById("next").click();
	}
});
const turns = [
`)
		for _, v := range t.TurnData {
			fmt.Fprintf(f, "%q,\n", v)
		}
		fmt.Fprintln(f, `];`)
		//update button state
		fmt.Fprintf(f, `let currentTurn = 0
function updateButtons() {
	document.getElementById("turnLabel").textContent =
		"Turn " + (currentTurn + 1) + " / " + turns.length;
	document.getElementById("prev").disabled = currentTurn === 0;
	document.getElementById("next").disabled = currentTurn === turns.length - 1;
}
`)
		//hide all trains, unhide trains that are active on current turn, relocate trains
		fmt.Fprintf(f, `function showTurn(currentTurn) {
	document.querySelectorAll(".train").forEach(train => {
		train.setAttribute("visibility", "hidden");
	});
	const entries = turns[currentTurn].split(" ");
	for (const entry of entries) {
		const [trainName, stationName] = entry.split("-");
		const station = document.getElementById(stationName);
		const train = document.getElementById(trainName);
		const x = station.cx.baseVal.value;
		const y = station.cy.baseVal.value;

		train.setAttribute("visibility", "visible");
		//x + offset + offset of offset (+12 +10), y (+10 -7)
		train.querySelector("circle").setAttribute("cx", x+10);
		train.querySelector("circle").setAttribute("cy", y+10);
		train.querySelector("text").setAttribute("x", x+22);
		train.querySelector("text").setAttribute("y", y+3);
		}
	}

	document.getElementById("prev").onclick = () => {
		if (currentTurn > 0) {
			currentTurn--;
			showTurn(currentTurn);
			updateButtons();
		}
	};
	document.getElementById("next").onclick = () => {
		if (currentTurn + 1 < turns.length) {
			currentTurn++;
			showTurn(currentTurn);
			updateButtons();
		}
	};
showTurn(currentTurn);
updateButtons();
`)
	}
	fmt.Fprintln(f, `</script>`)
	fmt.Fprintln(f, `</body>`)
	fmt.Fprintln(f, `</html>`)
}

// generate a grid where stations are 100 x/y apart + small random adjustment
// every station checks if it's not the last on its row/column, and if not, adds connections to those instead
func generateGrid(rows, cols int) Graph {
	rng := rand.New(rand.NewSource(62))

	g := Graph{make(map[string]Station), make(map[string][]string)}

	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			name := fmt.Sprintf("station_%03d", y*cols+x)
			px := x*100 + rng.Intn(31) - 15
			py := y*100 + rng.Intn(31) - 15
			g.Stations[name] = Station{name, px, py}
		}
	}

	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			curr := fmt.Sprintf("station_%03d", y*cols+x)
			if x+1 < cols {
				right := fmt.Sprintf("station_%03d", y*cols+x+1)
				g.Connections[curr] = append(g.Connections[curr], right)
			}
			if y+1 < rows {
				down := fmt.Sprintf("station_%03d", (y+1)*cols+x)
				g.Connections[curr] = append(g.Connections[curr], down)
			}
		}
	}

	return g
}
