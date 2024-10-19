-- up
CREATE TABLE IF NOT EXISTS receipts (
  id integer primary key,
  items text,
  subtotal real,
  serviceCharge real,
  serviceChargePercent real,
  taxPercent real,
  taxAmount real,
  totalAmount real
);

-- down
DROP TABLE receipts;