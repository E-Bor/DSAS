package main

import (
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migrationsPath, migrationsTable string

	flag.StringVar(
		&storagePath,
		"storage-path",
		"",
		"Path to a directory containing the migration files",
	)
	flag.StringVar(
		&migrationsPath,
		"migrations-path",
		"",
		"Path to a directory containing the migration files",
	)
	flag.StringVar(
		&migrationsTable,
		"migrations-table",
		"migrations",
		"Path to a directory containing the migration files",
	)
	flag.Parse()

	if storagePath == "" {
		panic("storage-path is required")
	}
	if migrationsPath == "" {
		panic("migrations-path is required")
	}

	migrator, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf(
			"sqlite3://%s?x-migrations-table=%s",
			storagePath,
			migrationsTable,
		),
	)
	if err != nil {
		panic(err)
	}

	if err := migrator.Up(); err != nil {
		if err != migrate.ErrNoChange {
			fmt.Println("no migrations to apply")
			return
		}
		panic(err)
	}
	fmt.Println("migrations applied")
}