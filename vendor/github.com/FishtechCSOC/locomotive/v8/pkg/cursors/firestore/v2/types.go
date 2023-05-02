package firestore

import (
	"time"
)

type CursorDocument struct {
	Collector  string    `firestore:"collector"`
	Customer   string    `firestore:"customer"`
	Instance   string    `firestore:"instance"`
	Shard      string    `firestore:"shard"`
	Value      string    `firestore:"value"`
	LastUpdate time.Time `firestore:"lastUpdate"`
	Expiration time.Time `firestore:"expiration"`
}
