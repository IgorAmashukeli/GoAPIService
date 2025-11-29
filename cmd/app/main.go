package main

import (
	"errors"
	"net/http"

	//"time"
	//"gorm.io/gorm"
	//"fmt"
	//"HighTalent/internal/api"
	"HighTalent/internal/api"
	"HighTalent/internal/db"
	"database/sql"
	"log"
)

// connecting to the sql database with goose
func Connect() *sql.DB {
	database, connect_err := db.GetDatabase()
	if connect_err != nil {
		log.Fatalf("%v", connect_err)
	}
	return database
}

// do migrations
func DoMigrations(database *sql.DB) {

	// up migrations
	up_err := db.MigrateUp(database)
	if up_err != nil {
		db.CloseDB(database)
		log.Fatalf("%v", up_err)
	}

	// checking the status of migtation
	status_err := db.MigrateStatus(database)
	if status_err != nil {
		db.CloseDB(database)
		log.Fatalf("%v", status_err)
	}
}

func main() {
	// connection
	database := Connect()
	log.Println("Succesfull connection")
	defer db.CloseDB(database)

	// migrations
	db.InitMigrationDir()
	DoMigrations(database)
	log.Println("Succesfull migration")

	// prepare db clients
	question_db, answer_db, ctx := db.PrepareDBClients(database)

	if ctx == nil {
		log.Fatalf("%v", errors.New("gorm prepare error"))
	}

	//create api
	mux := api.CreateApi(question_db, answer_db)

	// listen
	log.Println("Server starting on port 8080...")
	http.ListenAndServe(":8080", mux)
}
