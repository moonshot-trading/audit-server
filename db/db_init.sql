
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
  id                  serial PRIMARY KEY,
  time                TIMESTAMP,
  type                log_type,
  server              VARCHAR(10),
  t_id                INT,
  command             command,
  u_id                INT,
  user_name           VARCHAR(20),
  stock_symbol        VARCHAR(3),
  filename            VARCHAR(20),
  funds               money,
  cryptokey           VARCHAR(20),
  amount              NUMERIC,
  price               money,
  quote_server_time   INT,
  action              VARCHAR(20),
  error_message       VARCHAR(100),
  debug_message       VARCHAR(100)
);