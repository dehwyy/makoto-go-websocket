package main

import (
	"fmt"
	"log"
	"net/http"

	l "github.com/dehwyy/makoto-go-websocket.git/logger"
	"github.com/dehwyy/makoto-go-websocket.git/ws"
)

func main() {
	// creating new instance of Hub and running it in the goroutine
	hubs := make([]*ws.Hub, 0)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		requestedHubUID := r.URL.Query().Get("hub")

		// looking for requested Hub whether it exists
		for _, hub := range hubs {
			if hub.UID == requestedHubUID {
				ws.ServeWS(hub, w, r)
				return
			}
		}

		// if not, create new hub
		hub := ws.NewHub(requestedHubUID)
		l.Log(fmt.Sprintf("Created new Hub: %s", hub.UID))

		// adding new hub to the list
		hubs = append(hubs, hub)

		//
		go hub.Run()

		//
		ws.ServeWS(hub, w, r)

	})
	err := http.ListenAndServe(":7070", nil)
	if err != nil {
		log.Fatalf("Error occured when starting server op port %d: %v", 6969, err)
	}

}
