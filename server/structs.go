package server

type ServerInit struct {
	Algorithms   []string            `json:"algorithms"`
	FileNames    []string            `json:"names"`
	FileStations map[string][]string `json:"stations"`
}
