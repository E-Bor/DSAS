package core

import (
	"DSAS/internal/dsas_errors"
	"DSAS/internal/report_planner"
	"DSAS/internal/reports_registry"
	"DSAS/internal/storage"
	"context"
	"errors"
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

func (c *DSASCore) AddReportToQueue(
	datasource, reportType string,
	estimatedDate,
	dateFrom,
	dateTo time.Time,
) (
	string,
	error,
) {
	const op = "core.AddReportToQueue"
	if datasource == "" || reportType == "" {
		return "", dsas_errors.NewExternalError(
			op,
			errors.New("empty datasource or report type"),
			"got datasource %s, got report type %s",
			datasource,
			reportType,
		)
	}
	if dateFrom.IsZero() || dateTo.IsZero() || dateFrom.After(dateTo) {
		return "", dsas_errors.NewExternalError(
			op,
			errors.New("error date range provided"),
			"got date from %s, got date to %s",
			dateFrom.String(),
			dateTo.String(),
		)
	}
	if estimatedDate.IsZero() {
		return "", dsas_errors.NewExternalError(
			op,
			errors.New("empty estimated date provided"),
			"got estimated date %s",
			estimatedDate.String(),
		)
	}

	reportFunc, err := c.reportMap.Get(
		datasource,
		reportType,
	)

	if err != nil {
		return "", err
	}
	reportName := fmt.Sprintf(
		"%s_%s",
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
