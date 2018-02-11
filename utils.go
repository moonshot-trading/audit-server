package main

import (
	"fmt"
	"net/http"
	"encoding/xml"
	"reflect"
	"database/sql"
	"strconv"
	"io"
	"os"
)

func runningInDocker() bool {
	_, err := os.Stat("/.dockerenv")
	if err == nil {
		return true
	}
	return false
}

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s", msg, err)
		panic(err)
	}
}

func failWithStatusCode(err error, msg string, w http.ResponseWriter, statusCode int) {
	failGracefully(err, msg)
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, msg)
}

func failGracefully(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s", msg, err)
	}
}

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
			} else {
				continue
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
