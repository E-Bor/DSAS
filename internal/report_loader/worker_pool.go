package report_loader

import (
	"DSAS/internal/report_loader/worker"
	"DSAS/internal/report_planner"
	"DSAS/internal/reports_registry"
	"context"
	"log/slog"
	"time"
)

type WorkerPool struct {
	log                  *slog.Logger
	workersExpectedCount int
	workersCurrenCount   int
	inputChan            <-chan *report_planner.ReportQueueItem
	OutputChan           chan *reports_registry.ReportResultItem
	workersCancel        context.CancelFunc
	statStorage          worker.AverageLoadingStorage
	workerSleepDuration  int
}

func New(
	logger *slog.Logger,
	inputChan <-chan *report_planner.ReportQueueItem,
	workerCount int,
	statStorage worker.AverageLoadingStorage,
	reportWorkerPoolChannelBuffer int,
	workerSleepDuration int,
) *WorkerPool {
	outputChan := make(
		chan *reports_registry.ReportResultItem,
		reportWorkerPoolChannelBuffer,
	)
	return &WorkerPool{
		log:                  logger,
		workersExpectedCount: workerCount,
		inputChan:            inputChan,
		OutputChan:           outputChan,
		statStorage:          statStorage,
		workerSleepDuration:  workerSleepDuration,
	}
}

func (w *WorkerPool) Start(ctx context.Context) {
	workerCtx, cancel := context.WithCancel(ctx)
	w.workersCancel = cancel
	for i := 0; i < w.workersExpectedCount; i++ {
		w.workersCurrenCount++
		wkr := worker.New(
			w.log,
			i,
			w.inputChan,
			w.OutputChan,
			time.Duration(w.workerSleepDuration)*time.Second,
			w.statStorage,
		)
		go wkr.Start(workerCtx)
	}
}

func (w *WorkerPool) Stop() {
	w.workersCancel()
	close(w.OutputChan)
}
