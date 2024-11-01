package core

import (
	"DSAS/internal/report_planner"
	"DSAS/internal/reports_registry"
	"DSAS/internal/storage"
	"context"
	"fmt"
	"time"
)

type ReportMap interface {
	Get(dataSource, reportType string) (
		reports_registry.ReportFunction,
		error,
	)
}

type ReportPlanner interface {
	Add(
		reportName string,
		reportFunc reports_registry.ReportFunction,
		estimatedDate time.Time,
		dateFrom time.Time,
		dateTo time.Time,
	) string

	StartPlannedQueue() <-chan *report_planner.ReportQueueItem
}

type ReportLoader interface {
	Start(ctx context.Context)
	Stop()
}

type ReportWriter interface {
	StoreAllData()
}

type DSASCore struct {
	reportMap     ReportMap
	storage       storage.Storage
	reportPlanner ReportPlanner
	reportLoader  ReportLoader
	reportWriter  ReportWriter
}

func New(
	reportMap ReportMap,
	storage storage.Storage,
	reportPlanner ReportPlanner,
	reportLoader ReportLoader,
	reportWriter ReportWriter,
) *DSASCore {
	return &DSASCore{
		reportMap:     reportMap,
		storage:       storage,
		reportPlanner: reportPlanner,
		reportLoader:  reportLoader,
		reportWriter:  reportWriter,
	}
}

func (c DSASCore) AddReportToQueue(
	datasource, reportType string,
	estimatedDate,
	dateFrom,
	dateTo time.Time,
) (
	string,
	error,
) {
	reportFunc, err := c.reportMap.Get(
		datasource,
		reportType,
	)

	if err != nil {
		return "", err
	}
	reportName := fmt.Sprintf(
		"%s-%s",
		datasource,
		reportType,
	)

	traceId := c.reportPlanner.Add(
		reportName,
		reportFunc,
		estimatedDate,
		dateFrom,
		dateTo,
	)

	if traceId != "" {
		return traceId, nil
	}
	return "", nil
}
