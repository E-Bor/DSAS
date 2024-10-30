package report_writer

import (
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
) *Writer {

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
			slog.Default().Error(
				"Failed to create sqlite writer",
				"error",
				err,
			)
		}
		writer.storage = sqlWriter
	}
	return &writer
}

func (w Writer) StoreAllData() {
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
}
