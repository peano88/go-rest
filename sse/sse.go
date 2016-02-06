package sse

import (
	"fmt"
	"log"
	"net/http"
)

type Broker struct {
	clients map[chan string]bool
	addedClients chan chan string
	removedClients chan chan string
	Messages chan string
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

			case s := <-b.addedClients:
				b.clients[s] = true
				log.Println("Added new client")

			case s := <-b.removedClients:
				delete(b.clients, s)
				close(s)
				log.Println("Removed client")

			case msg := <-b.Messages:
				for s, _ := range b.clients {
					s <- msg
				}
				log.Printf("Broadcast message to %d clients", len(b.clients))
			}
		}
	}()
}

func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log.Println("Start HTTP request at ", r.URL.Path)

	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	messageChan := make(chan string)
	b.addedClients <- messageChan

	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		b.removedClients <- messageChan
		log.Println("HTTP connection just closed.")
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		msg, open := <-messageChan

		if !open {
			break
		}

		fmt.Fprintf(w, "data: Message: %s\n\n", msg)
		f.Flush()
	}

	log.Println("Finished HTTP request at ", r.URL.Path)
}
