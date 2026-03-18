-- queries/users.sql

-- name: CreateUser :one
INSERT INTO users (full_name, email, password_hash, role)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, updated_at;

-- name: GetUserByID :one
SELECT id, full_name, email, password_hash, role, created_at, updated_at
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, full_name, email, password_hash, role, created_at, updated_at
FROM users
WHERE email = $1;

-- name: ListUsers :many
SELECT id, full_name, email, password_hash, role, created_at, updated_at
FROM users
ORDER BY created_at ASC;

-- name: UpdateUser :one
UPDATE users
SET full_name = $1, role = $2, updated_at = NOW()
WHERE id = $3
RETURNING updated_at;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: EmailExists :one
SELECT EXISTS(
    SELECT 1 FROM users WHERE email = $1
) AS exists;

-- name: CountUsersByRole :one
SELECT COUNT(*) AS count
FROM users
WHERE role = $1;
