package data

import (
	"database/sql"
	"errors"
)

// Custom errors
var (
	ErrNotFound = errors.New("record not found")
	ErrEditConflict = errors.New("record not found")
	ErrDuplicateEmail = errors.New("diplicate email")
)


// This Models struct all the models in the database
type Models struct {
	Movements MovementModel
	Users UserModel
}

// This returns the Models struct with the initialized movement model
func NewModels(db *sql.DB) Models {
	return Models {
		Movements: MovementModel{DB: db},
		Users: UserModel{DB: db},
	}
}