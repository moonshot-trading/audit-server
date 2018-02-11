DO
$body$
BEGIN
  IF NOT EXISTS (
      SELECT                       -- SELECT list can stay empty for this
      FROM   pg_catalog.pg_user
      WHERE  usename = 'moonshot') THEN

    CREATE ROLE moonshot LOGIN PASSWORD 'hodl';
  END IF;
END
$body$;

DROP DATABASE IF EXISTS  "moonshot-audit";
CREATE DATABASE "moonshot-audit";
GRANT ALL PRIVILEGES ON DATABASE "moonshot-audit" TO moonshot;

\connect moonshot-audit

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
  'debugEvent'
);

CREATE TABLE IF NOT EXISTS audit_log (
  t_id              SERIAL PRIMARY KEY,
  transactionNum    INT,
  timestamp         TIMESTAMP,
  logType           log_type,
  server            VARCHAR(20),
  command           command,
  username          VARCHAR(20),
  stockSymbol       VARCHAR(3),
  filename          VARCHAR(20),
  funds             INT,
  cryptokey         VARCHAR(64),
  price             INT,
  quoteServerTime   BIGINT,
  action            VARCHAR(20),
  errorMessage      VARCHAR(100),
  debugMessage      VARCHAR(100)
);
