package appError

import (
	"fmt"
	"net/http"
)

type badRequestError struct {
	errorMessage string
	statusCode   int
}

func (e *badRequestError) Error() string {
	return e.errorMessage
}

func (e *badRequestError) StatusCode() int {
	return e.statusCode
}

func NewBadRequestError(errorMessageFormat string, args ...interface{}) CommonError {
	return &badRequestError{
		errorMessage: fmt.Sprintf(errorMessageFormat, args...),
		statusCode:   http.StatusBadRequest,
	}
}
