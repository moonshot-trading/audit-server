
CREATE TYPE command AS ENUM (
  'ADD',
  'QUOTE',
  'BUY',
  'COMMIT_BUY',
  'CANCEL_BUY',
  'SELL',
  'COMMIT_SELL',
  'CANCEL_SELL',
  'SET_BUY_AMOUNT',
  'CANCEL_SET_BUY',
  'SET_BUY_TRIGGER',
  'SET_SELL_AMOUNT',
  'SET_SELL_TRIGGER',
  'CANCEL_SET_SELL',
  'DUMPLOG',
  'DISPLAY_SUMMARY'
);

CREATE TYPE log_type as ENUM (
  'userCommand',
  'quoteServer',
  'accountTransaction',
  'systemEvent',
  'errorEvent',
  'debug'
);

CREATE TABLE IF NOT EXISTS audit_log (
  id                serial PRIMARY KEY,
  timestamp         TIMESTAMP,
  logType           log_type,
  server            VARCHAR(20),
  transactionNum    INT,
  command           command,
  username          VARCHAR(20),
  stockSymbol       VARCHAR(3),
  filename          VARCHAR(20),
  funds             money,
  cryptokey         VARCHAR(20),
  price             money,
  quoteServerTime   BIGINT,
  action            VARCHAR(20),
  errorMessage      VARCHAR(100),
  debugMessage      VARCHAR(100)
);