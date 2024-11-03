package dsas_errors

import "fmt"

type ExternalError struct {
	Err       error
	Message   string
	Operation string
}

func (e *ExternalError) Unwrap() error {
	return e.Err
}

func (e *ExternalError) Error() string {
	return e.Message
}

func NewExternalError(
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
