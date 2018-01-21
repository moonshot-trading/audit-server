package main

import (
	"net/http"
	"encoding/json"
	"database/sql"
)

func checkErrors(res sql.Result, err error, w http.ResponseWriter){
	if err != nil {
		failWithStatusCode(err, "Error inserting log", w, http.StatusInternalServerError)
	}

	numRows, err := res.RowsAffected()

	if numRows < 1 {
		failWithStatusCode(err, "Unable to insert log entry", w, http.StatusInternalServerError)
	}
}

func userCommandHandler(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)
	d := userCommand{}
	err := decoder.Decode(&d)

	if err != nil {
		failWithStatusCode(err, http.StatusText(http.StatusBadRequest), w, http.StatusBadRequest)
		return
	}

	queryString := "INSERT INTO audit_log(timestamp, server, transactionNum, command, username, stockSymbol, filename, funds) VALUES (current_timestamp, $1, $2, $3, $4, $5, $6, $7)"
	stmt, err := db.Prepare(queryString)
	
	res, err := stmt.Exec(d.Timestamp, d.Server, d.TransactionNum, d.Command, d.Username, d.StockSymbol, d.Filename, d.Funds)

	checkErrors(res, err, w)
}

func quoteServerHandler(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)
	d := quoteServer{}
	err := decoder.Decode(&d)

	if err != nil {
		failWithStatusCode(err, http.StatusText(http.StatusBadRequest), w, http.StatusBadRequest)
		return
	}

	queryString := "INSERT INTO audit_log(timestamp, server, transactionNum, price, stockSymbol, username, quoteServerTime, cryptokey) VALUES (current_timestamp, $1, $2, $3, $4, $5, $6, $7)"
	stmt, err := db.Prepare(queryString)

	res, err := stmt.Exec(d.Timestamp, d.Server, d.TransactionNum, d.Price, d.StockSymbol, d.Username, d.QuoteServerTime, d.Cryptokey)

	checkErrors(res, err, w)
}

func accountTransactionHandler(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)
	d := accountTransaction{}
	err := decoder.Decode(&d)

	if err != nil {
		failWithStatusCode(err, http.StatusText(http.StatusBadRequest), w, http.StatusBadRequest)
		return
	}

	queryString := "INSERT INTO audit_log(timestamp, server, transactionNum, action, username, funds) VALUES (current_timestamp, $1, $2, $3, $4, $5)"
	stmt, err := db.Prepare(queryString)

	res, err := stmt.Exec(d.Timestamp, d.Server, d.TransactionNum, d.Action, d.Username, d.Funds)

	checkErrors(res, err, w)
}

func systemEventHandler(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)
	d := systemEvent{}
	err := decoder.Decode(&d)

	if err != nil {
		failWithStatusCode(err, http.StatusText(http.StatusBadRequest), w, http.StatusBadRequest)
		return
	}

	queryString := "INSERT INTO audit_log(timestamp, server, transactionNum, command, username, stockSymbol, filename, funds) VALUES (current_timestamp, $1, $2, $3, $4, $5, $6, $7)"
	stmt, err := db.Prepare(queryString)

	res, err := stmt.Exec(d.Timestamp, d.Server, d.TransactionNum, d.Command, d.Username, d.StockSymbol, d.Filename, d.Funds)

	checkErrors(res, err, w)
}

func errorEventHandler(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)
	d := errorEvent{}
	err := decoder.Decode(&d)

	if err != nil {
		failWithStatusCode(err, http.StatusText(http.StatusBadRequest), w, http.StatusBadRequest)
		return
	}

	queryString := "INSERT INTO audit_log(timestamp, server, transactionNum, command, username, stockSymbol, filename, funds, errorMessage) VALUES (current_timestamp, $1, $2, $3, $4, $5, $6, $7, $8)"
	stmt, err := db.Prepare(queryString)

	res, err := stmt.Exec(d.Timestamp, d.Server, d.TransactionNum, d.Command, d.Username, d.StockSymbol, d.Filename, d.Funds, d.ErrorMessage)

	checkErrors(res, err, w)
}

func debugEventHandler(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)
	d := debugEvent{}
	err := decoder.Decode(&d)

	if err != nil {
		failWithStatusCode(err, http.StatusText(http.StatusBadRequest), w, http.StatusBadRequest)
		return
	}

	queryString := "INSERT INTO audit_log(timestamp, server, transactionNum, command, username, stockSymbol, filename, funds, debugMessage) VALUES (current_timestamp, $1, $2, $3, $4, $5, $6, $7, $8)"
	stmt, err := db.Prepare(queryString)

	res, err := stmt.Exec(d.Timestamp, d.Server, d.TransactionNum, d.Command, d.Username, d.StockSymbol, d.Filename, d.Funds, d.DebugMessage)

	checkErrors(res, err, w)
}
