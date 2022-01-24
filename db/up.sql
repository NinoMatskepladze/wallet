DROP TABLE IF EXISTS wallets;
CREATE TABLE wallets (
  id VARCHAR(64) PRIMARY KEY,
  balance INT CHECK(balance >= 0),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

DROP TABLE IF EXISTS transactions;
CREATE TABLE transactions (
  id VARCHAR(64) PRIMARY KEY,
  wallet_id VARCHAR(64) NOT NULL,
  amount INT NOT NULL,
  balance INT CHECK(balance >= 0),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  FOREIGN KEY (wallet_id) REFERENCES wallets (id)
);