package main

import (
	"log"
	"sync"
	"time"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
	messages   []Message    // 存储消息历史
	mutex      sync.RWMutex // 用于保护messages切片
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan Message, 256), // 增加缓冲区大小
		register:   make(chan *Client),
		unregister: make(chan *Client),
		messages:   make([]Message, 0),
	}
}

func (h *Hub) addMessage(message Message) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.messages = append(h.messages, message)
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

			// send system message to all clients
			joinMessage := Message{
				Type:      SystemMessage,
				Content:   client.username + " joined the chat",
				Timestamp: time.Now(),
			}
			h.addMessage(joinMessage)

			// send join message to all clients
			for c := range h.clients {
				select {
				case c.send <- joinMessage:
				default:
					log.Printf("Failed to send join message to client: %s", c.username)
				}
			}

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				leaveMessage := Message{
					Type:      SystemMessage,
					Content:   client.username + " left the chat",
					Timestamp: time.Now(),
				}
				h.addMessage(leaveMessage)

				// send leave message to all clients
				for c := range h.clients {
					select {
					case c.send <- leaveMessage:
						log.Printf("Sent leave message to client: %s", c.username)
					default:
						log.Printf("Failed to send leave message to client: %s", c.username)
					}
				}
			}

		case message := <-h.broadcast:
			// add message to history
			h.addMessage(message)
			successCount := 0

			// send message to all clients
			for client := range h.clients {
				select {
				case client.send <- message:
					successCount++
				default:
				}
			}
		}
	}
}

func (h *Hub) Unregister(client *Client) {
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)
	}
}
