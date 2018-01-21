package main

//todo: make this less gross
//todo: convert funds to decimal before logging

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

var db = loadDB()

func loadDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5432, "moonshot", "hodl", "moonshot-audit")

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {failGracefully(err, "Failed to open Postgres ")}

	err = db.Ping()
	if err != nil {failGracefully(err, "Failed to ping Postgres ")}

	return db
}

func main() {
	initRoutes()
}