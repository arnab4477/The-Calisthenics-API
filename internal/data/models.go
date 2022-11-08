package data

import (
	"database/sql"
	"errors"
)

// Error to return when a record (like a movement) is not found in the database using the Get() method
var (
	ErrNotFound = errors.New("record not found")
)

// This Models struct all thee models in the database
type Models struct {
	Movements MovementModel
}

// This returns the Models struct with the initialized movement model
func NewModels(db *sql.DB) Models {
	return Models {
		Movements: MovementModel{DB: db},
	}
}