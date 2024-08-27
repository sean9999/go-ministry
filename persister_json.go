package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sean9999/harebrain"
)

var ErrJsonPersister = errors.New("jsonPersister")

type jsonPersister struct {
	store *harebrain.Table
}

var _ Persister = (*jsonPersister)(nil)

func (j *jsonPersister) Delete(k key) error {
	return j.store.Delete(k.String())
}

func (j *jsonPersister) Get(k key) (Message, error) {
	raw, err := j.store.Get(k.String())
	if err != nil {
		return NilMessage, fmt.Errorf("%w: retreiving message. %w", ErrJsonPersister, err)
	}
	msg := new(Message)
	err = json.Unmarshal(raw, msg)
	if err != nil {
		return NilMessage, fmt.Errorf("%w: unmarshalling json. %w", ErrJsonPersister, err)
	}
	return *msg, nil
}

func (j *jsonPersister) Has(k key) bool {
	_, err := j.Get(k)
	return err == nil
}

func (j *jsonPersister) Set(val Message) {
	j.store.Insert(&val)
}

func NewJsonPersister(rootPath string) (*jsonPersister, error) {
	db := harebrain.NewDatabase()
	if err := db.Open("data"); err != nil {
		return nil, fmt.Errorf("%w: opening database. %w", ErrJsonPersister, err)
	}
	j := jsonPersister{
		store: db.Table("wal"),
	}
	return &j, nil
}
