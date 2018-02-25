package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
)

var (
	config = auditConfig{func() string {
		if runningInDocker() {
			return "audit-db"
		} else {
			return "localhost"
		}
	}()}
	db            = loadDB()
	semaphoreChan = make(chan struct{}, 40)
)

const SERVER = "1"
const FILENAME = "10userWorkLoad"

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

	fmt.Println("Connected to DB at " + config.db)
	return db
}

func main() {
	initQueues()
	initRoutes()
}
