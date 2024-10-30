package storage

import (
	"DSAS/internal/storage/sqlite"
	"context"
	"fmt"
	"time"
)

type StorageType string

const (
	SQLite StorageType = "sqlite"
)

type Storage interface {
	GetAverageLoadingTime(
		ctx context.Context,
		reportName string,
	) (
		time.Duration,
		error,
	)
	SaveAverageLoadingTime(
		reportName string,
		loadDuration time.Duration,
	)
}

func NewStorage(
	storageType StorageType,
	storagePath string,
) (
	Storage,
	error,
) {
	switch storageType {
	case SQLite:
		storage, err := sqlite.New(storagePath)
		return storage, err
	default:
		return nil, fmt.Errorf(
			"unknown storage type: %s",
			storageType,
		)
	}

}
