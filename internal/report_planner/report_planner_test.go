package report_planner

import (
	"DSAS/internal/report_planner/mocks"
	"DSAS/internal/reports_registry"
	"container/list"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"reflect"
	"testing"
	"time"
)

//go:generate go run github.com/vektra/mockery/v2@v2.46.3 --name=AverageLoadingStorage

func TestNewReportPlanner(t *testing.T) {
	type args struct {
		logger                *slog.Logger
		averageLoadingStorage AverageLoadingStorage
	}
	logger := &slog.Logger{}
	storage := mocks.NewAverageLoadingStorage(t)

	tests := []struct {
		name string
		args args
		want *ReportPlanner
	}{
		{
			name: "init test",
			args: args{
				logger:                logger,
				averageLoadingStorage: storage,
			},
			want: &ReportPlanner{
				log:                   logger,
				averageLoadingStorage: storage,
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				got := NewReportPlanner(
					tt.args.logger,
					tt.args.averageLoadingStorage,
				)

				if !reflect.DeepEqual(
					got.averageLoadingStorage,
					tt.want.averageLoadingStorage,
				) || !reflect.DeepEqual(
					got.log,
					tt.want.log,
				) {
					t.Errorf(
						"NewReportPlanner() = %v, want %v",
						got,
						tt.want,
					)
				}
				if reflect.TypeOf(got.reportsLocalQueue) != reflect.TypeOf(&list.List{}) {
					t.Errorf(
						"NewReportPlanner() = %v, want %v",
						got.reportsLocalQueue,
						reflect.TypeOf(&list.List{}),
					)
				}
			},
		)
	}
}

func TestReportPlanner_Add(t *testing.T) {
	type reports struct {
		reportName                      string
		reportFunc                      reports_registry.ReportFunction
		estimatedDate                   time.Time
		dateFrom                        time.Time
		dateTo                          time.Time
		averageLoadingDurationForOneDay time.Duration
	}
	tests := []struct {
		name                string
		loadingStorageCall  *mock.Call
		reports             []reports
		expectedAddSequence [][]string
	}{
		{
			name: "Test should add items in queue in right order",
			reports: []reports{
				{
					reportName:                      "Report_1",
					reportFunc:                      func(traceId string) *reports_registry.ReportResultItem { return nil },
					estimatedDate:                   time.Now().Add(20 * time.Hour),
					dateFrom:                        time.Now().Add(-3 * 24 * time.Hour),
					dateTo:                          time.Now().Add(-3 * 24 * time.Hour),
					averageLoadingDurationForOneDay: time.Minute,
				},
				{
					reportName:                      "Report_2",
					reportFunc:                      func(traceId string) *reports_registry.ReportResultItem { return nil },
					estimatedDate:                   time.Now().Add(20 * time.Hour),
					dateFrom:                        time.Now().Add(-3 * 24 * time.Hour),
					dateTo:                          time.Now().Add(-3 * 24 * time.Hour),
					averageLoadingDurationForOneDay: time.Minute,
				},
				{
					reportName:                      "Report_3",
					reportFunc:                      func(traceId string) *reports_registry.ReportResultItem { return nil },
					estimatedDate:                   time.Now().Add(1 * time.Hour),
					dateFrom:                        time.Now().Add(-4 * 24 * time.Hour),
					dateTo:                          time.Now().Add(-3 * 24 * time.Hour),
					averageLoadingDurationForOneDay: 15 * time.Minute,
				},
				{
					reportName:                      "Report_4",
					reportFunc:                      func(traceId string) *reports_registry.ReportResultItem { return nil },
					estimatedDate:                   time.Now().Add(2 * time.Minute),
					dateFrom:                        time.Now().Add(-4 * 24 * time.Hour),
					dateTo:                          time.Now().Add(-3 * 24 * time.Hour),
					averageLoadingDurationForOneDay: 2 * time.Minute,
				},
				{
					reportName:                      "Report_5",
					reportFunc:                      func(traceId string) *reports_registry.ReportResultItem { return nil },
					estimatedDate:                   time.Now().Add(10 * time.Minute),
					dateFrom:                        time.Now().Add(-4 * 24 * time.Hour),
					dateTo:                          time.Now().Add(-3 * 24 * time.Hour),
					averageLoadingDurationForOneDay: 5 * time.Minute,
				},
			},
			expectedAddSequence: [][]string{
				{"Report_1"},
				{
					"Report_1",
					"Report_2",
				},
				{
					"Report_3",
					"Report_1",
					"Report_2",
				},
				{
					"Report_4",
					"Report_3",
					"Report_1",
					"Report_2",
				},
				{
					"Report_4",
					"Report_5",
					"Report_3",
					"Report_1",
					"Report_2",
				},
			},
		},
	}

	localStorage := mocks.NewAverageLoadingStorage(t)
	reportPlanner := NewReportPlanner(
		slog.Default(),
		localStorage,
	)
	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				for i, report := range tt.reports {
					localStorage.On(
						"GetAverageLoadingTime",
						mock.AnythingOfType("context.backgroundCtx"),
						report.reportName,
					).Return(
						report.averageLoadingDurationForOneDay,
						nil,
					)
					reportPlanner.Add(
						report.reportName,
						report.reportFunc,
						report.estimatedDate,
						report.dateFrom,
						report.dateTo,
					)
					currentSequence := reportPlanner.GetAllSequence()
					var currentReportNames []string
					for _, report := range currentSequence {
						currentReportNames = append(
							currentReportNames,
							report.ReportName,
						)
					}

					if !reflect.DeepEqual(
						currentReportNames,
						tt.expectedAddSequence[i],
					) {
						t.Errorf(
							"Unexpected sequence: %v, want %v",
							currentReportNames,
							tt.expectedAddSequence[i],
						)
					}
				}
			},
		)
	}
}

func TestReportPlanner_Get(t *testing.T) {

	tests := []struct {
		name string
		want chan *ReportQueueItem
	}{
		{
			name: "get should return channel",
			want: make(chan *ReportQueueItem),
		},
	}
	localStorage := mocks.NewAverageLoadingStorage(t)
	reportPlanner := NewReportPlanner(
		slog.Default(),
		localStorage,
	)

	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				currentChannel := reportPlanner.Get()
				if reflect.TypeOf(currentChannel) != reflect.TypeOf(tt.want) {
					t.Errorf(
						"Get() = %v, want %v",
						currentChannel,
						tt.want,
					)
				}
			},
		)
	}
}
