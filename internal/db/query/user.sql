-- name: CreateUser :one
INSERT INTO users (
    id, first_name, last_name, email, mobile, password_hash, role, user_type, code, verified, expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1 AND deleted_at IS NULL;

-- name: GetUserByMobile :one
SELECT *
FROM users
WHERE mobile = $1 AND deleted_at IS NULL;

-- name: UpdateUser :exec
UPDATE users
SET
    first_name = COALESCE($2, first_name),
    last_name = COALESCE($3, last_name),
    email = COALESCE($4, email),
    mobile = COALESCE($5, mobile),
    password_hash = COALESCE($6, password_hash),
    role = COALESCE($7, role),
    user_type = COALESCE($8, user_type),
    code = COALESCE($9, code),
    verified = COALESCE($10, verified),
    updated_at = now(),
    expires_at = COALESCE($11, expires_at)
WHERE id = $1 AND deleted_at IS NULL;

-- name: SoftDeleteUser :exec
UPDATE users
SET deleted_at = now()
WHERE id = $1;

-- name: ListUsers :many
SELECT *
FROM users
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: VerifyUser :exec
UPDATE users
SET verified = true, updated_at = now()
WHERE id = $1;

-- name: ResetUserCode :exec
UPDATE users
SET code = $2, updated_at = now()
WHERE id = $1;
