package data

import (
	"github.com/arnab4477/Parkour_API/internal/validator"
)

// This filters struct is to be used in the query parameters of urls
type Filters struct {
	Sort string
	Page int
	PageSize int
	SortSafeList []string
}

// Function that validates the filters of the query parameters
func ValidateFilters(v *validator.Validator, f Filters) {

	// Check that the sort parameters only contain the valid values
	// All the valid values are in the SortSafeList slice of the Filter struct
	v.Check(!validator.In(f.Sort, f.SortSafeList...), "sort", "invalid sort value")

	// Check that the page and page_size parameters contain sensible values
	v.Check(f.Page < 0, "page", "must be greater than zero") 
	v.Check(f.Page > 10_000_000, "page", "must be lower than 10 million") 
	v.Check(f.PageSize < 0, "page_size", "must be greater than zero") 
	v.Check(f.PageSize > 100, "page_size", "must be lower than one hundred") 
}

