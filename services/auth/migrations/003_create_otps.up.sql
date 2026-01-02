-- Migration: Create OTPs table
-- Version: 003
-- Description: Stores one-time passwords for phone verification

CREATE TABLE otps (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    -- Phone number the OTP was sent to
    phone VARCHAR(15) NOT NULL,

    -- The OTP code (6 digits)
    -- In production, you might want to hash this too
    code VARCHAR(6) NOT NULL,

    -- When the OTP expires (typically 5 minutes)
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,

    -- Whether the OTP has been verified
    verified BOOLEAN NOT NULL DEFAULT FALSE,

    -- Number of failed verification attempts
    attempts INTEGER NOT NULL DEFAULT 0,

    -- When the OTP was created
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Index for looking up the latest OTP for a phone
CREATE INDEX idx_otps_phone_created ON otps(phone, created_at DESC);

-- Index for cleanup job (delete expired OTPs)
CREATE INDEX idx_otps_expires_at ON otps(expires_at);

COMMENT ON TABLE otps IS 'One-time passwords for phone verification';
COMMENT ON COLUMN otps.attempts IS 'Failed verification attempts (max 3)';
