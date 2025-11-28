package db


import (
	"fmt"
	"os"
	"database/sql"
	_ "github.com/lib/pq"

)

const port = 5432

func getEnv(env_key string, default_value string) string {
	value := os.Getenv(env_key)
	if value == "" {
		return default_value
	}
	return value
}


func getDSN() string {
	user := getEnv("POSTGRES_USER", "user")
	password := getEnv("POSTGRES_PASSWORD", "supersecret123")
	host := getEnv("POSTGRES_HOST", "localhost")
	db := getEnv("POSTGRES_DB", "mydb")

	DSN := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user, password, host, port, db) 

	return DSN
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


