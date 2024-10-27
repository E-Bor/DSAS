package main

import (
	"DSAS/internal/reports_registry"
	"errors"
	"log/slog"
)

func Report(traceId string) *reports_registry.ReportResultItem {
	err := errors.New("Facebook Error succsess with error")
	slog.Info("created error facebook get_pet_report")
	return &reports_registry.ReportResultItem{
		TraceId: traceId,
		Err:     err,
		Result:  nil,
	}
}
