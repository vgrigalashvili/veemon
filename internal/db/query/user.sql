-- name: AddUser :one
INSERT INTO users (
    id, created_at, updated_at, role, first_name, last_name, email, mobile, password, code, verified, user_type, expires_at
) VALUES (
    $1, $2, $3, COALESCE(sql.narg(1), 'user'), $4, $5,
    $6, $7, $8, COALESCE(sql.narg(2), 0), COALESCE(sql.narg(3), FALSE),
    COALESCE(sql.narg(4), 'trial'), $9
) RETURNING *;

-- name: FindUserByID :one
SELECT
    id,
    created_at,
    updated_at,
    COALESCE(role, sql.narg(1)) AS role,
    first_name,
    last_name,
    email,
    mobile,
    password,
    COALESCE(code, sql.narg(2)) AS code,
    COALESCE(verified, sql.narg(3)) AS verified,
    COALESCE(user_type, sql.narg(4)) AS user_type,
    expires_at
FROM users
WHERE id = $1;

-- name: FindUserByMobile :one
SELECT
    id,
    created_at,
    updated_at,
    COALESCE(role, sql.narg(1)) AS role,
    first_name,
    last_name,
    email,
    mobile,
    password,
    COALESCE(code, sql.narg(2)) AS code,
    COALESCE(verified, sql.narg(3)) AS verified,
    COALESCE(user_type, sql.narg(4)) AS user_type,
    expires_at
FROM users
WHERE mobile = $1;

-- name: FindUserByEmail :one
SELECT
    id,
    created_at,
    updated_at,
    COALESCE(role, sql.narg(1)) AS role,
    first_name,
    last_name,
    email,
    mobile,
    password,
    COALESCE(code, sql.narg(2)) AS code,
    COALESCE(verified, sql.narg(3)) AS verified,
    COALESCE(user_type, sql.narg(4)) AS user_type,
    expires_at
FROM users
WHERE email = $1;

-- name: GetAllUsers :many
SELECT
    id,
    created_at,
    updated_at,
    COALESCE(role, sql.narg(1)) AS role,
    first_name,
    last_name,
    email,
    mobile,
    password,
    COALESCE(code, sql.narg(2)) AS code,
    COALESCE(verified, sql.narg(3)) AS verified,
    COALESCE(user_type, sql.narg(4)) AS user_type,
    expires_at
FROM users;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: UpdateUser :one
UPDATE users
SET
    updated_at = $2,
    role = COALESCE(sql.narg(1), role),
    first_name = $3,
    last_name = $4,
    email = $5,
    mobile = $6,
    password = $7,
    code = COALESCE(sql.narg(2), code),
    verified = COALESCE(sql.narg(3), verified),
    user_type = COALESCE(sql.narg(4), user_type),
    expires_at = $8
WHERE id = $1
RETURNING *;
