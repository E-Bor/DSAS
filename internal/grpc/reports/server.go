package reports

import (
	"context"
	reportsv1 "github.com/E-Bor/DSAS_Proto/gen/go/report_handlers"
)
import "google.golang.org/grpc"

type serverAPI struct {
	reportsv1.UnimplementedReportServer
}

func RegisterReport(gRPC *grpc.Server) {
	reportsv1.RegisterReportServer(
		gRPC,
		&serverAPI{},
	)
}

func (a *serverAPI) Start(
	ctx context.Context,
	req *reportsv1.StartRequest,
) (
	*reportsv1.StartResponse,
	error,
) {

	return nil, nil
}
