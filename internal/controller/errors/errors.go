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

type UnknownMacroError struct {
	Macro string
}

func (e UnknownMacroError) Error() string {
	return fmt.Sprintf("unknown macro: %s", e.Macro)
}

// Is implements error matching for errors.Is
func (e UnknownMacroError) Is(target error) bool {
	_, ok := target.(*UnknownMacroError)
	return ok
}

func NewUnknownMacroError(macro string) error {
	return &UnknownMacroError{
		Macro: macro,
	}
}

// Sentinel error for errors.Is checking
var ErrUnknownMacro = &UnknownMacroError{}
