package appError

import (
	"fmt"
	"net/http"
)

type internalServerError struct {
	errorMessage string
	statusCode   int
}

func (e *internalServerError) Error() string {
	return e.errorMessage
}

func (e *internalServerError) StatusCode() int {
	return e.statusCode
}

func NewInternalServerError(errorMessageFormat string, args ...interface{}) CommonError {
	return &internalServerError{
		errorMessage: fmt.Sprintf(errorMessageFormat, args...),
		statusCode:   http.StatusInternalServerError,
	}
}
