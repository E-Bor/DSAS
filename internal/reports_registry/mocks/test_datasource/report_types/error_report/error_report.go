package main

import (
	"DSAS/internal/reports_registry"
	"errors"
)

func Report(traceId string) *reports_registry.ReportResultItem {
	err := errors.New("Test report with error")
	return &reports_registry.ReportResultItem{
		TraceId: traceId,
		Result:  nil,
		Err:     err,
	}
}
