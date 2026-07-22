package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	i "pathfinder/internal"
	"slices"
)

func StartServer(serverInit ServerInit) {
	mux := http.NewServeMux()

	mux.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))))
	mux.Handle("/output/",
		http.StripPrefix("/output/",
			http.FileServer(http.Dir("output"))))
	mux.HandleFunc("/", HandleRoot)
	mux.HandleFunc("/api/init", serverInit.HandleInitData)

	fmt.Println("running on http://localhost:8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatal(err)
	}
}

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form:", err)
	}

	data := i.PageContent{}

	if len(r.Form) != 0 {
		args := make([]string, 5)
		args[0] = "networks/" + r.Form["file"][0]
		args[1] = r.Form["start"][0]
		args[2] = r.Form["end"][0]
		args[3] = r.Form["trains"][0]
		args[4] = "server"
		algorithms := r.Form["algorithms"]

		data, err = i.RunLogic(args, algorithms)
		if err != nil {
			data.ErrorMessage = "Error: No paths found between starting and ending stations"
		}
	}

	funcMap := template.FuncMap{
		"increment1": func(i int) int { return i + 1 },
	}

	tmp, err := template.New("base.html").Funcs(funcMap).ParseFiles("templates/base.html")
	if err != nil {
		log.Println("Error parsing template:", err)
	}

	tmp.Execute(w, data)
}

func (s *ServerInit) HandleInitData(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(s)
}

func (s *ServerInit) GetAllFilesData() error {
	root := os.DirFS("networks/")

	txtFiles, err := fs.Glob(root, "*.txt")
	if err != nil {
		return err
	}

	fileNames := make(map[string]string)

	s.FileNames = make([]string, 0)
	s.FileStations = make(map[string][]string, 0)

	for _, v := range txtFiles {
		fileNames[v] = "networks/" + v
	}

	s.FileStations = make(map[string][]string, 0)
	for key := range fileNames {
		s.FileNames = append(s.FileNames, key)
	}
	slices.Sort(s.FileNames)

	for _, name := range s.FileNames {
		statioNames, err := i.GetStationNames(fileNames[name])
		if err != nil {
			return err
		}
		s.FileStations[name] = append(s.FileStations[name], statioNames...)
	}

	return nil
}
