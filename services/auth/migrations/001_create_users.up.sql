-- Migration: Create users table
-- Version: 001
-- Description: Initial users table for authentication
--
-- MICROSERVICES PATTERN: Database Per Service
-- ============================================
-- Each microservice owns its own database. The auth service has
-- exclusive access to the auth_db database.
--
-- Benefits:
-- - Services are loosely coupled
-- - Can use different database technologies per service
-- - Independent scaling and schema evolution
-- - Better fault isolation

-- Enable UUID extension for generating UUIDs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create users table
CREATE TABLE users (
    -- Primary key: UUID is better than auto-increment for distributed systems
    -- Benefits:
    -- - Can be generated client-side (no DB round-trip)
    -- - No sequential pattern (harder to guess/enumerate)
    -- - Works well with database sharding
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    -- Phone number is the primary identifier for Malaysian users
    -- Format: +60XXXXXXXXXX
    phone VARCHAR(15) NOT NULL UNIQUE,

    -- Email is optional but useful for notifications
    email VARCHAR(255),

    -- Password hash (bcrypt produces ~60 character strings)
    password_hash VARCHAR(255) NOT NULL,

    -- User's display name
    full_name VARCHAR(255) NOT NULL,

    -- Account status: pending, active, inactive, banned
    status VARCHAR(20) NOT NULL DEFAULT 'pending',

    -- Timestamps for auditing
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for common query patterns
-- INDEX STRATEGY:
-- - Phone lookup: Used for login (very frequent)
-- - Email lookup: Used for password reset (less frequent)
-- - Status filter: Used for admin queries

CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_email ON users(email) WHERE email IS NOT NULL;
CREATE INDEX idx_users_status ON users(status);

-- Create updated_at trigger
-- This automatically updates the updated_at column when a row is modified
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments for documentation
COMMENT ON TABLE users IS 'User accounts for the parking super app';
COMMENT ON COLUMN users.phone IS 'Malaysian phone number in +60 format';
COMMENT ON COLUMN users.status IS 'Account status: pending, active, inactive, banned';
