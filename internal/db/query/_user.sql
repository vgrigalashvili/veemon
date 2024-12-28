-- Create a new user
-- name: CreateUser :exec
INSERT INTO users (
    id,
    created_at,
    updated_at,
    role,
    first_name,
    last_name,
    email,
    mobile,
    password,
    code,
    verified,
    user_type,
    expires_at
) VALUES (
    $1, -- UUID of the user
    $2, -- Timestamp of creation
    $3, -- Timestamp of last update
    $4, -- Role (e.g., 'admin', 'user')
    $5, -- First name
    $6, -- Last name
    $7, -- Email
    $8, -- Mobile phone number
    $9, -- Hashed password
    $10, -- Verification or identification code
    $11, -- Verification status (true/false)
    $12, -- User type (e.g., 'trial')
    $13  -- Expiration date (nullable)
);

-- Get a user by ID
-- name: GetUserByID :one
SELECT
    id,
    created_at,
    updated_at,
    role,
    first_name,
    last_name,
    email,
    mobile,
    password,
    code,
    verified,
    user_type,
    expires_at
FROM users
WHERE id = $1; -- UUID of the user

-- Get all users
-- name: GetAllUsers :many
SELECT
    id,
    created_at,
    updated_at,
    role,
    first_name,
    last_name,
    email,
    mobile,
    password,
    code,
    verified,
    user_type,
    expires_at
FROM users;

-- Check if a user exists by mobile
-- name: CheckUserExistsByMobile :one
SELECT
    id
FROM users
WHERE mobile = $1; -- Mobile phone number

-- Check if a user exists by email
-- name: CheckUserExistsByEmail :one
SELECT
    id
FROM users
WHERE email = $1; -- Email address

-- Update a user by ID
-- name: UpdateUser :exec
UPDATE users
SET
    updated_at = COALESCE($2, updated_at),  -- Timestamp of last update
    role = COALESCE($3, role),              -- Role (e.g., 'admin', 'user')
    first_name = COALESCE($4, first_name),  -- First name
    last_name = COALESCE($5, last_name),    -- Last name
    email = COALESCE($6, email),            -- Email
    mobile = COALESCE($7, mobile),          -- Mobile phone number
    password = COALESCE($8, password),      -- Hashed password
    code = COALESCE($9, code),              -- Verification or identification code
    verified = COALESCE($10, verified),     -- Verification status (true/false)
    user_type = COALESCE($11, user_type),   -- User type (e.g., 'trial')
    expires_at = COALESCE($12, expires_at)  -- Expiration date (nullable)
WHERE id = $1;                              -- UUID of the user

-- Delete a user by ID (Soft Delete)
-- name: SoftDeleteUser :exec
UPDATE users
SET deleted_at = $2 -- Timestamp of deletion
WHERE id = $1;      -- UUID of the user

-- Permanently delete a user by ID
-- name: HardDeleteUser :exec
DELETE FROM users
WHERE id = $1;      -- UUID of the user
