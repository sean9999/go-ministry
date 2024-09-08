package graph

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var ErrMessage = errors.New("message")
var ErrSubject = fmt.Errorf("%w: subject", ErrMessage)
var ErrID = fmt.Errorf("%w: ID", ErrMessage)
var ErrPayload = fmt.Errorf("%w: payload", ErrMessage)

type Message struct {
	From     string          `json:"from,omitempty"`
	To       string          `json:"to,omitempty"`
	ID       *uuid.UUID      `json:"id"`
	ThreadID *uuid.UUID      `json:"thread_id,omitempty"`
	Subject  string          `json:"subject"`
	Payload  json.RawMessage `json:"payload,omitempty"`
	Conn     *websocket.Conn `json:"-"`
}

func (m Message) Valid() (bool, error) {
	var errs []error
	if m.Subject == "" {
		errs = append(errs, fmt.Errorf("%w: %q", ErrSubject, m.Subject))
	}
	if m.ID == nil {
		errs = append(errs, fmt.Errorf("%w: nil", ErrID))
	}
	if !json.Valid(m.Payload) {
		errs = append(errs, fmt.Errorf("%w: invalid JSON", ErrPayload))
	}
	if len(errs) > 0 {
		var err error
		for _, e := range errs {
			err = fmt.Errorf("%w, %w", err, e)
		}
		return false, err
	}
	return true, nil
}

func NewMessage() Message {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	msg := Message{
		ID: &id,
	}
	return msg
}

func (m *Message) SetPayload(p json.Marshaler) error {
	j, err := p.MarshalJSON()
	if err != nil {
		return err
	}
	m.Payload = j
	return nil
}

func (m *Message) GetPayload() any {
	var val any
	json.Unmarshal(m.Payload, &val)
	return val
}

func (m *Message) Reply() Message {
	r := NewMessage()
	r.To = m.From
	m.From = r.To
	if m.ThreadID == nil {
		r.ThreadID = m.ID
	} else {
		r.ThreadID = m.ThreadID
	}
	r.Subject = m.Subject
	return r
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
