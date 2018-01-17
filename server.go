package main

import (
	"database/sql"
	"fmt"
)

func loadDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5432, "moonshot", "hodl", "moonshot-audit")
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		failGracefully(err, "Failed to open Postgres")
	}
	err = db.Ping()
	if err != nil {
		failGracefully(err, "Failed to Ping Postgres")
	} else {
		fmt.Println("Connectd to DB")
	}
	return db
}

func main() {

}