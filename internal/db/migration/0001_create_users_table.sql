-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP NULL,
    role VARCHAR(255) NOT NULL DEFAULT 'user',
    first_name TEXT,
    last_name TEXT,
    email CITEXT UNIQUE,
    mobile TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    code INTEGER CHECK (code >= 100000 AND code <= 999999), -- Example 6-digit OTP
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    user_type TEXT NOT NULL DEFAULT 'trial',
    expires_at TIMESTAMP NULL
);

-- Add indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_mobile ON users(mobile);
CREATE INDEX idx_users_verified ON users(verified);
CREATE INDEX idx_users_expires_at ON users(expires_at);
