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

// HandleRoot handles the mainpage, calling RunLogic() on form submissions, and executing the template with the results
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
			data.ErrorMessage = err.Error()
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

// GetAllFilesData gets all of the file names and the station names contained within, appending those to the ServerInit struct
func (s *ServerInit) GetAllFilesData() error {
	root := os.DirFS("networks/")

	txtFiles, err := fs.Glob(root, "*.txt")
	if err != nil {
		return err
	}

	s.FileNames = make([]string, 0)
	s.FileStations = make(map[string][]string, 0)

	for _, file := range txtFiles {
		s.FileNames = append(s.FileNames, file)

		stations, err := i.GetStationNames("networks/" + file)
		if err != nil {
			return fmt.Errorf("Error trying to fetch stations from %s: %v", file, err)
		}
		s.FileStations[file] = stations
	}
	slices.Sort(s.FileNames)

	return nil
}
