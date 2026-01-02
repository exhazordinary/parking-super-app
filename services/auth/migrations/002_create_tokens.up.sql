-- Migration: Create refresh_tokens table
-- Version: 002
-- Description: Stores refresh tokens for JWT token rotation

CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    -- Reference to the user who owns this token
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- SHA-256 hash of the actual token (we never store the raw token)
    -- SECURITY: If DB is compromised, attackers can't use these hashes
    token_hash VARCHAR(64) NOT NULL UNIQUE,

    -- When the token expires (typically 7 days from creation)
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,

    -- Whether the token has been revoked (logout, token rotation)
    revoked BOOLEAN NOT NULL DEFAULT FALSE,

    -- When the token was created
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    -- When the token was revoked (null if not revoked)
    revoked_at TIMESTAMP WITH TIME ZONE,

    -- Security metadata for tracking sessions
    user_agent TEXT,       -- Browser/app that created this token
    ip_address VARCHAR(45) -- IPv4 or IPv6 address
);

-- Indexes for common query patterns
-- Token lookup by hash (for validation)
CREATE INDEX idx_tokens_hash ON refresh_tokens(token_hash);

-- Find all tokens for a user (for "logout everywhere")
CREATE INDEX idx_tokens_user_id ON refresh_tokens(user_id);

-- Find expired tokens (for cleanup job)
CREATE INDEX idx_tokens_expires_at ON refresh_tokens(expires_at);

-- Find active tokens (not revoked, not expired)
CREATE INDEX idx_tokens_active ON refresh_tokens(user_id, revoked, expires_at)
    WHERE revoked = FALSE;

COMMENT ON TABLE refresh_tokens IS 'Stores refresh token hashes for JWT rotation';
COMMENT ON COLUMN refresh_tokens.token_hash IS 'SHA-256 hash of the refresh token';
