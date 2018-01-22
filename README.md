# audit-server

[![Build Status](https://travis-ci.org/moonshot-trading/audit-server.svg?branch=master)](https://travis-ci.org/moonshot-trading/audit-server)

<!-- Any time a user's account is touched, an account message is printed. Appropriate actions are "add" or "remove". -->
/accountTransaction
{
	"server":"string",
	"action": "add" or "remove",
	"username": "string",
	"funds": int  //cents
}

<!-- System events can be current user commands, interserver communications, or the execution of previously set triggers -->
/systemEvent
{
	"server":"string",
	"command": "string", //see user command list
	"stockSymbol": "string", //max three chars
	"username": "string"
	"filename": "string",
	"funds": int  //cents
}

<!-- Error messages contain all the information of user commands, in addition to an optional error message -->
/errorEvent
{
	"server":"string",
	"command": "string", //see user command list
	"stockSymbol": "string", //max three chars
	"filename": "string",
	"funds": int, //cents
	"username": "string",
	"errorMessage": "string"
}

<!-- Debugging messages contain all the information of user commands, in addition to an optional debug message -->
/debugEvent
{
	"server":"string",
	"command": "string", //see user command list
	"stockSymbol": "string", //max three chars
	"filename": "string",
	"funds": int, //cents
	"username": "string",
	"debugMessage": "string"
}

<!-- Every hit to the quote server requires a log entry with the results. The price, symbol, username, timestamp and cryptokey are as returned by the quote server -->
/quoteServer
{
	"server":"string",
	"price": int, //cents
	"stockSymbol": "string", //max three chars
	"username": "string",
	"quoteServerTime": int, //unix time min=1514764800000 max=1525132800000
	"cryptokey": "string"
}

<!-- User commands come from the user command files or from manual entries in the students' web forms -->
/userCommand
{
	"server":"string",
	"command": "string", //see user command list
	"username": "string",
	"stockSymbol": "string", //max three chars
	"filename": "string",
	"funds": int, //cents
}