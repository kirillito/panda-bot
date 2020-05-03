package db

import bolt "go.etcd.io/bbolt"

// Tx represents a BoltDB transaction
type Tx struct {
	*bolt.Tx
}

// Vacation retrieves a Vacation from the database with the given userId.
func (tx *Tx) Vacation(userId []byte) (*Vacation, error) {
	vacation := &Vacation{
		Tx: tx,
		UserId: userId,
	}

	return vacation, vacation.Load()
}
