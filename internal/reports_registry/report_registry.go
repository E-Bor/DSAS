package reports_registry

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"plugin"
	"strings"
)

const logPath = "internal.core.report_registry"

type ReportResultItem struct {
	TraceId string
	Result  []map[string]interface{}
	Err     error
}

type ReportFunction func(traceId string) *ReportResultItem

type ReportRegistry struct {
	baseDir   string
	reportMap map[string]ReportFunction
}

func NewReportRegistry(baseDir string) (
	*ReportRegistry,
	error,
) {
	registry := &ReportRegistry{baseDir: baseDir}
	err := registry.getReportsFromAllIntegrations()
	return registry, err
}

func (r *ReportRegistry) getReportsFromAllIntegrations() error {
	result := make(map[string]ReportFunction)
	slog.Info(
		"created new report map",
		"path",
		logPath,
	)
	// walk in datasource dir
	err := filepath.Walk(
		r.baseDir,
		func(
			path string,
			info os.FileInfo,
			err error,
		) error {
			if err != nil {
				return err
			}
			// try to find compiled reports as plugins with .so ext
			if !info.IsDir() && strings.HasSuffix(
				info.Name(),
				".so",
			) {
				pluginPath := path
				p, err := plugin.Open(pluginPath)
				if err != nil {
					return fmt.Errorf(
						"failed to open plugin %s: %v",
						pluginPath,
						err,
					)
				}
				symReport, err := p.Lookup("Report")
				if err != nil {
					return fmt.Errorf(
						"failed to find Report in %s: %v",
						pluginPath,
						err,
					)
				}

				rt, ok := symReport.(func(traceId string) *ReportResultItem)
				if !ok {
					return fmt.Errorf(
						"symbol Report is not of type func() in %s",
						pluginPath,
					)
				}
				reportFunc := ReportFunction(rt)

				// key in format "datasource_name_some_report"
				parts := strings.Split(
					filepath.Dir(path),
					string(os.PathSeparator),
				)
				dirName := parts[len(parts)-4] // take datasource name by dir name
				fileName := strings.Split(
					info.Name(),
					".",
				)[0]
				key := fmt.Sprintf(
					"%s_%s",
					dirName,
					fileName,
				)

				result[key] = reportFunc
			}

			return nil
		},
	)

	if err != nil {
		return err
	}
	r.reportMap = result
	return nil
}

func (r *ReportRegistry) Get(dataSource, reportType string) (
	ReportFunction,
	error,
) {
	reportName := fmt.Sprintf(
		"%s_%s",
		dataSource,
		reportType,
	)

	reportFunc, ok := r.reportMap[reportName]
	if !ok {
		slog.Error(
			"failed to find report for ",
			"report name",
			reportName,
			"path",
			logPath,
		)
		return nil, fmt.Errorf(
			"failed to find report %s",
			reportName,
		)
	}
	return reportFunc, nil
}
