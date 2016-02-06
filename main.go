package main

import (
	"log"
	"fmt"
	"time"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/alvalea/go-rest/rest"
	"github.com/alvalea/go-rest/sse"
)

func setupStatic(r *mux.Router) {
	r.PathPrefix("/web/").Handler(http.StripPrefix("/web/",
		http.FileServer(http.Dir("web"))))
}

func setupRest(r *mux.Router) {
	r.HandleFunc("/api/notes", rest.GetNoteHandler).Methods("GET")
}

func setupRouter() *mux.Router {
	r := mux.NewRouter()

	setupStatic(r)
	setupRest(r)

	return r
}

func setupBroker(r *mux.Router) *sse.Broker{
	b := sse.NewBroker()
	b.Start()
	r.Handle("/api/notes/events", b)
	return b
}

func testBroker(b *sse.Broker) {
	go func() {
		for i := 0; ; i++ {
			b.Messages <- fmt.Sprintf("%d - the time is %v", i, time.Now())

			log.Printf("Sent message %d ", i)
			time.Sleep(5 * 1e9)

		}
	}()
}

func main() {
	r := setupRouter()
	b := setupBroker(r)
	testBroker(b)
	http.ListenAndServe(":8080", r)
}
