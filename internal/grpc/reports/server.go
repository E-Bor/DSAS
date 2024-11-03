package reports

import (
	"context"
	reportsv1 "github.com/E-Bor/DSAS_Proto/gen/go/report_handlers"
	"time"
)
import "google.golang.org/grpc"

type DsasCore interface {
	AddReportToQueue(
		datasource, reportType string,
		estimatedDate,
		dateFrom,
		dateTo time.Time,
	) (
		string,
		error,
	)
}

type serverAPI struct {
	reportsv1.ReportServer
	dsasCore DsasCore
}

func RegisterReport(
	gRPC *grpc.Server,
	dsasCore DsasCore,
) {
	reportsv1.RegisterReportServer(
		gRPC,
		&serverAPI{dsasCore: dsasCore},
	)
}

func (a *serverAPI) Start(
	ctx context.Context,
	req *reportsv1.StartRequest,
) (
	*reportsv1.StartResponse,
	error,
) {
	traceId, err := a.dsasCore.AddReportToQueue(
		req.GetDatasourceName(),
		req.GetReportTypeName(),
		req.GetEta().AsTime(),
		req.GetDateFrom().AsTime(),
		req.GetDateTo().AsTime(),
	)

	if err != nil {
		return &reportsv1.StartResponse{
			Status:  reportsv1.Status_failed,
			TraceId: traceId,
		}, err
	}

	return &reportsv1.StartResponse{
		Status:  reportsv1.Status_success,
		TraceId: traceId,
	}, nil
}
