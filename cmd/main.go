package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dehwyy/makoto-go-websocket/config"
	l "github.com/dehwyy/makoto-go-websocket/logger"
	"github.com/dehwyy/makoto-go-websocket/ws"
)

func main() {
	config := config.New()

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
	err := http.ListenAndServe(fmt.Sprintf(":%s", config.Env.Port), nil)
	if err != nil {
		log.Fatalf("Error occured when starting server op port %d: %v", 6969, err)
	}

}
