// main.go
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// allow all CORS requests, in production environment, you need to check it more strictly
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// get username from query parameters
	username := r.URL.Query().Get("username")
	if username == "" {
		log.Println("Username is required")
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	// upgrade HTTP connection to WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade connection", http.StatusInternalServerError)
		return
	}

	// create new client
	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan Message, 256), // channel buffer size
		username: username,
	}

	// register client to hub
	client.hub.register <- client

	// start goroutines to handle read and write
	go client.writePump()
	go client.readPump()
}

func main() {
	// create new hub
	hub := NewHub()
	go hub.Run()

	// set WebSocket route
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	// static file service
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	// start server
	log.Printf("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
