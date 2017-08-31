package main

import "log"

type (
	Hub struct {
		clients map [int]*Client
		messages chan *TextMessage
		register chan *Client
		unregister chan *Client
	}
)

func (hub *Hub) Run() {
	for {
		select {
			case message := <- hub.messages:
				log.Println("Got message", message.Text)
				for _, client := range hub.clients {
					client.send <- message
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

func NewHub() *Hub {
	return &Hub{
		make(map[int]*Client),
		make(chan *TextMessage),
		make(chan *Client),
		make(chan *Client),
	}
}