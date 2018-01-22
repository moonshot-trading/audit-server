package main

import (
	"net/http"
	"os"
	"encoding/xml"
	"reflect"
	"database/sql"
	"strconv"
	"fmt"
	"io"
)

func (s logEntry) MarshalXML( e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = s["logType"]
	delete(s, "logType")
	err := e.EncodeToken(start)
	if err != nil { failGracefully(err, "could not encode token ") }
	for k, v := range s { e.Encode(xmlMapEntry{XMLName: xml.Name{Local: k}, Value: v}) }
	return e.EncodeToken(start.End())
}

func structToMap(i interface{}) (values logEntry){
	values = make(logEntry)
	iVal := reflect.ValueOf(i).Elem()
	typ := iVal.Type()
	for i := 0; i <iVal.NumField(); i++ {
		f := iVal.Field(i)
		tag := typ.Field(i).Tag.Get("paramName")
		var v string

		switch f.Interface().(type) {
		case sql.NullString:
			if f.Field(1).Bool(){
				v = f.Field(0).String()
			} else {
				continue
			}
		case sql.NullInt64:
			if f.Field(1).Bool(){
				v = strconv.FormatInt(f.Field(0).Int(), 10)
			} else {
				continue
			}
		case sql.NullFloat64:
			if f.Field(1).Bool(){
				v = strconv.FormatFloat(f.Field(0).Float(), 'f', 2, 64)
			}
		default:
			v = fmt.Sprint(f.Interface())
		}

		values[tag] = v
	}
	return values
}

func writeToXML(w io.Writer, r logDB) {
	var output []byte
	var err error

	rMap := structToMap(&r)
	output, err = xml.MarshalIndent(rMap, "    ", "    ")
	if err != nil { failGracefully(err, "Failed to write to file ") }

	w.Write(output)
	w.Write([]byte("\n"))
}

func dumpLogHandler(w http.ResponseWriter, r *http.Request) {

	f, err := os.Create("log/log.xml")
	if err != nil { failGracefully(err, "Failed to open log file ") }
	defer f.Close()

	rows, err := db.Query("SELECT logType, (extract(EPOCH FROM timestamp) * 1000)::BIGINT as timestamp, server, transactionNum, command, username, stockSymbol, filename, (funds::DECIMAL)/100, cryptokey, (price::DECIMAL)/100, quoteServerTime, action, errorMessage, debugMessage FROM audit_log;")
	if err != nil { failGracefully(err, "Failed to query audit DB ") }
	defer rows.Close()

	f.Write([]byte("<log>\n"))
	for rows.Next() {
		var r logDB
		err = rows.Scan(&r.LogType, &r.Timestamp, &r.Server, &r.TransactionNum, &r.Command, &r.Username, &r.StockSymbol, &r.Filename, &r.Funds, &r.Cryptokey, &r.Price, &r.QuoteServerTime, &r.Action, &r.ErrorMessage, &r.DebugMessage)
		writeToXML(f, r)
	}
	f.Write([]byte("</log>\n"))

}
