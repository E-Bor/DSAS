package app

import (
	"DSAS/configs"
	grpc_server "DSAS/internal/app/grpc"
	"DSAS/internal/config"
	"DSAS/internal/report_writer"
	"DSAS/internal/reports_registry"
	"DSAS/internal/storage"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type dsasCore struct {
	reportMap *reports_registry.ReportRegistry
	storage   storage.Storage
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

	sqlite, err := storage.NewStorage(
		storage.SQLite,
		mainConfig.SQLiteStoragePath,
	)
	if err != nil {
		return nil, err
	}

	testCH := make(
		chan *reports_registry.ReportResultItem,
		7,
	)

	testData := []reports_registry.ReportResultItem{
		{
			TraceId:    "trace-12345",
			ReportName: "sales_report",
			Result: []map[string]interface{}{
				{
					"product":  "Laptop",
					"quantity": 5,
					"price":    1000.0,
				},
				{
					"product":  "Smartphone",
					"quantity": 3,
					"price":    500.0,
				},
			},
			Err: nil,
		},
		{
			TraceId:    "trace-67890",
			ReportName: "inventory_report",
			Result: []map[string]interface{}{
				{
					"product": "Tablet",
					"stock":   20,
				},
				{
					"product": "Monitor",
					"stock":   15,
				},
			},
			Err: nil,
		},
		{
			TraceId:    "trace-111213",
			ReportName: "employee_report",
			Result: []map[string]interface{}{
				{
					"employee":     "John Doe",
					"hours_worked": 40,
					"department":   "Engineering",
				},
				{
					"employee":     "Jane Smith",
					"hours_worked": 35,
					"department":   "Marketing",
				},
			},
			Err: nil,
		},
		{
			TraceId:    "trace-141516",
			ReportName: "error_report",
			Result:     nil,
			Err:        errors.New("failed to generate report due to database error"),
		},
		{
			TraceId:    "trace-171819",
			ReportName: "financial_report",
			Result: []map[string]interface{}{
				{
					"revenue":  50000,
					"expenses": 30000,
					"profit":   20000,
				},
				{
					"revenue":  70000,
					"expenses": 40000,
					"profit":   30000,
				},
			},
			Err: nil,
		},
	}
	for _, val := range testData {
		testCH <- &val
	}

	writer := report_writer.New(
		testCH,
		report_writer.SQLiteWriter,
		mainConfig.SQLiteStoragePath,
		"report_errors",
		log,
	)
	go writer.StoreAllData()
	time.Sleep(10 * time.Second)
	reportMap, err := reports_registry.NewReportRegistry(dsasConfig.IntegrationsDir)
	if err != nil {
		return nil, err
	}
	core := dsasCore{
		reportMap: reportMap,
		storage:   sqlite,
	}

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
