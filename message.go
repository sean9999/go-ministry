package main

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Message struct {
	ID       uuid.UUID       `json:"id,omitempty"`
	ThreadID uuid.UUID       `json:"thread_id,omitempty"`
	Subject  string          `json:"subject,omitempty"`
	Payload  map[string]any  `json:"payload,omitempty"`
	Conn     *websocket.Conn `json:"-"`
}

func NewMessage() *Message {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	msg := Message{
		ID: id,
	}
	return &msg
}

func (m *Message) Hash() string {
	return m.ID.String()
}

func (m *Message) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Message) UnmarshalBinary(p []byte) error {
	return json.Unmarshal(p, m)
}

var NilMessage Message
