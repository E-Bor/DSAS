package main

import "DSAS/internal/reports_registry"

func Report(traceId string) *reports_registry.ReportResultItem {
	return &reports_registry.ReportResultItem{TraceId: traceId}
}
