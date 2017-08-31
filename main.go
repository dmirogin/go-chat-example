package main

import (
	"net/http"
	"github.com/gorilla/websocket"
	"log"
	"flag"
	"strconv"
)

var (
	websocketPortPointer = flag.Int("port", 3000, "Port which will be used by websocket")
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// Chat websocket handler
func ChatHandler (hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(conn)
	hub.register <- client

	go client.Read(hub)
	go client.Write(hub)
}

// Client entry point
func IndexHandler (w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func main() {
	flag.Parse()

	hub := NewHub()
	go hub.Run()

	http.HandleFunc("/", IndexHandler)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ChatHandler(hub, w, r)
	})

	// File Server
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	log.Fatal(http.ListenAndServe("localhost:" + strconv.Itoa(*websocketPortPointer), nil))
}