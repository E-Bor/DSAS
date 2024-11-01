package app

import (
	"DSAS/configs"
	"DSAS/internal/app/core"
	grpc_server "DSAS/internal/app/grpc"
	"DSAS/internal/config"
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

const workerPoolCount = 3

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
	server := grpc_server.NewGRPCServer(
		log,
		dsasConfig.CoreApiConfig.StartupHost,
	)

	sqlite, err := storage.NewStorage(
		storage.SQLite,
		mainConfig.SQLiteStoragePath,
	)
	if err != nil {
		return nil, err
	}

	reportPlanner := report_planner.NewReportPlanner(
		log,
		sqlite,
	)
	queueChannel := reportPlanner.StartPlannedQueue()

	reportLoadWorkerPool := report_loader.New(
		log,
		queueChannel,
		workerPoolCount,
		sqlite,
	)

	reportLoadWorkerPool.Start(context.Background())

	reportWriter := report_writer.New(
		reportLoadWorkerPool.OutputChan,
		report_writer.SQLiteWriter,
		mainConfig.SQLiteStoragePath,
		"report_errors",
		log,
	)

	reportWriter.StoreAllData()

	reportMap, err := reports_registry.NewReportRegistry(dsasConfig.IntegrationsDir)
	if err != nil {
		return nil, err
	}

	dsasCore := core.New(
		reportMap,
		sqlite,
		reportPlanner,
		reportLoadWorkerPool,
		reportWriter,
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
