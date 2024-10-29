package worker

import (
	"DSAS/internal/report_loader/worker/mocks"
	"DSAS/internal/report_planner"
	"DSAS/internal/reports_registry"
	"context"
	"errors"
	"log/slog"
	"reflect"
	"testing"
)

//go:generate go run github.com/vektra/mockery/v2@v2.46.3 --name=AverageLoadingStorage
func TestWorker_Start(t *testing.T) {

	inputChan := make(
		chan *report_planner.ReportQueueItem,
		3,
	)
	outputChan := make(
		chan *reports_registry.ReportResultItem,
		3,
	)
	mockAnalyticDB := mocks.NewAverageLoadingStorage(t)
	wkr := &Worker{
		workerId:         1,
		log:              slog.Default(),
		inputReportChan:  inputChan,
		outputReportChan: outputChan,
		workerSleepTime:  0,
		statStorage:      mockAnalyticDB,
	}

	testCases := []struct {
		reportItem *report_planner.ReportQueueItem
		resultItem *reports_registry.ReportResultItem
	}{
		{
			reportItem: &report_planner.ReportQueueItem{
				TraceId: "trace1",
				ReportFunction: func(traceId string) *reports_registry.ReportResultItem {
					return &reports_registry.ReportResultItem{
						TraceId: traceId,
						Err:     nil,
						Result: []map[string]interface{}{
							{
								"test1": "test1",
							},
						},
					}
				},
			},
			resultItem: &reports_registry.ReportResultItem{
				TraceId: "trace1",
				Err:     nil,
				Result: []map[string]interface{}{
					{
						"test1": "test1",
					},
				},
			},
		},
		{
			reportItem: &report_planner.ReportQueueItem{
				TraceId: "trace2",
				ReportFunction: func(traceId string) *reports_registry.ReportResultItem {
					return &reports_registry.ReportResultItem{
						TraceId: traceId,
						Err:     errors.New("test error"),
						Result:  []map[string]interface{}{},
					}
				},
			},
			resultItem: &reports_registry.ReportResultItem{
				TraceId: "trace2",
				Err:     errors.New("test error"),
				Result:  []map[string]interface{}{},
			},
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	go wkr.Start(ctx)

	for _, testCase := range testCases {
		inputChan <- testCase.reportItem
		if currentVal := <-outputChan; currentVal.TraceId != testCase.resultItem.TraceId || !reflect.DeepEqual(
			currentVal.Err,
			testCase.resultItem.Err,
		) || !reflect.DeepEqual(
			currentVal.Result,
			testCase.resultItem.Result,
		) {
			t.Errorf(
				"Wrong result item received. Expected %v, got %v",
				testCase.resultItem,
				currentVal,
			)
		}
	}
	cancel()
}
