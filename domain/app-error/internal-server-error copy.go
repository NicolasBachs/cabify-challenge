package appError

import (
	"fmt"
	"net/http"
)

type conflictResourceStateError struct {
	errorMessage string
	statusCode   int
}

func (e *conflictResourceStateError) Error() string {
	return e.errorMessage
}

func (e *conflictResourceStateError) StatusCode() int {
	return e.statusCode
}

func NewConflictResourceStateError(errorMessageFormat string, args ...interface{}) CommonError {
	return &internalServerError{
		errorMessage: fmt.Sprintf(errorMessageFormat, args...),
		statusCode:   http.StatusConflict,
	}
}
