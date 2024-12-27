-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP NULL,
    role VARCHAR(255) NOT NULL DEFAULT 'user',
    first_name TEXT,
    last_name TEXT,
    email TEXT UNIQUE,
    mobile TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    code INTEGER,
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    user_type TEXT NOT NULL DEFAULT 'trial',
    expires_at TIMESTAMP NULL
);

-- Add indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_mobile ON users(mobile);
