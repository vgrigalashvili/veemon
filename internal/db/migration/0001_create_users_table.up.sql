-- ============================================
-- USERS TABLE DEFINITION
-- ============================================
CREATE TABLE users (
    id UUID PRIMARY KEY,                         			-- Unique identifier for each user
    created_at TIMESTAMP NOT NULL DEFAULT now(),			-- Timestamp when the user was created
    updated_at TIMESTAMP NOT NULL DEFAULT now(),       		-- Timestamp when the user was last updated
    deleted_at TIMESTAMP NULL,                         		-- Timestamp for soft deletion
    first_name TEXT,                                   		-- User's first name
    last_name TEXT,                                    		-- User's last name
    type TEXT NOT NULL DEFAULT 'individual',           		-- User type (e.g., individual, company)
    role VARCHAR(255) NOT NULL DEFAULT 'user',         		-- Role of the user (e.g., user, admin)
    email TEXT UNIQUE,                                 		-- Unique email address
    mobile TEXT NOT NULL UNIQUE,                      		-- Unique mobile number
    password_hash TEXT NOT NULL,                      		-- Hashed password
    pin INTEGER CHECK (pin >= 100000 AND pin <= 999999),	-- 6-digit One-Time Pin (OTP)
    verified BOOLEAN NOT NULL DEFAULT FALSE,          		-- Indicates if the user is verified
    expires_at TIMESTAMP NULL                         		-- Expiration time for certain operations or accounts
);

-- ============================================
-- INDEXES FOR PERFORMANCE
-- ============================================
CREATE INDEX idx_users_email ON users(email);          		-- Index for fast email lookups
CREATE INDEX idx_users_mobile ON users(mobile);        		-- Index for fast mobile lookups
CREATE INDEX idx_users_verified ON users(verified);    		-- Index for filtering verified users
CREATE INDEX idx_users_expires_at ON users(expires_at);		-- Index for filtering expired accounts