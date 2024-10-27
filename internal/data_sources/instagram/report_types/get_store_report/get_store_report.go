package main

import (
	"DSAS/internal/reports_registry"
	"fmt"
	"time"
)

func Report(traceId string) *reports_registry.ReportResultItem {
	fmt.Println("Instagram report start")
	time.Sleep(5 * time.Second)
	fmt.Println("Instagram report process")
	time.Sleep(5 * time.Second)
	fmt.Println("Instagram report end")
	return &reports_registry.ReportResultItem{
		TraceId: traceId,
		Err:     nil,
		Result:  nil,
	}
}
