CREATE TABLE IF NOT EXISTS transactions (
    transaction_id TEXT PRIMARY KEY,
    amount         NUMERIC(19, 4) NOT NULL,
    store_id       TEXT NOT NULL,
    timestamp      TIMESTAMPTZ NOT NULL,
    customer_id    TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS customers (
    customer_id      TEXT PRIMARY KEY,
    identifier_type  TEXT NOT NULL,
    customer_active  BOOLEAN NOT NULL DEFAULT TRUE
);
