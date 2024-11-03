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
	result := []map[string]interface{}{
		{
			"ads": 23,
			"apps": []string{
				"app1",
				"app2",
			},
			"total_count": 4,
		},
		{
			"ads":         11,
			"apps":        []string{"app1"},
			"total_count": 5,
		},
		{
			"ads": 22,
			"apps": []string{
				"app1, app2",
				"app3",
			},
			"total_count": 1,
		},
	}
	fmt.Println("Instagram report end")
	return &reports_registry.ReportResultItem{
		TraceId: traceId,
		Err:     nil,
		Result:  result,
	}
}
