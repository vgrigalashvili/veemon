-- ============================================
-- QUERIES FOR USER MANAGEMENT
-- ============================================

-- Add a new user
-- Inserts a new user record and returns the full record.
-- Parameters: id, first_name, last_name, type, role, email, mobile, password_hash, pin, verified, expires_at
-- name: AddUser :one
INSERT INTO users (
    id, first_name, last_name, type, role, email, mobile, password_hash, pin, verified, expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- Retrieve user by ID
-- Fetches a user by ID, excluding soft-deleted users.
-- Parameters: id
-- name: WhoIsBID :one
SELECT *
FROM users
WHERE id = $1 AND deleted_at IS NULL;

-- Retrieve user by email
-- Fetches a user by email, excluding soft-deleted users.
-- Parameters: email
-- name: WhoIsBEmail :one
SELECT *
FROM users
WHERE email = $1 AND deleted_at IS NULL;

-- Retrieve user by mobile
-- Fetches a user by mobile, excluding soft-deleted users.
-- Parameters: mobile
-- name: WhoIsBMobile :one
SELECT *
FROM users
WHERE mobile = $1 AND deleted_at IS NULL;

-- Update user details
-- Updates user fields only if new values are provided.
-- Parameters: id, first_name, last_name, email, mobile, password_hash, role, type, pin, verified, expires_at
-- name: UpdateUser :exec
UPDATE users
SET
    first_name = COALESCE($2, first_name),
    last_name = COALESCE($3, last_name),
    email = COALESCE($4, email),
    mobile = COALESCE($5, mobile),
    password_hash = COALESCE($6, password_hash),
    role = COALESCE($7, role),
    type = COALESCE($8, type),
    pin = COALESCE($9, pin),
    verified = COALESCE($10, verified),
    expires_at = COALESCE($11, expires_at),
    updated_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- Soft delete user
-- Marks a user as deleted by setting the deleted_at timestamp.
-- Parameters: id
-- name: SoftDeleteUser :exec
UPDATE users
SET deleted_at = now()
WHERE id = $1;

-- List all active users
-- Fetches a paginated list of active users, ordered by creation date.
-- Parameters: limit, offset
-- name: ListUsers :many
SELECT *
FROM users
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- Verify user
-- Marks a user as verified.
-- Parameters: id
-- name: VerifyUser :exec
UPDATE users
SET verified = true, updated_at = now()
WHERE id = $1;

-- Reset user's pin
-- Updates a user's PIN with a new value.
-- Parameters: id, pin
-- name: ResetUserPin :exec
UPDATE users
SET pin = $2, updated_at = now()
WHERE id = $1;

-- Update user's expiration time
-- Sets or updates the expiration time for a user.
-- Parameters: id, expires_at
-- name: UserExpiresAt :exec
UPDATE users
SET expires_at = $2, updated_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- Retrieve user's role
-- Fetches only the role of a user by ID, excluding soft-deleted users.
-- Parameters: id
-- name: GetUserRole :one
SELECT role
FROM users
WHERE id = $1 AND deleted_at IS NULL;

-- Update user's role
-- Updates the role of a user.
-- Parameters: id, role
-- name: SetupUserRole :exec
UPDATE users
SET role = $2, updated_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- List soft-deleted users
-- Fetches a paginated list of soft-deleted users for recovery or auditing.
-- Parameters: limit, offset
-- name: ListSoftDeletedUsers :many
SELECT *
FROM users
WHERE deleted_at IS NOT NULL
ORDER BY deleted_at DESC
LIMIT $1 OFFSET $2;

-- Reactivate soft-deleted user
-- Removes the deleted_at timestamp to restore a soft-deleted user.
-- Parameters: id
-- name: ReactivateUser :exec
UPDATE users
SET deleted_at = NULL, updated_at = now()
WHERE id = $1;

-- Search users
-- Searches for users by name, email, or mobile, excluding soft-deleted users.
-- Parameters: search_term, limit, offset
-- name: SearchUsers :many
SELECT *
FROM users
WHERE deleted_at IS NULL
  AND (email ILIKE '%' || $1 || '%' OR mobile ILIKE '%' || $1 || '%' OR first_name ILIKE '%' || $1 || '%' OR last_name ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;