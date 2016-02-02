package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/alvalea/go-rest/lib"
)

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/web/").Handler(http.StripPrefix("/web/",
		http.FileServer(http.Dir("web"))))
	r.HandleFunc("/api/notes", lib.GetNoteHandler).Methods("GET")
	http.ListenAndServe(":8080", r)
}
