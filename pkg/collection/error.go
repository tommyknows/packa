package collection

import (
	"fmt"
)

type Error struct {
	errors map[string]error
}

// Add an error to the collection
func (e *Error) Add(name string, err error) {
	if e.errors == nil {
		e.errors = make(map[string]error)
	}
	e.errors[name] = err
}

// Merge two error collections. err will overwrite
// possible entries in e
func (e *Error) Merge(err Error) {
	if e.errors == nil {
		e.errors = make(map[string]error)
	}
	for name, err := range err.errors {
		e.errors[name] = err
	}
}

// IfNotEmpty returns an error if there is at least one
// error collected, and nil if not
func (e *Error) IfNotEmpty() error {
	if e.errors == nil {
		return nil
	}
	return e
}

// Implements the error interface
func (e Error) Error() (s string) {
	// TODO: make this tabular?
	for name, err := range e.errors {
		s += fmt.Sprintf("\n%v:\t%s", name, err.Error())
	}
	return s
}
