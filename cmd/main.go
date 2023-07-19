package main

import (
	"log"
	"net/http"

	"github.com/dehwyy/makoto-go-websocket.git/ws"
)

func main() {
	// creating new instance of Hub and running it in the Goroutine
	hub := ws.NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWS(hub, w, r)
	})
	err := http.ListenAndServe(":7070", nil)
	if err != nil {
		log.Fatalf("Error occured when starting server op port %d: %v", 6969, err)
	}

}
