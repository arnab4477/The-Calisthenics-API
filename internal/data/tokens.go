package data

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"

	"github.com/arnab4477/Parkour_API/internal/validator"
)

// Constants for the tokens scope
const (
	ScopeActivation = "activation"
)

// Token struct to hold data for individual tokens
type Token struct {
	PlainText string
	Hash []byte
	User_id int64
	Expiry time.Time
	Scope string
}

// Function to generate tpkems
func generateToken(userId int64, ttl time.Duration, scope string) (*Token, error) {
	// Add the given data to the token struct
	token := &Token{
		User_id: userId,
		Expiry: time.Now().Add(ttl),
		Scope: scope,
	}

	// Initialize a 16 bytes long 0 value byte slice
	// and fill it with random bytes
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	// Encode the byte slice to a base-32 encoded string and add
	// that to the token's plain text field
	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	// generate a sha-256 hash of the token string
	hash := sha256.Sum256([]byte(token.PlainText))
	token.Hash = hash[:]

	return token, nil
}

// Function that validates the token's length
func ValidateTokenPlainText(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext == "", "token", "must be provided")
	v.Check(len(tokenPlaintext) < 26, "token", "must be 26 bytes long")
}

// The token model type that warps a database connection
type TokenModel struct {
	DB *sql.DB
}

// Function to add a token to the database
func (m TokenModel) InsertOneToken(token *Token) error {
	// SQL query for the insertion
	query := `
		INSERT INTO tokens (hash, expiry, user_id, scope)
		VALUES ($1, $2, $3, $4)`
	
	// Create the values slicr and execute the query
	args := []interface{}{token.Hash, token.Expiry, token.User_id, token.Scope}
	_, err := m.DB.Exec(query, args...)
	return err
}

// Shortcut function that creates a new token and inserts it into the database
func (m TokenModel) NewToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	// Generate the token
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	// Insert the token into the database
	err = m.InsertOneToken(token)
	return token, err
}

// Function to delete all tokens for a specific user and scope
func (m TokenModel) DeleteTokens(userID int64, scope string) error {
	// SQL query for the deletion
	query := `
		DELETE FROM tokens
		WHERE user_id = $1 AND scope = $2`

	// Execute the query
	_, err := m.DB.Exec(query, userID, scope)
	return err
}