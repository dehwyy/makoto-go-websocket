package ws

import (
	"fmt"
	"log"
	"net/http"
	"time"

	l "github.com/dehwyy/makoto-go-websocket.git/logger"
	"github.com/gorilla/websocket"
)

const (
	// time for write
	writeWaitTime = time.Second * 5
	// time of waiting either message or period
	messageWaitTime = time.Second * 10
	// send ping with period
	messagePingPeriod = messageWaitTime * 9 / 10
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Client struct {
	// the room to which the client is connected
	hub *Hub

	// ws connection
	conn *websocket.Conn

	// buffered channel to broadcast users' messages
	send chan []byte
}

func (c *Client) write() {

	// ticker which will ping the client in period
	ticker := time.NewTicker(messagePingPeriod)

	// stopping ticker and closing send-message chan at the end of the function
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			// set max write time
			c.conn.SetWriteDeadline(time.Now().Add(writeWaitTime))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			}

			// getting Writer
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)

			// iterating through all messages
			chanlen := len(c.send)
			for i := 0; i < chanlen; i++ {
				w.Write([]byte("\n"))
				w.Write(<-c.send)
			}

			// closing writer
			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:

			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Fatalf("Failed to ping: %v", err)
			}
		}
	}
}

func (c *Client) read() {
	defer func() {
		c.hub.unreg <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(messageWaitTime))

	// on messageType == ping
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(messageWaitTime))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()

		//
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected error: %v", err)
			}

			// If connection was aborted (for example: on browser close), then it would be printed here as an error
			l.Log(fmt.Sprintf("Error: %v, %v", err, c))
			break
		}
		// message = bytes.TrimSpace(bytes.Replace(message, newline, spp-ace, -1))
		c.hub.broadcast <- message
	}
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// initializing new client and registering it in the Hub
	c := &Client{
		hub:  hub,
		conn: ws,
		send: make(chan []byte, 256),
	}

	c.hub.reg <- c

	//
	go c.read()
	go c.write()
}
