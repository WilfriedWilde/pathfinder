package server

import (
	"fmt"
	"net/http"
)

func StartServer() {
	mux := http.NewServeMux()

	mux.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))))
	mux.Handle("/output/",
		http.StripPrefix("/output/",
			http.FileServer(http.Dir("output"))))
	mux.HandleFunc("/", HandleRoot)

	fmt.Println("running on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}

func HandleRoot(w http.ResponseWriter, r *http.Request) {

}
