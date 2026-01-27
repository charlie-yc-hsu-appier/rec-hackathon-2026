package errors

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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

type InternalServerError struct {
	Message string
}

func (e InternalServerError) Error() string {
	return e.Message
}

func InternalServerErrorf(format string, a ...interface{}) error {
	return &InternalServerError{
		Message: fmt.Sprintf(format, a...),
	}
}

// Sentinel error for errors.Is checking
var ErrUnknownMacro = &UnknownMacroError{}

func ToGRPCStatus(err error) error {
	if err == nil {
		return nil
	}

	var badReqErr *BadRequestError
	if errors.As(err, &badReqErr) {
		return status.Errorf(codes.InvalidArgument, "%v", err)
	}

	var unknownMacroErr *UnknownMacroError
	if errors.As(err, &unknownMacroErr) {
		return status.Errorf(codes.Internal, "%v", err)
	}

	return status.Errorf(codes.Internal, "Internal error: %v", err)
}
