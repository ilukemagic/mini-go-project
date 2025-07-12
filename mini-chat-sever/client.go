package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan Message
	username string
}

func (c *Client) readPump() {
	defer func() {
		log.Printf("Client %s disconnecting...", c.username)
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		var message Message
		err := c.conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading from client %s: %v", c.username, err)
			}
			break
		}

		message.Timestamp = time.Now()
		message.Sender = c.username
		if message.Type == "" {
			message.Type = ChatMessage
		}

		// send message to broadcast channel
		c.hub.broadcast <- message
		log.Printf("Message from %s sent to broadcast channel", c.username)
	}
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		message, ok := <-c.send
		if !ok {
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		err := c.conn.WriteJSON(message)
		if err != nil {
			log.Printf("Error writing to client %s: %v", c.username, err)
			return
		}
		log.Printf("Successfully wrote message to client %s", c.username)
	}
}
