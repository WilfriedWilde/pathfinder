package server

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

type ServerInit struct {
	Algorithms   []string            `json:"algorithms"`
	FileNames    []string            `json:"names"`
	FileStations map[string][]string `json:"stations"`
}
