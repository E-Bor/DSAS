package sqlite_writer

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type SqliteWriteStorage struct {
	db           *sql.DB
	errTableName string
}

func New(
	storagePath string,
	errTableName string,
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
			reportName = fmt.Sprintf(
				"temp_table_%s_%s",
				reportName,
				TraceId,
			)
		}
	}

	insert := s.getQueryForBatchInsert(
		resultRows,
		reportName,
	)
	stmt, err := s.db.Prepare(insert)
	defer stmt.Close()
	if err != nil {
		return
	}
	_, err = stmt.ExecContext(
		ctx,
	)
}

func (s SqliteWriteStorage) SaveReportFailedResult(
	reportName string,
	TraceId string,
	err error,
) {
	ctx := context.Background()

	stmt, err := s.db.Prepare(
		fmt.Sprintf(
			"INSERT INTO %s(report_name, trace_id, load_error) VALUES (?, ?, ?)",
			s.errTableName,
		),
	)
	defer stmt.Close()
	if err != nil {
		return
	}

	_, err = stmt.ExecContext(
		ctx,
		reportName,
		TraceId,
		err,
	)

}

func (s *SqliteWriteStorage) tableExists(
	ctx context.Context,
	tableName string,
) bool {
	stmt, err := s.db.Prepare("SELECT name FROM sqlite_master WHERE type='table' AND name= ?")
	defer stmt.Close()
	if err != nil {
		return false
	}
	row := stmt.QueryRowContext(
		ctx,
		stmt,
		tableName,
	)
	var table string
	err = row.Scan(&table)
	if err != nil || table == "" {
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

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

// crunch, don`t secure
func (s SqliteWriteStorage) getQueryForBatchInsert(
	resultRows []map[string]interface{},
	tableName string,
) string {
	var inserts []string
	for _, resultRow := range resultRows {
		if len(resultRow) == 0 {
			continue
		}
		date, _ := resultRow["date"].(string)
		if date == "" {
			year, month, day := time.Now().UTC().Date()
			date = fmt.Sprintf(
				"%d-%d-%d",
				year,
				month,
				day,
			)
		}
		jsonRow, err := json.Marshal(resultRow)
		if err != nil {
			return ""
		}

		insterRow := fmt.Sprintf(
			"INSERT INTO %s (report_date, info) VALUES (%s, %s);",
			tableName,
			date,
			jsonRow,
		)
		inserts = append(
			inserts,
			insterRow,
		)
	}
	if len(inserts) > 0 {
		query := fmt.Sprintf(
			"BEGIN TRANSACTION;\n %s \n COMMIT;",
			strings.Join(
				inserts,
				"\n ",
			),
		)
		return query
	}
	return ""
}
