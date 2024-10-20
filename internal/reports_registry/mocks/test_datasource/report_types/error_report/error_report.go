package main

import (
	"errors"
)

func Report() error {
	err := errors.New("Test report with error")
	return err
}
