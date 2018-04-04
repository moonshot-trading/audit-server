package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var AlreadyDumped = 0

func dumpLogCommand() {
	//logDumpCommand(w, r)
	var err error
	var f *os.File
	if AlreadyDumped < 14 {
		AlreadyDumped++
		return
	}
	f, err = os.Create("log.xml")

	if err != nil {
		failGracefully(err, "Failed to open log file ")
	}

	rows, err := db.Query("SELECT logType, (extract(EPOCH FROM timestamp) * 1000)::BIGINT as timestamp, server, transactionNum, command, username, stockSymbol, filename, (funds::DECIMAL)/100 as funds, cryptokey, (price::DECIMAL)/100 as price, quoteServerTime, action, errorMessage, debugMessage FROM audit_log;")
	if err != nil {
		failGracefully(err, "Failed to query audit DB ")
	}
	defer rows.Close()
	f.Write([]byte("<?xml version=\"1.0\"?>\n"))
	f.Write([]byte("<log>\n"))
	for rows.Next() {
		var l logDB
		err = rows.Scan(&l.LogType, &l.Timestamp, &l.Server, &l.TransactionNum, &l.Command, &l.Username, &l.StockSymbol, &l.Filename, &l.Funds, &l.Cryptokey, &l.Price, &l.QuoteServerTime, &l.Action, &l.ErrorMessage, &l.DebugMessage)
		writeToXML(f, l)
	}
	f.Write([]byte("</log>\n"))
	f.Close()
}

func dumpLogHandler(w http.ResponseWriter, r *http.Request) {
	logDumpCommand(w, r)
	f, err := os.Create("log.xml")
	if err != nil {
		failGracefully(err, "Failed to open log file ")
	}

	rows, err := db.Query("SELECT logType, (extract(EPOCH FROM timestamp) * 1000)::BIGINT as timestamp, server, transactionNum, command, username, stockSymbol, filename, (funds::DECIMAL)/100 as funds, cryptokey, (price::DECIMAL)/100 as price, quoteServerTime, action, errorMessage, debugMessage FROM audit_log;")
	if err != nil {
		failGracefully(err, "Failed to query audit DB ")
	}
	defer rows.Close()
	f.Write([]byte("<?xml version=\"1.0\"?>\n"))
	f.Write([]byte("<log>\n"))
	for rows.Next() {
		var l logDB
		err = rows.Scan(&l.LogType, &l.Timestamp, &l.Server, &l.TransactionNum, &l.Command, &l.Username, &l.StockSymbol, &l.Filename, &l.Funds, &l.Cryptokey, &l.Price, &l.QuoteServerTime, &l.Action, &l.ErrorMessage, &l.DebugMessage)
		writeToXML(f, l)
	}
	f.Write([]byte("</log>\n"))
	f.Close()
	fmt.Println("Log Dumped")
}

func dumpLogUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	res := struct {
		Username       string `json:"username"`
		TransactionNum int    `json:"transactionNum"`
		Filename       string `json:"filename"`
		Server         string `json:"server"`
	}{"", -1, FILENAME, SERVER}
	err := decoder.Decode(&res)
	if err != nil {
		failWithStatusCode(err, http.StatusText(http.StatusBadRequest), w, http.StatusBadRequest)
		return
	}

	queryString := "INSERT INTO audit_log(timestamp, transactionnum, server, command, filename, logtype) VALUES (now(), $1, $2, $3, $4, 'userCommand')"
	stmt, err := db.Prepare(queryString)

	resdb, err := stmt.Exec(res.TransactionNum, res.Server, "DUMPLOG", res.Filename)

	checkErrors(resdb, err, w)

	f, err := os.Create("log" + res.Username + ".xml")
	if err != nil {
		failGracefully(err, "Failed to open log file ")
	}

	queryString = "SELECT logType, (extract(EPOCH FROM timestamp) * 1000)::BIGINT as timestamp, server, transactionNum, command, username, stockSymbol, filename, (funds::DECIMAL)/100 as funds, cryptokey, (price::DECIMAL)/100 as price, quoteServerTime, action, errorMessage, debugMessage FROM audit_log WHERE username = $1;"
	stmt, err = db.Prepare(queryString)

	rows, err := stmt.Query(res.Username)
	if err != nil {
		failGracefully(err, "Failed to query audit DB ")
	}
	defer rows.Close()

	f.Write([]byte("<?xml version=\"1.0\"?>\n"))
	f.Write([]byte("<log>\n"))
	for rows.Next() {
		var l logDB
		err = rows.Scan(&l.LogType, &l.Timestamp, &l.Server, &l.TransactionNum, &l.Command, &l.Username, &l.StockSymbol, &l.Filename, &l.Funds, &l.Cryptokey, &l.Price, &l.QuoteServerTime, &l.Action, &l.ErrorMessage, &l.DebugMessage)
		writeToXML(f, l)
	}
	f.Write([]byte("</log>\n"))
	f.Close()
	fmt.Println("User Log Dumped")
}

func logDumpCommand(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	res := struct {
		Username       string `json:"username"`
		TransactionNum int    `json:"transactionNum"`
		Filename       string `json:"filename"`
		Server         string `json:"server"`
	}{"", -1, FILENAME, SERVER}
	err := decoder.Decode(&res)
	if err != nil {
		failWithStatusCode(err, http.StatusText(http.StatusBadRequest), w, http.StatusBadRequest)
		return
	}

	queryString := "INSERT INTO audit_log(timestamp, transactionnum, server, command, filename, logtype) VALUES (now(), $1, $2, $3, $4, 'userCommand')"
	stmt, err := db.Prepare(queryString)

	DBres, err := stmt.Exec(res.TransactionNum, res.Server, "DUMPLOG", res.Filename)
	checkErrors(DBres, err, w)
}
