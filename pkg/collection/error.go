package collection

import (
	"fmt"
)

// Error is a collection of errors in a map.
// To properly return an error-type, use the
// IfNotEmpty method
type Error map[string]error

// Add an error to the collection
func (e *Error) Add(name string, err error) {
	if *e == nil {
		*e = make(map[string]error)
	}
	(*e)[name] = err
}

// Merge two error collections. err will overwrite
// possible entries in e
func (e *Error) Merge(err Error) {
	if *e == nil {
		*e = make(map[string]error)
	}
	for name, err := range err {
		(*e)[name] = err
	}
}

// IfNotEmpty returns an error if there is at least one
// error collected, and nil if not
func (e *Error) IfNotEmpty() error {
	if *e == nil {
		return nil
	}
	return e
}

// Implements the error interface
func (e Error) Error() (s string) {
	// TODO: make this tabular?
	for name, err := range e {
		s += fmt.Sprintf("\n%v:\t%s", name, err.Error())
	}
	return s
}
