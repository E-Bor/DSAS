package grpc_server

import (
	"DSAS/internal/grpc/reports"
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type GRPCServer struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

// TODO: Implement auth

func NewGRPCServer(
	log *slog.Logger,
	port int,
) *GRPCServer {
	server := grpc.NewServer()

	reports.RegisterReport(server)

	return &GRPCServer{
		log:        log,
		gRPCServer: server,
		port:       port,
	}
}

func (app *GRPCServer) Start() {
	const path = "app.start"
	log := app.log.With(
		slog.String(
			"op",
			path,
		),
	)
	l, err := net.Listen(
		"tcp",
		fmt.Sprintf(
			":%d",
			app.port,
		),
	)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	log.Info(
		"starting server in port",
		"port",
		l.Addr().String(),
	)

	if err := app.gRPCServer.Serve(l); err != nil {
		log.Error(err.Error())
		return
	}
}

func (app *GRPCServer) Stop() error {
	const path = "app.stop"
	log := app.log.With(
		slog.String(
			"op",
			path,
		),
	)
	log.Info("stopping server")
	app.gRPCServer.GracefulStop()
	return nil
}
