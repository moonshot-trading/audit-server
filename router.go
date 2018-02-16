package main

import (
	"net/http"
)

func initRoutes() {
	port := ":44417"
	http.HandleFunc("/dumpLog", dumpLogHandler)
	http.HandleFunc("/dumpLogUser", dumpLogUserHandler)
	// http.HandleFunc("/userCommand", userCommandHandler)
	// http.HandleFunc("/quoteServer", quoteServerHandler)
	// http.HandleFunc("/accountTransaction", accountTransactionHandler)
	// http.HandleFunc("/systemEvent", systemEventHandler)
	// http.HandleFunc("/errorEvent", errorEventHandler)
	// http.HandleFunc("/debugEvent", debugEventHandler)

	http.ListenAndServe(port, nil)
}
