package data

import (
	"database/sql"
	"errors"

	"github.com/arnab4477/Parkour_API/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

// Cystom password type to hold the user's password
type password struct {
	plain string
	hash []byte
}
// The User struct
type User struct {
	ID int64 `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password password `json:"-"`
	Activated bool `json:"activated"`
	Version int `json:"-"`
}


// function to hash user's password and store it in the password struct
func (p *password) SetHash(plainTextPassowrd string) error {
	// Hash the plainTextPassowrd using bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassowrd), 10)

	if err != nil {
		return err
	}

	// Set the values to the struct
	p.plain = plainTextPassowrd
	p.hash = hash

	return nil
}

 // Function to chek if the provided password matches its hash
 func (p *password) Matchhash(plainTextPassowrd string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainTextPassowrd))
	if err != nil {
		switch{
			// error if the hash does not match
			case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
				return false, nil
			default:
				return false, err
		}

	}
	return true, nil
 }

 // function to validate user's email address
 func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email == "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRegEx), "email", "must be a valid email address")
 }
 // function to validate user's plain text passwords
 func ValidatePlainPassword(v *validator.Validator, password string) {
	v.Check(password == "", "password", "must be provided")
	v.Check(len(password) <= 8, "password", "must be longer than 8 charaters")
	v.Check(len(password) >= 16, "password", "must be less than 16 charaters")
 }

 // Function to validate the user
 func ValidateUser(v *validator.Validator, user *User) {
	// Validate the user's username
	v.Check(user.Username == "", "username", "must be provided")
	v.Check(len(user.Username) <= 3, "username", "must be longer than 100 characters")
	v.Check(len(user.Username) >= 100, "username", "must be less than 100 characters")

	// Validate the user's email
	ValidateEmail(v, user.Email)

	// Validate the user's plain text password, if it is stored
	if user.Password.plain != "" {
		ValidatePlainPassword(v, user.Password.plain)
	}

	// User's password hash must never not be nil, if it ia that
	// means there is logic error in the code and a panic needs to be raised
	if user.Password.hash == nil {
		panic("user's password must never be nil")
	}
 }

// UserModel struct which warps a SQL connectopn pool
 type UserModel struct {
	DB *sql.DB
 }

 // Function to insert a new user record to the database
 func (m UserModel) InsertOneuser(user *User) error {
	// SQL query to insert an yser
	query := `
		INSERT INTO users (username, email, password_hash, activated)
		VALUES ($1, $2, $3, $4)
		RETURNING id, version`

	// Values to put into the database
	args := []interface{}{
		user.Username,
		user.Email,
		user.Password.hash,
		user.Activated,
	}
	
	// Execute the query. If an user provides a duplicate email then it will return
	// an error because of the UNIQYE constraint
	err := m.DB.QueryRow(query, args...).Scan(&user.ID, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil;
}

// Function to get one user by an  unique email
func (m UserModel) GetOneUserByEmail(email string) (*User, error) {
	// SQL query to retrieve one user from the database with email
	query := `
		SELECT * from Users
		WHERE email=$1`
	
	// An instance of the user struct
	var user User

	// Execute the query
	err := m.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.Version,
	)

	if err != nil {
		switch {
			// error in case of no records found
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

// Function to update one user
func (m UserModel) UpdateOneUser(user *User) error {
	// SQL query to update one user record in the database
	query := `
		UPDATE users
		SET username = $1, password_hash = $2, email = $3, activated = $4, version = version + 1
		WHERE id = $5 AND version = $6
		RETURNING version`

	args := []interface{}{
		user.Username,
		user.Password.hash,
		user.Email, user.Activated,
		user.ID,
		user.Version,
	}
	
	// Execute the SQLquery
	err := m.DB.QueryRow(query, args...).Scan(&user.Version)

	if err != nil {
		switch {
		// If the user triedd to update email to a duplicate one
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
}
return nil
}


