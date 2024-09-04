package main

import (
	"github.com/sean9999/go-oracle"
	"github.com/sean9999/harebrain"
)

type peerRecord oracle.Peer

func (p peerRecord) Hash() string {
	return oracle.Peer(p).Nickname() + ".json"
}

func (p peerRecord) MarshalBinary() ([]byte, error) {
	return oracle.Peer(p).MarshalJSON()
}

func (p peerRecord) UnmarshalBinary(b []byte) error {
	var orc oracle.Peer
	orc.UnmarshalJSON(b)
	copy(p[:], orc[:])
	return nil
}

func loadAllPeerRecords() ([][]byte, error) {
	db := harebrain.NewDatabase()
	db.Open("data")
	kv, err := db.Table("peers").GetAll()
	if err != nil {
		return nil, err
	}
	rows := make([][]byte, 0, len(kv))
	for _, v := range kv {
		rows = append(rows, v)
	}
	return rows, nil
}

func loadAllPeers() ([]oracle.Peer, error) {
	records, err := loadAllPeerRecords()
	if err != nil {
		return nil, err
	}
	peers := make([]oracle.Peer, len(records))
	var pbuf oracle.Peer
	for i, rec := range records {
		err = pbuf.UnmarshalJSON(rec)
		if err != nil {
			return nil, err
		}
		peers[i] = pbuf
	}
	return peers, nil
}
