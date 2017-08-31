package main

import (
	"net/http"
	"log"
	"flag"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/gotvitch/go-sse"
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

	http.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
		sseConnection, err := sse.Upgrade(w, r, sse.DefaultOptions)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for {
			select {
			case <-time.After(time.Second):
				sseConnection.Send("time", time.Now())
			case <-sseConnection.Closed:
				return
			}
		}
	})

	log.Fatal(http.ListenAndServe("localhost:" + strconv.Itoa(*websocketPortPointer), nil))
}