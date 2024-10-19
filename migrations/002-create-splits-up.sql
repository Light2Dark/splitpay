CREATE TABLE IF NOT EXISTS splits (
    id integer primary key,
    link text,
    receipt_id integer,
    FOREIGN KEY(receipt_id) REFERENCES receipts(id)
);