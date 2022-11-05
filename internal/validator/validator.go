package validator

import (
	"regexp"
)

// A map to contain all the validation errors
type Validator struct {
	Errors map[string]string
}

// This function creates a new Validator instance with an empty map
func NewValidator() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// This function returns true if the Errors map is empty
func (v *Validator) NoErrors() bool {
	return len(v.Errors) == 0
}

// This function adds an unique error message to the map
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// This adds an error message to the map only if an error message is not 'ok'
func (v *Validator) Check(ok bool, key, message string) {
	if ok {
		v.AddError(key, message)
	}
}

// This returns true if a particular v alue is in a list
func In(value string, list ...string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}
	return false
}

// This returns true if a string matches a RegEx
func Matches(value string, regex *regexp.Regexp) bool {
	return regex.MatchString(value)
}

//this returns true if all the strings in a slice are unique
func IsUnique(values []string) bool {
	// Create a map that will only contain the unique values of the slice
	uniqueValue := make(map[string]bool)
	for _, value := range values {
		// If a value is repeated then it won't be added to the map
		// Thus making its length shorter than the original slice
		uniqueValue[value] = true
	}
	return len(values) == len(uniqueValue)
}

// This is the Regular Expression for validaing email addresses
// This was taken from https://html.spec.whatwg.org/multipage/input.html#valid-e-mail-address
var (
	EmailRegEx = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)


