package errors

import "fmt"

type BadRequestError struct {
	Message string
}

func (e BadRequestError) Error() string {
	return e.Message
}

func BadRequestErrorf(format string, a ...interface{}) error {
	return &BadRequestError{
		Message: fmt.Sprintf(format, a...),
	}
}
