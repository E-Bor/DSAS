package sqlite_writer

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

type SqliteWriteStorage struct {
	db           *sql.DB
	errTableName string
	log          *slog.Logger
}

func New(
	storagePath string,
	errTableName string,
	log *slog.Logger,
) (
	*SqliteWriteStorage,
	error,
) {
	db, err := sql.Open(
		"sqlite3",
		storagePath,
	)
	if err != nil {
		return nil, err
	}

	return &SqliteWriteStorage{
		db:           db,
		errTableName: errTableName,
		log:          log,
	}, nil
}

func (s *SqliteWriteStorage) SaveReportSuccessResult(
	reportName string,
	TraceId string,
	resultRows []map[string]interface{},
) {
	ctx := context.Background()
	tableExist := s.tableExists(
		ctx,
		reportName,
	)

	if !tableExist {
		err := s.createTableByReportName(
			ctx,
			reportName,
		)
		if err != nil {
			slog.Error(
				"failed to create sqlite table",
				"error",
				err,
				"tried table",
				reportName,
			)
			reportName = fmt.Sprintf(
				"temp_table_%s_%s",
				reportName,
				TraceId,
			)
		}
		err = s.createTableByReportName(
			ctx,
			reportName,
		)
		if err != nil {
			slog.Error(
				"retry to create sqlite table failed",
				"error",
				err,
				"tried table",
				reportName,
			)
		}
	}

	tx, err := s.db.BeginTx(
		ctx,
		nil,
	)
	if err != nil {
		s.log.Error(
			"failed to begin tx",
			"error",
			err,
		)
	}

	inserts := s.getBatchInsertQueries(
		resultRows,
		reportName,
	)
	for _, insert := range inserts {
		_, err := tx.ExecContext(
			ctx,
			insert.query,
			insert.args...,
		)
		if err != nil {
			tx.Rollback()
			s.log.Error(
				"failed to execute transaction",
				"query",
				insert.query,
				"error",
				err,
			)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		s.log.Error(
			"failed to commit transaction",
			"error",
			err,
		)
	}

}

func (s *SqliteWriteStorage) SaveReportFailedResult(
	reportName string,
	traceId string,
	repError error,
) {
	ctx := context.Background()
	query := fmt.Sprintf(
		"INSERT INTO %s(report_type_name, trace_id, load_error) VALUES (?, ?, ?)",
		s.errTableName,
	)
	stmt, err := s.db.Prepare(
		fmt.Sprintf(query),
	)
	defer stmt.Close()
	if err != nil {
		return
	}

	_, err = stmt.ExecContext(
		ctx,
		reportName,
		traceId,
		repError.Error(),
	)
	if err != nil {
		s.log.Error(
			"error save failed result",
			"err",
			err,
		)
		return
	}
}

func (s *SqliteWriteStorage) tableExists(
	ctx context.Context,
	tableName string,
) bool {
	stmt, err := s.db.Prepare("SELECT name FROM sqlite_master WHERE type='table' AND name= ?")
	defer stmt.Close()
	if err != nil {
		s.log.Error(
			"error prepare request to fetch exist tables",
			err,
		)
		return false
	}
	row := stmt.QueryRowContext(
		ctx,
		tableName,
	)
	var table string
	err = row.Scan(&table)

	if err != nil {
		if errors.Is(
			err,
			sql.ErrNoRows,
		) {
			return false
		}
		s.log.Error(
			"error fetch exist tables",
			"err",
			err,
			"tableName",
			tableName,
		)
		return false
	}
	return true
}

func (s *SqliteWriteStorage) createTableByReportName(
	ctx context.Context,
	tableName string,
) error {
	query := fmt.Sprintf(
		`
    CREATE TABLE IF NOT EXISTS %s (
        report_date TEXT NOT NULL,
        info JSON NOT NULL
    );`,
		tableName,
	)

	_, err := s.db.ExecContext(
		ctx,
		query,
	)
	if err != nil {
		s.log.Error(
			"failed to execute create table statement",
			"error",
			err,
			"tableName",
			tableName,
		)
		return err
	}

	return nil
}

func (s *SqliteWriteStorage) getBatchInsertQueries(
	resultRows []map[string]interface{},
	tableName string,
) []struct {
	query string
	args  []interface{}
} {
	var inserts []struct {
		query string
		args  []interface{}
	}

	for _, resultRow := range resultRows {
		if len(resultRow) == 0 {
			continue
		}

		date, _ := resultRow["date"].(string)
		if date == "" {
			year, month, day := time.Now().UTC().Date()
			date = fmt.Sprintf(
				"%d-%02d-%02d",
				year,
				month,
				day,
			)
		}

		jsonRow, err := json.Marshal(resultRow)
		if err != nil {
			continue
		}

		insertRow := struct {
			query string
			args  []interface{}
		}{
			query: fmt.Sprintf(
				"INSERT INTO %s (report_date, info) VALUES (?, ?);",
				tableName,
			),
			args: []interface{}{
				date,
				string(jsonRow),
			},
		}
		inserts = append(
			inserts,
			insertRow,
		)
	}

	return inserts
}
