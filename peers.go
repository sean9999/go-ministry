package main

import (
	"io"

	"github.com/sean9999/go-oracle"
)

func newPeer(randy io.Reader) oracle.Peer {
	return oracle.New(randy).AsPeer()
}
