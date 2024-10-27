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
	app.log.Info(
		"starting server in port",
		"port",
		l.Addr().String(),
	)

	if err := app.gRPCServer.Serve(l); err != nil {
		app.log.Error(err.Error())
		return
	}
}

func (app *GRPCServer) Stop() error {
	app.log.Info("stopping server")
	app.gRPCServer.GracefulStop()
	return nil
}
