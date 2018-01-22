package main
import (
	"database/sql"
	"encoding/xml"
)

type logEntry map[string]string

type xmlMapEntry struct {
	XMLName xml.Name
	Value string `xml:",chardata"`
}

type logDB struct {
	LogType			string				`paramName:"logType"`
	Timestamp		int64				`paramName:"timestamp"`
	Server			string				`paramName:"server"`
	TransactionNum	int					`paramName:"transactionNum"`
	Command			sql.NullString		`paramName:"command"`
	Username		sql.NullString		`paramName:"username"`
	StockSymbol		sql.NullString		`paramName:"stockSymbol"`
	Filename		sql.NullString		`paramName:"filename"`
	Cryptokey		sql.NullString		`paramName:"cryptokey"`
	Action			sql.NullString		`paramName:"action"`
	ErrorMessage	sql.NullString		`paramName:"errorMessage"`
	DebugMessage	sql.NullString		`paramName:"debugMessage"`
	Funds			sql.NullFloat64		`paramName:"funds"`
	Price			sql.NullFloat64		`paramName:"price"`
	QuoteServerTime	sql.NullInt64		`paramName:"quoteServerTime"`
}

type userCommand struct {
	Timestamp		int64		`json:"timestamp"`
	Server			string		`json:"server"`
	TransactionNum	int			`json:"transactionNum"`
	Command			string		`json:"command"`
	Username		string		`json:"username"`
	StockSymbol		string		`json:"stockSymbol"`
	Filename		string		`json:"filename"`
	Funds			float64		`json:"funds"`
}

type quoteServer struct {
	Timestamp		int64		`json:"timestamp"`
	Server			string		`json:"server"`
	TransactionNum	int			`json:"transactionNum"`
	Price			float64		`json:"price"`
	StockSymbol		string		`json:"stockSymbol"`
	Username		string		`json:"username"`
	QuoteServerTime	int64		`json:"quoteServerTime"`
	Cryptokey		string		`json:"cryptokey"`
}

type accountTransaction struct {
	Timestamp		int64		`json:"timestamp"`
	Server			string		`json:"server"`
	TransactionNum	int			`json:"transactionNum"`
	Action			string		`json:"action"`
	Username		string		`json:"username"`
	Funds			float64		`json:"funds"`
}

type systemEvent struct {
	Timestamp		int64		`json:"timestamp"`
	Server			string		`json:"server"`
	TransactionNum	int			`json:"transactionNum"`
	Command			string		`json:"command"`
	Username		string		`json:"username"`
	StockSymbol		string		`json:"stockSymbol"`
	Filename		string		`json:"filename"`
	Funds			float64		`json:"funds"`
}

type errorEvent	struct {
	Timestamp		int64		`json:"timestamp"`
	Server			string		`json:"server"`
	TransactionNum	int			`json:"transactionNum"`
	Command			string		`json:"command"`
	Username		string		`json:"username"`
	StockSymbol		string		`json:"stockSymbol"`
	Filename		string		`json:"filename"`
	Funds			float64		`json:"funds"`
	ErrorMessage	string		`json:"errorMessage"`
}

type debugEvent	struct {
	Timestamp		int64		`json:"timestamp"`
	Server			string		`json:"server"`
	TransactionNum	int			`json:"transactionNum"`
	Command			string		`json:"command"`
	Username		string		`json:"username"`
	StockSymbol		string		`json:"stockSymbol"`
	Filename		string		`json:"filename"`
	Funds			float64		`json:"funds"`
	DebugMessage	string		`json:"debugMessage"`
}
