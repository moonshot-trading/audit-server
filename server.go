package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

var config = func() auditConfig { if runningInDocker() {return auditConfig{"audit-db"}} else {return auditConfig{"localhost"}}}()
var db = loadDB()

const SERVER = "1"
const FILENAME = "10userWorkLoad"

func loadDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.db, 5432, "moonshot", "hodl", "moonshot-audit")

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {failGracefully(err, "Failed to open Postgres ")}

	err = db.Ping()
	if err != nil {failGracefully(err, "Failed to ping Postgres ")}

	fmt.Println("Connected to DB at " + config.db)
	return db
}

func main() {
	initRoutes()
}
