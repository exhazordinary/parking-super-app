-- Parking Service: Tables for parking session management

CREATE TYPE session_status AS ENUM ('active', 'completed', 'cancelled', 'failed');

-- Parking sessions table
CREATE TABLE parking_sessions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    provider_id UUID NOT NULL,
    location_id UUID NOT NULL,
    external_session_id VARCHAR(255),
    vehicle_plate VARCHAR(20) NOT NULL,
    vehicle_type VARCHAR(50) NOT NULL DEFAULT 'car',
    entry_time TIMESTAMPTZ NOT NULL,
    exit_time TIMESTAMPTZ,
    duration_minutes INT DEFAULT 0,
    amount DECIMAL(19, 4) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'MYR',
    status session_status NOT NULL DEFAULT 'active',
    payment_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Vehicles table: Registered vehicles for users
CREATE TABLE vehicles (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    plate VARCHAR(20) NOT NULL,
    type VARCHAR(50) NOT NULL DEFAULT 'car',
    make VARCHAR(100),
    model VARCHAR(100),
    color VARCHAR(50),
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_sessions_user_id ON parking_sessions(user_id);
CREATE INDEX idx_sessions_provider_id ON parking_sessions(provider_id);
CREATE INDEX idx_sessions_status ON parking_sessions(status);
CREATE INDEX idx_sessions_entry_time ON parking_sessions(entry_time DESC);
CREATE INDEX idx_sessions_external_id ON parking_sessions(external_session_id);
CREATE INDEX idx_vehicles_user_id ON vehicles(user_id);
CREATE INDEX idx_vehicles_plate ON vehicles(plate);
