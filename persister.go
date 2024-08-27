package main

import (
	"errors"

	"github.com/google/uuid"
)

var ErrNoRecord = errors.New("no such record")

type key = uuid.UUID

type Persister interface {
	Set(Message)
	Get(key) (Message, error)
	Has(key) bool
	Delete(key) error
}
