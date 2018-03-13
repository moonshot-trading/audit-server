package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func checkErrors(res sql.Result, err error, w http.ResponseWriter) {
	if err != nil {
		failWithStatusCode(err, "Error inserting log ", w, http.StatusInternalServerError)
	}

	numRows, err := res.RowsAffected()

	if numRows < 1 {
		failWithStatusCode(err, "Unable to insert log entry ", w, http.StatusInternalServerError)
	}
}

func errorCheck(res sql.Result, err error) {
	if err != nil {
		failGracefully(err, "Error inserting log ")
	}

	numRows, err := res.RowsAffected()

	if numRows < 1 {
		failGracefully(err, "Unable to inserting log ")
	}

}

var (
	semaphoreChan = make(chan struct{}, 80)

	errStmt   string = "INSERT INTO audit_log(timestamp, server, transactionnum, command, username, stockSymbol, filename, funds, errorMessage, logType) VALUES (now(), $1, $2, $3, $4, $5, $6, $7, $8, 'errorEvent')"
	accStmt   string = "INSERT INTO audit_log(timestamp, server, transactionnum, action, username, funds, logtype) VALUES (now(), $1, $2, $3, $4, $5, 'accountTransaction')"
	quoteStmt string = "INSERT INTO audit_log(timestamp, server, transactionnum, price, stockSymbol, username, quoteServerTime, cryptokey, logType) VALUES (now(), $1, $2, $3, $4, $5, $6, $7, 'quoteServer')"
	userStmt  string = "INSERT INTO audit_log(timestamp, server, transactionNum, command, username, stockSymbol, filename, funds, logtype) VALUES (now(), $1, $2, $3, $4, $5, $6, $7, 'userCommand')"
)

// func userCommandHandler(w http.ResponseWriter, r *http.Request) {
// 	decoder := json.NewDecoder(r.Body)
// 	d := userCommand{}
// 	err := decoder.Decode(&d)

// 	if err != nil {
// 		failWithStatusCode(err, http.StatusText(http.StatusBadRequest), w, http.StatusBadRequest)
// 		return
// 	}

// 	queryString := "INSERT INTO audit_log(timestamp, server, transactionNum, command, username, stockSymbol, filename, funds, logtype) VALUES (now(), $1, $2, $3, $4, $5, $6, $7, 'userCommand')"

// 	stmt, err := db.Prepare(queryString)

// 	res, err := stmt.Exec(d.Server, d.TransactionNum, d.Command, d.Username, d.StockSymbol, d.Filename, d.Funds)

// 	checkErrors(res, err, w)
// }

func userCommandHandler(d []userCommand) {
	semaphoreChan <- struct{}{}
	go func() {
		bulkInsertUser(d)
		<-semaphoreChan
	}()
}

func bulkInsertUser(u []userCommand) {

	// txn, err := db.Begin()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// stmt, _ := txn.Prepare(pq.CopyIn("audit_log", "timestamp", "server", "transactionNum", "command", "username", "stockSymbol", "filename", "funds", "logtype"))
	// for _, d := range u {
	// 	res, err := stmt.Exec(time.Now(), d.Server, d.TransactionNum, d.Command, d.Username, d.StockSymbol, d.Filename, d.Funds, "userCommand")
	// 	errorCheck(res, err)
	// }
	// //res, err2 := stmt.Exec()
	// stmt.Exec()
	// err2 := stmt.Close()
	// err2 = txn.Commit()
	// //errorCheck(res, err2)
	// if err2 != nil {
	// 	failGracefully(err, "Error inserting log ")
	// }

	stmt, err := db.Prepare(userStmt)
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range u {
		res, err := stmt.Exec(d.Server, d.TransactionNum, d.Command, d.Username, d.StockSymbol, d.Filename, d.Funds)
		errorCheck(res, err)
	}

	stmt.Close()
}

// func userCommandHandler(d userCommand) {

// 	queryString := "INSERT INTO audit_log(timestamp, server, transactionNum, command, username, stockSymbol, filename, funds, logtype) VALUES (now(), $1, $2, $3, $4, $5, $6, $7, 'userCommand')"

// 	stmt, err := db.Prepare(queryString)

// 	res, err := stmt.Exec(d.Server, d.TransactionNum, d.Command, d.Username, d.StockSymbol, d.Filename, d.Funds)

// 	_ = stmt.Close()

// 	errorCheck(res, err)
// }

// func quoteServerHandler(w http.ResponseWriter, r *http.Request) {
// 	decoder := json.NewDecoder(r.Body)
// 	d := quoteServer{}
// 	err := decoder.Decode(&d)

// 	if err != nil {
// 		failWithStatusCode(err, http.StatusText(http.StatusBadRequest), w, http.StatusBadRequest)
// 		return
// 	}

// 	queryString := "INSERT INTO audit_log(timestamp, server, transactionnum, price, stockSymbol, username, quoteServerTime, cryptokey, logType) VALUES (now(), $1, $2, $3, $4, $5, $6, $7, 'quoteServer')"
// 	stmt, err := db.Prepare(queryString)

// 	res, err := stmt.Exec(d.Server, d.TransactionNum, d.Price, d.StockSymbol, d.Username, d.QuoteServerTime, d.Cryptokey)

// 	checkErrors(res, err, w)
// }

func bulkInsertQuote(u []quoteServer) {

	// txn, err := db.Begin()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// stmt, _ := txn.Prepare(pq.CopyIn("audit_log", "timestamp", "server", "transactionNum", "price", "stockSymbol", "username", "quoteServerTime", "cryptokey", "logType"))

	// for _, d := range u {
	// 	res, err := stmt.Exec(time.Now(), d.Server, d.TransactionNum, d.Price, d.StockSymbol, d.Username, d.QuoteServerTime, d.Cryptokey, "userCommand")
	// 	errorCheck(res, err)
	// }
	// //res, err2 := stmt.Exec()
	// stmt.Exec()
	// err2 := stmt.Close()
	// err2 = txn.Commit()
	// //errorCheck(res, err2)
	// if err2 != nil {
	// 	failGracefully(err, "Error inserting log ")
	// }
	stmt, err := db.Prepare(quoteStmt)
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range u {
		res, err := stmt.Exec(d.Server, d.TransactionNum, d.Price, d.StockSymbol, d.Username, d.QuoteServerTime, d.Cryptokey)
		errorCheck(res, err)
	}

	stmt.Close()
}

func quoteServerHandler(d []quoteServer) {
	semaphoreChan <- struct{}{}
	go func() {
		bulkInsertQuote(d)
		<-semaphoreChan
	}()
}

// func forceDumpQuote() {
// 	if len(quoteBulk) > 0 {
// 		tmp := make([]quoteServer, len(quoteBulk))
// 		copy(tmp, quoteBulk)
// 		go bulkInsertQuote(tmp)
// 		quoteBulk = nil
// 	}
// }

// func quoteServerHandler(d quoteServer) {

// 	queryString := "INSERT INTO audit_log(timestamp, server, transactionnum, price, stockSymbol, username, quoteServerTime, cryptokey, logType) VALUES (now(), $1, $2, $3, $4, $5, $6, $7, 'quoteServer')"
// 	stmt, err := db.Prepare(queryString)

// 	res, err := stmt.Exec(d.Server, d.TransactionNum, d.Price, d.StockSymbol, d.Username, d.QuoteServerTime, d.Cryptokey)

// 	_ = stmt.Close()

// 	errorCheck(res, err)
// }

// func accountTransactionHandler(w http.ResponseWriter, r *http.Request) {
// 	decoder := json.NewDecoder(r.Body)
// 	d := accountTransaction{}
// 	err := decoder.Decode(&d)

// 	if err != nil {
// 		failWithStatusCode(err, http.StatusText(http.StatusBadRequest), w, http.StatusBadRequest)
// 		return
// 	}

// 	queryString := "INSERT INTO audit_log(timestamp, server, transactionnum, action, username, funds, logtype) VALUES (now(), $1, $2, $3, $4, $5, 'accountTransaction')"
// 	stmt, err := db.Prepare(queryString)

// 	res, err := stmt.Exec(d.Server, d.TransactionNum, d.Action, d.Username, d.Funds)

// 	checkErrors(res, err, w)
// }

func bulkInsertTransaction(u []accountTransaction) {
	// fmt.Println("meme")
	// txn, err := db.Begin()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// stmt, _ := txn.Prepare(pq.CopyIn("audit_log", "timestamp", "server", "transactionNum", "action", "username", "funds", "logType"))

	// for _, d := range u {
	// 	_, err := stmt.Exec(time.Now(), d.Server, d.TransactionNum, d.Action, d.Username, d.Funds, "userCommand")
	// 	//errorCheck(res, err)
	// 	if err != nil {
	// 		failGracefully(err, "Error inserting log ")
	// 	}
	// }
	// //res, err2 := stmt.Exec()
	// stmt.Exec()
	// err2 := stmt.Close()
	// err2 = txn.Commit()
	// //errorCheck(res, err2)
	// if err2 != nil {
	// 	failGracefully(err, "Error inserting log ")
	// }
	//fmt.Println("meme", u)
	stmt, err := db.Prepare(accStmt)
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range u {
		//fmt.Println("FDasfdadaf", d.Server)
		res, err := stmt.Exec(d.Server, d.TransactionNum, d.Action, d.Username, d.Funds)
		errorCheck(res, err)
	}

	stmt.Close()

}

func accountTransactionHandler(d []accountTransaction) {
	// fmt.Println("meme1")
	// transactionBulk = append(transactionBulk, d)
	// if len(transactionBulk) > bulkAmount {
	// 	tmp := make([]accountTransaction, len(transactionBulk))
	// 	copy(tmp, transactionBulk)
	// 	go bulkInsertTransaction(tmp)
	// 	transactionBulk = nil
	// }
	semaphoreChan <- struct{}{}
	go func() {
		bulkInsertTransaction(d)
		<-semaphoreChan
	}()
}

// func accountTransactionHandler(d accountTransaction) {

// 	queryString := "INSERT INTO audit_log(timestamp, server, transactionnum, action, username, funds, logtype) VALUES (now(), $1, $2, $3, $4, $5, 'accountTransaction')"
// 	stmt, err := db.Prepare(queryString)

// 	res, err := stmt.Exec(d.Server, d.TransactionNum, d.Action, d.Username, d.Funds)

// 	_ = stmt.Close()

// 	errorCheck(res, err)
// }

func systemEventHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	d := systemEvent{}
	err := decoder.Decode(&d)

	if err != nil {
		failWithStatusCode(err, http.StatusText(http.StatusBadRequest), w, http.StatusBadRequest)
		return
	}

	queryString := "INSERT INTO audit_log(timestamp, server, transactionnum, command, username, stockSymbol, filename, funds, logType) VALUES (now(), $1, $2, $3, $4, $5, $6, $7, 'systemEvent')"
	stmt, err := db.Prepare(queryString)

	res, err := stmt.Exec(d.Server, d.TransactionNum, d.Command, d.Username, d.StockSymbol, d.Filename, d.Funds)

	_ = stmt.Close()

	checkErrors(res, err, w)
}

// func errorEventHandler(w http.ResponseWriter, r *http.Request){
// 	decoder := json.NewDecoder(r.Body)
// 	d := errorEvent{}
// 	err := decoder.Decode(&d)

// 	if err != nil {
// 		failWithStatusCode(err, http.StatusText(http.StatusBadRequest), w, http.StatusBadRequest)
// 		return
// 	}

// 	queryString := "INSERT INTO audit_log(timestamp, server, transactionnum, command, username, stockSymbol, filename, funds, errorMessage, logType) VALUES (now(), $1, $2, $3, $4, $5, $6, $7, $8, 'errorEvent')"
// 	stmt, err := db.Prepare(queryString)

// 	res, err := stmt.Exec(d.Server, d.TransactionNum, d.Command, d.Username, d.StockSymbol, d.Filename, d.Funds, d.ErrorMessage)

// 	checkErrors(res, err, w)
// }

func bulkInsertError(u []errorEvent) {

	// txn, err := db.Begin()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// stmt, _ := txn.Prepare(pq.CopyIn("audit_log", "timestamp", "server", "transactionNum", "command", "username", "stockSymbol", "filename", "funds", "errorMessage", "logType"))

	// for _, d := range u {
	// 	res, err := stmt.Exec(time.Now(), d.Server, d.TransactionNum, d.Command, d.Username, d.StockSymbol, d.Filename, d.Funds, d.ErrorMessage, "userCommand")
	// 	errorCheck(res, err)
	// }
	// //res, err2 := stmt.Exec()
	// stmt.Exec()
	// err2 := stmt.Close()
	// err2 = txn.Commit()
	// //errorCheck(res, err2)
	// if err2 != nil {
	// 	failGracefully(err, "Error inserting log ")
	// }

	stmt, err := db.Prepare(errStmt)
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range u {
		res, err := stmt.Exec(d.Server, d.TransactionNum, d.Command, d.Username, d.StockSymbol, d.Filename, d.Funds, d.ErrorMessage)
		errorCheck(res, err)
	}

	stmt.Close()
}

func errorEventHandler(d []errorEvent) {
	semaphoreChan <- struct{}{}
	go func() {
		bulkInsertError(d)
		<-semaphoreChan
	}()
}

// func forceDumpError() {
// 	if len(errorBulk) > 0 {
// 		tmp := make([]errorEvent, len(errorBulk))
// 		copy(tmp, errorBulk)
// 		go bulkInsertError(tmp)
// 		errorBulk = nil
// 	}
// }

// func errorEventHandler(d errorEvent) {

// 	queryString := "INSERT INTO audit_log(timestamp, server, transactionnum, command, username, stockSymbol, filename, funds, errorMessage, logType) VALUES (now(), $1, $2, $3, $4, $5, $6, $7, $8, 'errorEvent')"
// 	stmt, err := db.Prepare(queryString)

// 	res, err := stmt.Exec(d.Server, d.TransactionNum, d.Command, d.Username, d.StockSymbol, d.Filename, d.Funds, d.ErrorMessage)

// 	_ = stmt.Close()

// 	errorCheck(res, err)
// }

func debugEventHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	d := debugEvent{}
	err := decoder.Decode(&d)

	if err != nil {
		failWithStatusCode(err, http.StatusText(http.StatusBadRequest), w, http.StatusBadRequest)
		return
	}

	queryString := "INSERT INTO audit_log(timestamp, server, transactionnum, command, username, stockSymbol, filename, funds, debugMessage, logType) VALUES (now(), $1, $2, $3, $4, $5, $6, $7, $8, 'debugEvent')"
	stmt, err := db.Prepare(queryString)

	res, err := stmt.Exec(d.Server, d.TransactionNum, d.Command, d.Username, d.StockSymbol, d.Filename, d.Funds, d.DebugMessage)

	_ = stmt.Close()

	checkErrors(res, err, w)
}
