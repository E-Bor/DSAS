package app

import (
	"DSAS/configs"
	grpc_server "DSAS/internal/app/grpc"
	"DSAS/internal/config"
	"DSAS/internal/core"
	"DSAS/internal/dsas_errors"
	"DSAS/internal/report_loader"
	"DSAS/internal/report_planner"
	"DSAS/internal/report_writer"
	"DSAS/internal/reports_registry"
	"DSAS/internal/storage"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	gRPCServer *grpc_server.GRPCServer
	log        *slog.Logger
	core       *core.DSASCore
}

func NewApp(
	log *slog.Logger,
	mainConfig *configs.MainConfig,
	dsasConfig *config.DsasConfig,
) (
	*App,
	error,
) {
	const op = "NewApp"
	sqlite, err := storage.NewStorage(
		storage.SQLite,
		mainConfig.SQLiteStoragePath,
		dsasConfig.DsasCoreConfig.DefaultAverageLoadTime,
	)
	if err != nil {
		return nil, err
	}

	reportPlanner := report_planner.NewReportPlanner(
		log,
		sqlite,
		dsasConfig.DsasCoreConfig.TraceIdLength,
		dsasConfig.DsasCoreConfig.QueueSleepTime,
		dsasConfig.DsasCoreConfig.QueueLength,
	)
	queueChannel := reportPlanner.StartPlannedQueue()

	reportLoadWorkerPool := report_loader.New(
		log,
		queueChannel,
		dsasConfig.DsasCoreConfig.LoadWorkersCount,
		sqlite,
		dsasConfig.DsasCoreConfig.WorkerPoolChannelBufferSize,
		dsasConfig.DsasCoreConfig.WorkerSleepTime,
	)

	reportLoadWorkerPool.Start(context.Background())

	reportWriter, err := report_writer.New(
		reportLoadWorkerPool.OutputChan,
		report_writer.SQLiteWriter,
		mainConfig.SQLiteStoragePath,
		"report_errors",
		log,
	)
	if err != nil {
		return nil, err
	}

	reportWriter.StoreAllData()

	reportMap, err := reports_registry.NewReportRegistry(dsasConfig.IntegrationsDir)
	if err != nil {
		return nil, dsas_errors.NewInternalError(
			op,
			err,
			"failed to create new report registry",
		)
	}

	dsasCore := core.New(
		reportMap,
		sqlite,
		reportPlanner,
		reportLoadWorkerPool,
		reportWriter,
	)
	server := grpc_server.NewGRPCServer(
		log,
		dsasConfig.CoreApiConfig.StartupHost,
		dsasCore,
	)

	return &App{
		gRPCServer: server,
		log:        log,
		core:       dsasCore,
	}, nil
}

func (a App) Start() error {
	go a.gRPCServer.Start()

	stop := make(
		chan os.Signal,
		1,
	)
	signal.Notify(
		stop,
		os.Interrupt,
		syscall.SIGTERM,
	)

	// block until sys call
	<-stop
	err := a.gRPCServer.Stop()
	if err != nil {
		return err
	}
	return nil
}
