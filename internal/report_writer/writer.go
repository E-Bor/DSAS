package report_writer

import (
	"DSAS/internal/dsas_errors"
	"DSAS/internal/report_writer/sqlite_writer"
	"DSAS/internal/reports_registry"
	"log/slog"
)

type WriterType string

const (
	SQLiteWriter WriterType = "sqlite"
)

type storageForWriteReport interface {
	SaveReportFailedResult(
		string,
		string,
		error,
	)
	SaveReportSuccessResult(
		string,
		string,
		[]map[string]interface{},
	)
}

type Writer struct {
	storage storageForWriteReport
	dataCh  <-chan *reports_registry.ReportResultItem
	log     *slog.Logger
}

func New(
	dataCh <-chan *reports_registry.ReportResultItem,
	writerType WriterType,
	storagePath string,
	errTableName string,
	logger *slog.Logger,
) (
	*Writer,
	error,
) {
	const op = "writer.New"

	var writer Writer

	writer.dataCh = dataCh
	writer.log = logger

	switch writerType {
	case SQLiteWriter:
		sqlWriter, err := sqlite_writer.New(
			storagePath,
			errTableName,
			logger,
		)
		if err != nil {
			return nil, dsas_errors.NewInternalError(
				op,
				err,
				"Failed to create sqlite writer",
			)
		}
		writer.storage = sqlWriter
	}
	return &writer, nil
}

func (w Writer) StoreAllData() {
	go func() {
		for {
			select {
			case reportLoadResult := <-w.dataCh:
				if reportLoadResult.Err != nil {
					w.storage.SaveReportFailedResult(
						reportLoadResult.ReportName,
						reportLoadResult.TraceId,
						reportLoadResult.Err,
					)
				} else {
					w.storage.SaveReportSuccessResult(
						reportLoadResult.ReportName,
						reportLoadResult.TraceId,
						reportLoadResult.Result,
					)
				}
			}
		}
	}()
}
