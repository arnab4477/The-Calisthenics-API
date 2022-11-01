package data

import (
	"time"
)

// Declare the movement struct and the JSON alternative keys
type Movement struct {
	ID int64 `json:"id"`
	CreatedAt time.Time `json:"-"` // This field will not show up in the JSON response
	Name string `json:"name"`
	Description string `json:"description"`
	Image string `json:"image"`
	Tutorials []string `json:"tutorials"` // Array of helpful tutorial links (YouTube, blogs etc) for the movement
	Skilltype []string `json:"skilltype"` // The type the movement belongs to such as 'vault', 'climb' etc
	Muscles []string `json:"muscles"`
	Difficulty string `json:"difficulty"` // Beginner, Intermediate or Advance
	Equipments []string `json:"equipments"`
	Prerequisite []string `json:"prerequisite"`
	Version int32 `json:"version"` // Version will start at 1 and will be incremented each time the struct is updated
}