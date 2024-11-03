package worker

import (
	"DSAS/internal/report_planner"
	"DSAS/internal/reports_registry"
	"context"
	"log/slog"
	"time"
)

type Worker struct {
	workerId         int
	log              *slog.Logger
	inputReportChan  <-chan *report_planner.ReportQueueItem
	outputReportChan chan<- *reports_registry.ReportResultItem
	workerSleepTime  time.Duration
	statStorage      AverageLoadingStorage
}

type AverageLoadingStorage interface {
	SaveAverageLoadingTime(
		reportName string,
		loadDuration time.Duration,
	)
}

func New(
	log *slog.Logger,
	workerId int,
	inputChan <-chan *report_planner.ReportQueueItem,
	outputChan chan<- *reports_registry.ReportResultItem,
	workerSleepTime time.Duration,
	statStorage AverageLoadingStorage,
) *Worker {
	log = log.With(
		slog.Int(
			"worker_id",
			workerId,
		),
	)
	return &Worker{
		workerId:         workerId,
		log:              log,
		inputReportChan:  inputChan,
		outputReportChan: outputChan,
		workerSleepTime:  workerSleepTime,
		statStorage:      statStorage,
	}
}

func (w *Worker) Start(ctx context.Context) {
	w.log.Info(
		"starting worker",
		"worker_id",
		w.workerId,
	)
	for {
		select {
		case job := <-w.inputReportChan:
			w.log.Info(
				"got job, start",
				slog.String(
					"TraceId",
					job.TraceId,
				),
			)
			timeStart := time.Now().UTC()

			result := job.ReportFunction(job.TraceId)
			result.ReportName = job.ReportName
			w.outputReportChan <- result

			loadTime := time.Since(timeStart)
			go w.statStorage.SaveAverageLoadingTime(
				job.ReportName,
				loadTime,
			)
			w.log.Info(
				"job end",
				slog.String(
					"loadTime",
					loadTime.String(),
				),
			)
		case <-ctx.Done():
			return
		default:
			time.Sleep(w.workerSleepTime)
		}
	}
}
