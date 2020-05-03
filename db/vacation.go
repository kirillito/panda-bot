package db

// Errors
var (
	ErrVacationNotFound = &Error{"vacation not found", nil}
	ErrNoVacationId     = &Error{"no vacation id", nil}
)

// Vacation represents a Vacation data
type Vacation struct {
	Tx   *Tx
	Id   []byte
	Data []byte
}

func (v *Vacation) bucket() []byte {
	return []byte("Vacations")
}

func (v *Vacation) get() ([]byte, error) {
	data := v.Tx.Bucket(v.bucket()).Get(v.Id)
	if data == nil {
		return nil, ErrVacationNotFound
	}

	return data, nil
}

// Load retrieves a page from the database.
func (v *Vacation) Load() error {
	data, err := v.get()
	if err != nil {
		return err
	}

	v.Data = data

	return nil
}

// Save commits the Vacation data to the database.
func (v *Vacation) Save() error {
	if len(v.Id) == 0 {
		return ErrNoVacationId
	}

	return v.Tx.Bucket(v.bucket()).Put(v.Id, v.Data)
}
