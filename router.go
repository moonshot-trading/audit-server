package main

import (
	"net/http"
	"log"
)

func initRoutes() {
	port := ":44417"
	http.HandleFunc("/dumpLog", dumpLogHandler)
	http.HandleFunc("/userCommand", userCommandHandler)
	http.HandleFunc("/quoteServer", quoteServerHandler)
	http.HandleFunc("/accountTransaction", accountTransactionHandler)
	http.HandleFunc("/systemEvent", systemEventHandler)
	http.HandleFunc("/errorEvent", errorEventHandler)
	http.HandleFunc("/debugEvent", debugEventHandler)

	err := http.ListenAndServeTLS(port, "server.crt", "server.key", nil)

	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}
