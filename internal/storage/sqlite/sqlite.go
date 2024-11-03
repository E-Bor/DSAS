package sqlite

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"math"
	"time"
)

type Storage struct {
	db                     *sql.DB
	defaultAverageLoadTime int64
}

func New(
	storagePath string,
	defaultAverageLoadTime int64,
) (
	*Storage,
	error,
) {
	db, err := sql.Open(
		"sqlite3",
		storagePath,
	)
	if err != nil {
		return nil, err
	}

	return &Storage{
		db:                     db,
		defaultAverageLoadTime: defaultAverageLoadTime,
	}, nil
}
func (s *Storage) GetAverageLoadingTime(
	ctx context.Context,
	reportName string,
) (
	time.Duration,
	error,
) {
	stmt, err := s.db.Prepare("SELECT AVG(load_time_sec) FROM report_stat WHERE report_type_name = ?")
	defer stmt.Close()
	if err != nil {
		return 0, err
	}

	row := stmt.QueryRowContext(
		ctx,
		reportName,
	)
	var avgInSeconds int64
	err = row.Scan(&avgInSeconds)

	if errors.Is(
		err,
		sql.ErrNoRows,
	) {
		avgInSeconds = s.defaultAverageLoadTime
	}
	avgDuration := time.Duration(avgInSeconds) * time.Second
	return avgDuration, nil
}

func (s *Storage) SaveAverageLoadingTime(
	reportName string,
	loadDuration time.Duration,
) {
	ctx := context.Background()
	stmt, err := s.db.Prepare("INSERT INTO report_stat(report_type_name, load_time_sec) VALUES (?, ?)")
	defer stmt.Close()
	if err != nil {
		return
	}
	loadingSeconds := int64(math.Round(loadDuration.Seconds()))
	_, err = stmt.ExecContext(
		ctx,
		reportName,
		loadingSeconds,
	)
	return
}
