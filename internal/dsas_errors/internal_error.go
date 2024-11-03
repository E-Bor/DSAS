package dsas_errors

import "fmt"

type InternalError struct {
	Err       error
	Message   string
	Operation string
}

func (e *InternalError) Unwrap() error {
	return e.Err
}

func (e *InternalError) Error() string {
	return e.Message
}

func NewInternalError(
	operation string,
	err error,
	fmtMessage string,
	args ...interface{},
) error {
	message := fmt.Sprintf(
		fmtMessage,
		args...,
	)
	message = fmt.Sprintf(
		"internal err: %s, in operation: %s original error: %s",
		message,
		operation,
		err.Error(),
	)
	return &InternalError{
		Message:   message,
		Err:       err,
		Operation: operation,
	}
}
