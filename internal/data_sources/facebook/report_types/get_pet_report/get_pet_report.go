package main

import (
	"errors"
	"log/slog"
)

func Report() error {
	err := errors.New("Facebook Error succsess with error")
	slog.Info("created error facebook get_pet_report")
	return err
}
