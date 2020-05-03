package db

import bolt "go.etcd.io/bbolt"

// Tx represents a BoltDB transaction
type Tx struct {
	*bolt.Tx
}

// Vacation retrieves a Vacation from the database with the given name.
func (tx *Tx) Vacation(name []byte) (*Vacation, error) {
	v := &Vacation{
		Tx: tx,
		Id: name,
	}

	return v, v.Load()
}
