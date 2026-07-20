package server

import (
	"fmt"
	"net/http"
)

func startServer() {

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))))

	http.Handle("/output/",
		http.StripPrefix("/output/",
			http.FileServer(http.Dir("output"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	fmt.Println("running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
