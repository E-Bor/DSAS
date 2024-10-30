package report_writer

import "DSAS/internal/reports_registry"

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
}

func New(
	storage storageForWriteReport,
	dataCh <-chan *reports_registry.ReportResultItem,
) *Writer {
	return &Writer{
		storage: storage,
		dataCh:  dataCh,
	}
}

func (w Writer) storeAllData() {
	select {
	case reportLoadResult := <-w.dataCh:
		if reportLoadResult.Err != nil {
			w.storage.SaveReportFailedResult(
				reportLoadResult.ReportName,
				reportLoadResult.TraceId,
				reportLoadResult.Err,
			)
		}
		w.storage.SaveReportSuccessResult(
			reportLoadResult.ReportName,
			reportLoadResult.TraceId,
			reportLoadResult.Result,
		)
	}
}
