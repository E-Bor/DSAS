package app

import (
	"DSAS/configs"
	grpc_server "DSAS/internal/app/grpc"
	"DSAS/internal/config"
	"DSAS/internal/reports_registry"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type dsasCore struct {
	reportMap *reports_registry.ReportRegistry
}

type App struct {
	gRPCServer *grpc_server.GRPCServer
	log        slog.Logger
	core       dsasCore
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

	reportMap, err := reports_registry.NewReportRegistry(dsasConfig.IntegrationsDir)
	if err != nil {
		return nil, err
	}

	core := dsasCore{reportMap: reportMap}

	return &App{
		gRPCServer: server,
		log:        *log,
		core:       core,
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
