package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/stdlib"
)

var (
	config = auditConfig{func() string {
		if runningInDocker() {
			return os.Getenv("AS_DB_HOST")
		} else {
			return "localhost"
		}
	}()}
	db       = loadDB()
	SERVER   = os.Getenv("AS_SERVER_ENUM")
	FILENAME = os.Getenv("FILENAME")
)

func loadDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.db, 5432, "moonshot", "hodl", "moonshot-audit")

	db, err := sql.Open("pgx", psqlInfo)
	if err != nil {
		failGracefully(err, "Failed to open Postgres ")
	}

	err = db.Ping()
	if err != nil {
		failGracefully(err, "Failed to ping Postgres ")
	}

	if err == nil {
		fmt.Println("Connected to DB at " + config.db)
	}
	return db
}

func main() {
	initQueues()
	initRoutes()
}
