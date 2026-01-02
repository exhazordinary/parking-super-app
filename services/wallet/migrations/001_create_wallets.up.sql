-- Wallet Service: Core tables for wallet and transaction management
-- This migration creates the foundational tables for the digital wallet system

-- Create enum types for wallet and transaction statuses
CREATE TYPE wallet_status AS ENUM ('active', 'inactive', 'frozen');
CREATE TYPE transaction_type AS ENUM ('topup', 'payment', 'refund', 'transfer');
CREATE TYPE transaction_status AS ENUM ('pending', 'completed', 'failed', 'refunded');

-- Wallets table: Stores user wallet information
-- Each user has exactly one wallet (1:1 relationship enforced by unique constraint)
CREATE TABLE wallets (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE,
    balance DECIMAL(19, 4) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'MYR',
    status wallet_status NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Ensure balance is never negative (enforced at DB level for safety)
    CONSTRAINT positive_balance CHECK (balance >= 0)
);

-- Transactions table: Immutable ledger of all wallet transactions
-- Uses idempotency_key to prevent duplicate transactions
CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    wallet_id UUID NOT NULL REFERENCES wallets(id),
    type transaction_type NOT NULL,
    amount DECIMAL(19, 4) NOT NULL,
    balance_before DECIMAL(19, 4) NOT NULL,
    balance_after DECIMAL(19, 4) NOT NULL,
    reference_id VARCHAR(255),
    provider_id UUID,
    status transaction_status NOT NULL DEFAULT 'pending',
    description TEXT,
    idempotency_key VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Ensure transaction amounts are positive
    CONSTRAINT positive_amount CHECK (amount > 0)
);

-- Payment methods table: Stores saved payment methods (cards, bank accounts)
CREATE TABLE payment_methods (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    type VARCHAR(50) NOT NULL,
    provider VARCHAR(100) NOT NULL,
    token VARCHAR(500) NOT NULL,
    last_four VARCHAR(4),
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for common query patterns
CREATE INDEX idx_wallets_user_id ON wallets(user_id);
CREATE INDEX idx_transactions_wallet_id ON transactions(wallet_id);
CREATE INDEX idx_transactions_created_at ON transactions(created_at DESC);
CREATE INDEX idx_transactions_idempotency_key ON transactions(idempotency_key) WHERE idempotency_key IS NOT NULL;
CREATE INDEX idx_transactions_reference_id ON transactions(reference_id) WHERE reference_id IS NOT NULL;
CREATE INDEX idx_payment_methods_user_id ON payment_methods(user_id);

-- Unique constraint on idempotency key (only when not null)
CREATE UNIQUE INDEX idx_transactions_unique_idempotency ON transactions(idempotency_key) WHERE idempotency_key IS NOT NULL AND idempotency_key != '';
