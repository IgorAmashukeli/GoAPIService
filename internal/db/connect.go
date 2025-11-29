package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func getEnv(env_key string, default_value string) string {
	value := os.Getenv(env_key)
	if value == "" {
		return default_value
	}
	return value
}

func getDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
}

func OpenDB() (*sql.DB, error) {

	DSN := getDSN()

	return sql.Open("postgres", DSN)

}

func CloseDB(database *sql.DB) {
	database.Close()
}

func PingDB(database *sql.DB) error {
	return database.Ping()
}

func GetDatabase() (*sql.DB, error) {
	database, open_err := OpenDB()

	if open_err != nil {
		return database, open_err
	}
	ping_err := PingDB(database)
	if ping_err != nil {
		CloseDB(database)
		return nil, ping_err

	}

	return database, nil
}
