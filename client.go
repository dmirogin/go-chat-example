package main

import (
	"github.com/gorilla/websocket"
	"log"
	"encoding/json"
)

type (
	TextMessage struct {
		Author int `json:"author"`
		Text string `json:"text"`
	}
	Client struct {
		id int
		conn *websocket.Conn
		send chan *TextMessage
	}
)

var (
	generateID = IDGenerator()
)

// Handler for any readable data from this client
func (client *Client) Read(hub *Hub) {
	defer func() {
		hub.unregister <- client
		client.conn.Close()
	}()
	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		hub.messages <- &TextMessage{client.id, string(message)}
	}
}

// Handler for any receivable data for this client
func (client *Client) Write(hub *Hub) {
	defer func() {
		hub.unregister <- client
		client.conn.Close()
	}()
	for {
		select {
			case message := <- client.send:
				log.Println("Message to client", client.id)
				w, err := client.conn.NextWriter(websocket.TextMessage)

				if err != nil {
					return
				}

				data, _ := json.Marshal(message)
				w.Write(data)

				if err := w.Close(); err != nil {
					return
				}
			}
	}
}

// Create new client
func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		generateID(),
		conn,
		make(chan *TextMessage),
	}
}

// IDs generating
func IDGenerator() func() int {
	counter := 0
	return func () int {
		counter += 1
		return counter
	}
}
