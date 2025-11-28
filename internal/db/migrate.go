package db

import (
	"embed"
	"database/sql"
	"github.com/pressly/goose/v3"
	_ "github.com/lib/pq"
)


//go:embed migrations/*.sql
var migrationsFs embed.FS

func InitMigrationDir() {
	goose.SetBaseFS(migrationsFs)
}

func MigrateUp(db *sql.DB) error {
	return goose.Up(db, "migrations")
}

func MigrateDown(db *sql.DB) error {
	return goose.DownTo(db, "migrations", 0)
}

func MigrateStatus(db *sql.DB) error {
	return goose.Status(db, "migrations")
}






