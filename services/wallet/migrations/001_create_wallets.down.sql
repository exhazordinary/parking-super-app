-- Rollback wallet service tables
DROP TABLE IF EXISTS payment_methods;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS wallets;

DROP TYPE IF EXISTS transaction_status;
DROP TYPE IF EXISTS transaction_type;
DROP TYPE IF EXISTS wallet_status;
