package db

import (
	"database/sql"
	"embed"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationsFs embed.FS

func InitMigrationDir() {
	goose.SetBaseFS(&migrationsFs)
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
