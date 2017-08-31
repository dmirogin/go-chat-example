package main

import (
	"net/http"
	"github.com/gorilla/websocket"
	"log"
)

type (
	Client struct {
		id int
		conn *websocket.Conn
		send chan []byte
	}
	Hub struct {
		clients map [int]*Client
		messages chan []byte
		register chan *Client
		unregister chan *Client
	}
)

// IDs generating
func counter() func() int {
	counter := 0
	return func () int {
		counter += 1
		return counter
	}
}

var (
	nextInt = counter()
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// Chat websocket server
func ChatHandler (hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{nextInt(), conn, make(chan []byte)}
	hub.register<-client

	// client read
	go func(hub  *Hub, conn *websocket.Conn){
		defer func() {
			hub.unregister <- client
			conn.Close()
		}()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					log.Printf("error: %v", err)
				}
				break
			}

			hub.messages<-message
		}
	}(hub, conn)

	// client write
	go func(client *Client, hub  *Hub, conn *websocket.Conn){
		defer func() {
			hub.unregister <- client
			conn.Close()
		}()
		for {
			select {
				case message := <- client.send:
					log.Println("Message to client", client.id)
					w, err := conn.NextWriter(websocket.TextMessage)

					if err != nil {
						return
					}

					w.Write(message)

					if err := w.Close(); err != nil {
						return
					}
			}
		}
	}(client, hub, conn)
}


func (hub *Hub) route() {
	for {
		select {
			case message := <- hub.messages:
				log.Println("Got message", message)
				for _, client := range hub.clients {
					client.send<-message
				}
			case client := <- hub.register:
				log.Println("Client connected", client.id)
				hub.clients[client.id] = client
			case client := <- hub.unregister:
				log.Println("Client disconnected", client.id)
				delete(hub.clients, client.id)
		}
	}
}

// Client entry point
func IndexHandler (w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func main() {
	hub := &Hub{
		make(map[int]*Client),
		make(chan []byte),
		make(chan *Client),
		make(chan *Client),
	}

	go hub.route()

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		ChatHandler(hub, w, r)
	})
	http.HandleFunc("/", IndexHandler)

	log.Fatal(http.ListenAndServe("localhost:3000", nil))
}