-- Provider Service: Tables for parking provider management

CREATE TYPE provider_status AS ENUM ('active', 'inactive', 'pending');
CREATE TYPE credential_environment AS ENUM ('sandbox', 'production');

-- Providers table: Parking provider organizations
CREATE TABLE providers (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(20) NOT NULL UNIQUE,
    description TEXT,
    logo_url VARCHAR(500),
    status provider_status NOT NULL DEFAULT 'pending',
    mfe_url VARCHAR(500) NOT NULL,
    api_base_url VARCHAR(500) NOT NULL,
    webhook_secret VARCHAR(255),
    config JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Provider credentials: API keys for provider authentication
CREATE TABLE provider_credentials (
    id UUID PRIMARY KEY,
    provider_id UUID NOT NULL REFERENCES providers(id) ON DELETE CASCADE,
    api_key VARCHAR(255) NOT NULL UNIQUE,
    api_secret VARCHAR(255) NOT NULL,
    environment credential_environment NOT NULL DEFAULT 'sandbox',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ
);

-- Locations: Parking locations operated by providers
CREATE TABLE locations (
    id UUID PRIMARY KEY,
    provider_id UUID NOT NULL REFERENCES providers(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    address VARCHAR(500) NOT NULL,
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL,
    postal_code VARCHAR(20),
    latitude DECIMAL(10, 8) NOT NULL,
    longitude DECIMAL(11, 8) NOT NULL,
    total_spaces INT DEFAULT 0,
    amenities TEXT[] DEFAULT '{}',
    hourly_rate DECIMAL(10, 2) NOT NULL DEFAULT 0,
    daily_max DECIMAL(10, 2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'MYR',
    grace_period_min INT NOT NULL DEFAULT 15,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_providers_code ON providers(code);
CREATE INDEX idx_providers_status ON providers(status);
CREATE INDEX idx_provider_credentials_api_key ON provider_credentials(api_key);
CREATE INDEX idx_provider_credentials_provider_id ON provider_credentials(provider_id);
CREATE INDEX idx_locations_provider_id ON locations(provider_id);
CREATE INDEX idx_locations_city ON locations(city);
CREATE INDEX idx_locations_coordinates ON locations(latitude, longitude);
