package appError

import "fmt"

type CommonError interface {
	Error() string
	StatusCode() int
}

type commonError struct {
	errorMessage string
	statusCode   int
}

func (e *commonError) Error() string {
	return e.errorMessage
}

func (e *commonError) StatusCode() int {
	return e.statusCode
}

func NewCommonError(statusCode int, errorMessageFormat string, args ...interface{}) CommonError {
	return &commonError{
		errorMessage: fmt.Sprintf(errorMessageFormat, args...),
		statusCode:   statusCode,
	}
}
