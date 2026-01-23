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

// Sentinel error for errors.Is checking
var ErrUnknownMacro = &UnknownMacroError{}

// ToGRPCStatus converts internal errors to gRPC status codes
// This allows the same error types to be used for both HTTP and gRPC
// while converting them appropriately at the handler boundary
func ToGRPCStatus(err error) error {
	if err == nil {
		return nil
	}

	// Check for BadRequestError (validation errors, invalid input)
	var badReqErr *BadRequestError
	if errors.As(err, &badReqErr) {
		return status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// Check for UnknownMacroError (invalid macro in URL template)
	var unknownMacroErr *UnknownMacroError
	if errors.As(err, &unknownMacroErr) {
		return status.Errorf(codes.InvalidArgument, "%v", err)
	}

	// Default to Internal error for unexpected errors
	return status.Errorf(codes.Internal, "Internal error: %v", err)
}
