package data

import (
	"database/sql"
	"time"

	"github.com/arnab4477/Parkour_API/internal/validator"
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

// Check if the input data causes any validation error
func ValidateMovement(v *validator.Validator, input *Movement) {
	// Check for valudation errors and provide keys for error messages in case of an error
	v.Check(input.Name == "", "name", "must not be empty")
	v.Check(len(input.Name) >= 256, "name", "must not be over 256 bytes")
	
	v.Check(input.Description == "", "description", "must not be empty")
	v.Check(input.Image == "", "image", "must be provided")
	v.Check(input.Difficulty == "", "difficulty", "must be provided")

	v.Check(len(input.Tutorials) <= 0, "tutorial", "must be provided")
	v.Check(len(input.Skilltype) <= 0, "skilltype", "must be provided")
	v.Check(len(input.Muscles) <= 0, "muscles", "must be provided")

	v.Check(len(input.Equipments) <= 0, "equipments", "must be provided")
	v.Check(len(input.Difficulty) <= 0, "difficulty", "must be provided")
	
	// Check that the items in various slices are uniuye
	v.Check(!validator.IsUnique(input.Tutorials), "tutorials", "must not contain duplicate values")
	v.Check(!validator.IsUnique(input.Skilltype), "skilltype", "must not contain duplicate values")
	v.Check(!validator.IsUnique(input.Equipments), "equipment", "must not contain duplicate values")
	v.Check(!validator.IsUnique(input.Muscles), "muscles", "must not contain duplicate values")
	v.Check(!validator.IsUnique(input.Prerequisite), "prerequisites", "must not contain duplicate values")
}

// Define a MovementModel struct which warps a SQL connectopn pool
type MovementModel struct {
	DB *sql.DB
}

// Method for inserting a new movement to the movement table
func (m MovementModel) Insert(movement *Movement) error {
	return nil
}
// Method for getting a new movement to the movement table
func (m MovementModel) Get(id int64) (*Movement, error) {
	return nil, nil
}
// Method for updating a new movement to the movement table
func (m MovementModel) Update(movement *Movement) error {
	return nil
}
// Method for deleting a new movement to the movement table
func (m MovementModel) Delete(id int64) error {
	return nil
}