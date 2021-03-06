package sse

import (
	"fmt"
	"log"
	"net/http"
)

type Broker struct {
	clients map[chan string]bool
	newClients chan chan string
	closeClients chan chan string
	events chan string
}

func NewBroker() *Broker {
	return &Broker{
		make(map[chan string]bool),
		make(chan (chan string)),
		make(chan (chan string)),
		make(chan string),
	}
}

func (b *Broker) Start() {
	go func() {
		for {
			select {

			case client := <-b.newClients:
				b.addClient(client)

			case client := <-b.closeClients:
				b.removeClient(client)

			case ev := <-b.events:
				b.broadcastEvent(ev)
			}
		}
	}()
}

func (b *Broker) SendEvent(ev string) {
	b.events <- ev
}

func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	client := b.newClient()
	b.setupCloseNotifier(w, client)
	b.setupReponseWriter(w)
	b.sendEvents(w, f, client)

	log.Println("Finished HTTP request at ", r.URL.Path)
}

func (b *Broker) addClient(client chan string) {
	b.clients[client] = true
	log.Println("Added new client")

}

func (b *Broker) removeClient(client chan string) {
	delete(b.clients, client)
	close(client)
	log.Println("Removed client")
}

func (b *Broker) broadcastEvent(ev string) {
	for client, _ := range b.clients {
		client <- ev
	}
	log.Printf("Broadcast event to %d clients", len(b.clients))
}

func (b *Broker) setupCloseNotifier(w http.ResponseWriter, client chan string) {
	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		b.closeClients <- client
		log.Println("HTTP connection just closed.")
	}()
}

func (b *Broker) setupReponseWriter(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
}

func (b *Broker) sendEvents(w http.ResponseWriter, f http.Flusher, client chan string) {
	for {
		ev, open := <-client

		if !open {
			break
		}

		fmt.Fprintf(w, "data: Event: %s\n\n", ev)
		f.Flush()
	}
}

func (b *Broker) newClient() chan string{
	client := make(chan string)
	b.newClients <- client
	return client
}
