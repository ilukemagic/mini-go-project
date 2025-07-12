package main

import "time"

type MessageType string

const (
	ChatMessage   MessageType = "chat"
	SystemMessage MessageType = "system"
	StatusMessage MessageType = "status"
)

type Message struct {
	Type      MessageType `json:"type"`
	Content   string      `json:"content"`
	Sender    string      `json:"sender"`
	Timestamp time.Time   `json:"timestamp"`
}
