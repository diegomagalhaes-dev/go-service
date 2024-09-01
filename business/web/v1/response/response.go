// Package response provides utilities for handling API responses,
// including error handling and pagination structures.
package response

import "errors"

type ErrorDocument struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

// PageDocument is the form used for API responses from query API calls.
type PageDocument[T any] struct {
	Items       []T `json:"items"`
	Total       int `json:"total"`
	Page        int `json:"page"`
	RowsPerPage int `json:"rowsPerPage"`
}

type Error struct {
	Err    error
	Status int
}

func NewError(err error, status int) error {
	return &Error{err, status}
}

func (re *Error) Error() string {
	return re.Err.Error()
}

func IsError(err error) bool {
	var re *Error
	return errors.As(err, &re)
}

func GetError(err error) *Error {
	var re *Error
	if !errors.As(err, &re) {
		return nil
	}
	return re
}

// NewPageDocument constructs a response value for a web paging response.
func NewPageDocument[T any](items []T, total int, page int, rowsPrePage int) PageDocument[T] {
	return PageDocument[T]{
		Items:       items,
		Total:       total,
		Page:        page,
		RowsPerPage: rowsPrePage,
	}
}
