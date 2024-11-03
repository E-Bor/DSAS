package main

import (
	"DSAS/internal/reports_registry"
	"fmt"
	"time"
)

func Report(traceId string) *reports_registry.ReportResultItem {
	fmt.Println("Facebook report start")
	time.Sleep(10 * time.Second)

	result := []map[string]interface{}{
		{
			"date":    "2024-05-06",
			"revenue": 23,
			"apps": []string{
				"app1",
				"app2",
			},
			"total_count": 4,
		},
		{
			"date":        "2024-05-02",
			"revenue":     11,
			"apps":        []string{"app1"},
			"total_count": 5,
		},
		{
			"date":    "2024-05-10",
			"revenue": 22,
			"apps": []string{
				"app1, app2",
				"app3",
			},
			"total_count": 1,
		},
	}
	fmt.Println("Facebook report end")
	return &reports_registry.ReportResultItem{
		TraceId: traceId,
		Err:     nil,
		Result:  result,
	}
}
