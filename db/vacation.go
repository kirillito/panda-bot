package db

import (
	"time"
	"encoding/json"
)

// Errors
var (
	ErrVacationNotFound = &Error{"user vacation data not found", nil}
	ErrNoUserId     = &Error{"no user id", nil}
)

// Vacation represents a Vacation data
// Id - string id of Slack user
// Data - JSON with the list of saved vacations
type Vacation struct {
	Tx   *Tx
	UserId []byte
	Type string
	DateStart time.Time
	DateEnd time.Time
}

func (vacation *Vacation) bucket() []byte {
	return []byte("Vacations")
}

func (vacation *Vacation) get() ([]byte, error) {
	data := vacation.Tx.Bucket(vacation.bucket()).Get(vacation.UserId)
	if data == nil {
		return nil, ErrVacationNotFound
	}

	return data, nil
}

// Load retrieves vacation data from the database.
func (vacation *Vacation) Load() error {
	data, err := vacation.get()
	if err != nil {
		return err
	}

	// convert JSON into vacation data
	err = json.Unmarshal(data, &vacation)
	if err != nil {
		return err
	}
	
	return nil
}

// Save commits the Vacation data to the database.
func (vacation *Vacation) Save() error {
	if len(vacation.UserId) == 0 {
		return ErrNoUserId
	}

	// convert vacation data to JSON
	encodedData, err := json.Marshal(vacation)
	if err != nil {
			return err
	}

	return vacation.Tx.Bucket(vacation.bucket()).Put(vacation.UserId, encodedData)
}
