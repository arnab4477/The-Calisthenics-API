package data

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/arnab4477/Parkour_API/internal/validator"
	"github.com/lib/pq"
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
	Prerequisites []string `json:"prerequisite"`
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
	v.Check(!validator.IsUnique(input.Prerequisites), "prerequisites", "must not contain duplicate values")
}

// MovementModel struct which warps a SQL connectopn pool
type MovementModel struct {
	DB *sql.DB
}

// Method for getting all the movements from the database
func (m MovementModel) GetAllMovements(
	name string,
	difficulty string,
	skilltype []string, 
	muscles []string, 
	equipments []string,
	filters Filters) ([]*Movement, error) {

		// SQL query to get all the movements from the database
		// There is full text search implemented for the name of the movement
		// For documentation, visit: https://www.postgresql.org/docs/current/datatype-textsearch.html
		//The movements will be sorted according to the given parameter (if any)
		// The limit and offset handles the pagination of the returned data
		query := fmt.Sprintf(`
			SELECT * FROM movements
			WHERE (to_tsvector('english', name) @@ plainto_tsquery('english', $1) OR $1 = '')
			AND (LOWER(difficulty) = LOWER($2) OR $2 = '')
			AND (skilltype @> $3 OR $3 = '{}')
			AND (muscles @> $4 OR $4 = '{}')
			AND (equipments @> $5 OR $5 = '{}')
			order by %s %s, id ASC
			LIMIT %d OFFSET %d`, filters.sortColumns(), filters.sortDirection(),
								 filters.limit(), filters.offset())

		// Execute the SQL query
		rows, err := m.DB.Query(
			query, name, difficulty, pq.Array(skilltype),
			pq.Array(muscles), pq.Array(equipments))

		if err != nil {
			return nil, err
		}
		defer rows.Close()

		// The movements array that hold all the movements
		movements := []*Movement{}

		// Initialize an emoty movement struct and out the data in it
		// using rows.Next() to inerate over the rows
		for rows.Next() {
			var movement Movement

			err := rows.Scan(
				&movement.ID,
				&movement.CreatedAt,
				&movement.Name,
				&movement.Description,
				&movement.Image,
				pq.Array(&movement.Tutorials),
				pq.Array(&movement.Skilltype),
				pq.Array(&movement.Muscles),
				&movement.Difficulty,
				pq.Array(&movement.Equipments),
				pq.Array(&movement.Prerequisites),
				&movement.Version,
			)

			if err != nil {
				return nil, err
			}

			// Append the movement into the movements array
			movements = append(movements, &movement)
		}

		// Check for any error that might have occured during the iteration
		// If there is none then return the array
		if err = rows.Err(); err != nil {
			return nil, err
		}

		return movements, nil

	}

// Method for inserting a new movement to the movement table
func (m MovementModel) InsertOneMovement(movement *Movement) error {
	// SQL query for inserting new record to the Movements table
	// And returning system generated data
	query := `
		INSERT INTO movements (name, description, image, tutorials, skilltype, muscles, difficulty, equipments, prerequisite)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, createdAt, version`

	// Args slice that holds the values for the placeholders in the SQL query
	// These values are from the movement struct
	args := 
		[]interface{}{movement.Name, movement.Description, movement.Image, pq.Array(movement.Tutorials), pq.Array(movement.Skilltype), pq.Array(movement.Muscles), movement.Difficulty, pq.Array(movement.Equipments), pq.Array(movement.Prerequisites)}
	
	// Execute and return the QueryRow() method wuth the query and the args slice as parameters
	// The Scan() method is used to return the system generated values
	return m.DB.QueryRow(query, args...).Scan(&movement.ID, &movement.CreatedAt, &movement.Version)
}
// Method for getting a new movement to the movement table
func (m MovementModel) GetOneMovement(id int64) (*Movement, error) {
	if id < 1 {
		return nil, ErrNotFound
	}

	// The query to fetch data of a specific movement
	query := `
		SELECT id, name, description, image, tutorials, skilltype, muscles, difficulty, equipments, prerequisite, version
		FROM movements
		WHERE id = $1`

	// Struct to hold the data returned from the query
	var movement Movement

	// Execute the query passing in the id parameter
	// Scan the response data into the fields of the movement struct
	err := m.DB.QueryRow(query, id).Scan(
		&movement.ID,
		&movement.Name,
		&movement.Description,
		&movement.Image,
		pq.Array(&movement.Tutorials),
		pq.Array(&movement.Skilltype),
		pq.Array(&movement.Muscles),
		&movement.Difficulty,
		pq.Array(&movement.Equipments),
		pq.Array(&movement.Prerequisites),
		&movement.Version,
	)

	// Check if the there is any error regarding the query
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}

	// Return a pointer to the movement struct
	return &movement, nil
}

// Method for updating a new movement to the movement table
func (m MovementModel) UpdateOneMovement(movement *Movement) error {
	// SQL query to update movements in the database
	query := `
		UPDATE movements
		SET name = $1, description = $2, image = $3, tutorials = $4, skilltype = $5, muscles = $6, difficulty = $7, equipments = $8, prerequisite = $9, version = version + 1
		WHERE id = $10 and version = $11
		RETURNING version`

	// Interface to hold all the placeholder values for the query
	args := []interface{}{
		movement.Name,
		movement.Description,
		movement.Image,
		pq.Array(movement.Tutorials),
		pq.Array(movement.Skilltype),
		pq.Array(movement.Muscles),
		movement.Difficulty,
		pq.Array(movement.Equipments),
		pq.Array(movement.Prerequisites),
		movement.ID,
		movement.Version,
	}

	// Execute the QueryRow method to update the record and scan the version value to the struct 
	err := m.DB.QueryRow(query, args...).Scan(&movement.Version)
	if err != nil {
		// If no rows were affected that means there was an edit conflict
		// Handling this error enables optimistic conurrency locking which avoids
		// Such edit conflicts in the case of a data race
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEditConflict
		} else {
			return err
		}
	}

	return nil
}

// Method for deleting a new movement to the movement table
func (m MovementModel) DeleteOneMovement(id int64) error {

	if id < 1 {
		return ErrNotFound
	}

	// SQL query for deleting a specific movement
	query := `
		DELETE FROM movements
		WHERE id = $1`
	
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	// The RowsAffected() method returns the number if rows affected from the query
	// If no rows were affected that means that no record was deleted
	// Which means the no record with the given id exists
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	} else if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}