-- ============================================
-- QUERIES FOR USER MANAGEMENT
-- ============================================

-- name: CreateUser :one
INSERT INTO users (
    id, first_name, last_name, email, password, email_verified
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: Read :one
SELECT *
FROM users
WHERE id = $1 AND deleted_at IS NULL;
