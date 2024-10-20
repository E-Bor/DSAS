package reports_registry

import (
	"errors"
	"reflect"
	"testing"
)

//go:generate go build -buildmode=plugin -o ../reports_registry/mocks/test_datasource/report_types/test_error_report.so ../reports_registry/mocks/test_datasource/report_types/error_report/error_report.go
//go:generate go build -buildmode=plugin -o ../reports_registry/mocks/test_datasource/report_types/test_success_report.so ../reports_registry/mocks/test_datasource/report_types/success_report/success_report.go
func TestNewReportRegistry(t *testing.T) {
	type args struct {
		baseDir string
	}
	tests := []struct {
		name        string
		args        args
		wantMapType *ReportRegistry
		wantMapKeys []string
		wantErr     bool
	}{
		{
			name:        "Get report registry with mock reports",
			args:        args{baseDir: "mocks"},
			wantMapType: &ReportRegistry{},
			wantMapKeys: []string{
				"mocks_test_error_report",
				"mocks_test_success_report",
			},
			wantErr: false,
		},
		{
			name:        "Get report registry with doesnt existed dir",
			args:        args{baseDir: "test_not_exist_dir"},
			wantMapType: &ReportRegistry{},
			wantMapKeys: []string{},
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				got, err := NewReportRegistry(tt.args.baseDir)
				if (err != nil) != tt.wantErr {
					t.Errorf(
						"NewReportRegistry() error = %v, wantErr %v",
						err,
						tt.wantErr,
					)
					return
				}
				if reflect.TypeOf(got) != reflect.TypeOf(tt.wantMapType) {
					t.Errorf(
						"NewReportRegistry() got = %v, want %v",
						got,
						tt.wantMapType,
					)
				}
				for _, key := range tt.wantMapKeys {
					if gotFunc, ok := got.reportMap[key]; ok && reflect.TypeOf(gotFunc) == reflect.TypeOf(func() error { return nil }) {
						continue
					} else {
						t.Errorf(
							"NewReportRegistry() got = %v, want func() error",
							reflect.TypeOf(gotFunc),
						)
					}
				}
			},
		)
	}
}

func TestReportRegistry_Get(t *testing.T) {
	type fields struct {
		baseDir   string
		reportMap map[string]func() error
	}
	type args struct {
		dataSource string
		reportType string
	}
	ResultError := errors.New("test error")
	tests := []struct {
		name     string
		fields   fields
		args     args
		want     error
		wantType reflect.Type
		wantErr  bool
	}{
		{
			name: "Get report from registry",
			fields: fields{
				baseDir: "mocks",
				reportMap: map[string]func() error{
					"test_datasource1_test_report1": func() error { return nil },
				},
			},
			args: args{
				dataSource: "test_datasource1",
				reportType: "test_report1",
			},
			wantType: reflect.TypeOf(func() error { return nil }),
			want:     nil,
			wantErr:  false,
		},
		{
			name: "Get report from registry",
			fields: fields{
				baseDir: "mocks",
				reportMap: map[string]func() error{
					"test_datasource2_test_report2": func() error { return ResultError },
				},
			},
			args: args{
				dataSource: "test_datasource2",
				reportType: "test_report2",
			},
			wantType: reflect.TypeOf(func() error { return ResultError }),
			want:     ResultError,
			wantErr:  false,
		},
		{
			name: "Get report from registry",
			fields: fields{
				baseDir:   "mocks",
				reportMap: map[string]func() error{},
			},
			args: args{
				dataSource: "test_datasource2",
				reportType: "test_report2",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				r := &ReportRegistry{
					baseDir:   tt.fields.baseDir,
					reportMap: tt.fields.reportMap,
				}
				got, err := r.Get(
					tt.args.dataSource,
					tt.args.reportType,
				)
				if (err != nil) != tt.wantErr {
					t.Errorf(
						"Get() error = %v, wantErr %v",
						err,
						tt.wantErr,
					)
					return
				}

				if (err != nil) == tt.wantErr {
					return
				}

				callResult := got()

				if !(reflect.TypeOf(got) == tt.wantType && errors.Is(
					callResult,
					tt.want,
				)) {
					t.Errorf(
						"Get() got = %v, want %v",
						callResult,
						tt.want,
					)
				}
			},
		)
	}
}
