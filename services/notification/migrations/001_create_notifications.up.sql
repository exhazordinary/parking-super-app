-- Notification Service: Tables for notification management

CREATE TYPE notification_channel AS ENUM ('push', 'sms', 'email');
CREATE TYPE notification_status AS ENUM ('pending', 'sent', 'delivered', 'failed');
CREATE TYPE notification_priority AS ENUM ('low', 'normal', 'high');

-- Notifications table
CREATE TABLE notifications (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    channel notification_channel NOT NULL,
    type VARCHAR(100) NOT NULL,
    title VARCHAR(500) NOT NULL,
    body TEXT NOT NULL,
    data JSONB DEFAULT '{}',
    priority notification_priority NOT NULL DEFAULT 'normal',
    status notification_status NOT NULL DEFAULT 'pending',
    recipient VARCHAR(500) NOT NULL,
    provider_id VARCHAR(255),
    scheduled_at TIMESTAMPTZ,
    sent_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    failed_at TIMESTAMPTZ,
    error_msg TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Notification templates
CREATE TABLE notification_templates (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    channel notification_channel NOT NULL,
    type VARCHAR(100) NOT NULL,
    title VARCHAR(500) NOT NULL,
    body TEXT NOT NULL,
    variables TEXT[] DEFAULT '{}',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- User notification preferences
CREATE TABLE user_preferences (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE,
    push_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    sms_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    email_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    quiet_hours_start INT,
    quiet_hours_end INT,
    type_preferences JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_status ON notifications(status);
CREATE INDEX idx_notifications_created_at ON notifications(created_at DESC);
CREATE INDEX idx_notifications_scheduled ON notifications(scheduled_at) WHERE status = 'pending';
CREATE INDEX idx_templates_name ON notification_templates(name);
CREATE INDEX idx_templates_type ON notification_templates(type, channel);
CREATE INDEX idx_preferences_user_id ON user_preferences(user_id);
