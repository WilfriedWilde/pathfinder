package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	i "pathfinder/internal"
)

func StartServer(config i.Config) {
	mux := http.NewServeMux()

	mux.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))))
	mux.Handle("/output/",
		http.StripPrefix("/output/",
			http.FileServer(http.Dir("output"))))
	mux.HandleFunc("/", HandleRoot)
	mux.HandleFunc("/api/init", config.HandleInitData)

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

	tmp, err := template.ParseFiles("templates/base.html")
	if err != nil {
		log.Println("Error parsing template:", err)
	}

	tmp.Execute(w, data)
}
