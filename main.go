package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("web"))
	r.Handle("/", fs)
	http.ListenAndServe(":8000", r)
}
