package appError

import (
	"fmt"
	"net/http"
)

type notFoundError struct {
	errorMessage string
	statusCode   int
}

func (e *notFoundError) Error() string {
	return e.errorMessage
}

func (e *notFoundError) StatusCode() int {
	return e.statusCode
}

func NewNotFoundError(errorMessageFormat string, args ...interface{}) CommonError {
	return &notFoundError{
		errorMessage: fmt.Sprintf(errorMessageFormat, args...),
		statusCode:   http.StatusNotFound,
	}
}
