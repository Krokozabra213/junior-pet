-- name: CreateUser :one
INSERT INTO users (
    username,
    email,
    password,
    name,
    surname,
    is_male
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING id, created_at, updated_at;

-- name: GetUserByUsername :one
SELECT 
    id,
    username,
    email,
    password,
    name,
    surname,
    is_male,
    created_at,
    updated_at
FROM users
WHERE username = $1
  AND deleted_at IS NULL;

-- name: GetUserByID :one
SELECT 
    id,
    username,
    email,
    password,
    name,
    surname,
    is_male,
    created_at,
    updated_at
FROM users
WHERE id = $1
  AND deleted_at IS NULL;

-- name: GetUserByEmail :one
SELECT 
    id,
    username,
    email,
    password,
    name,
    surname,
    is_male,
    created_at,
    updated_at
FROM users
WHERE email = $1
  AND deleted_at IS NULL;

-- name: UpdatePassword :exec
UPDATE users
SET
    password   = $2,
    updated_at = NOW()
WHERE id = $1
  AND deleted_at IS NULL;

-- name: SoftDeleteUser :exec
UPDATE users
SET
    deleted_at = NOW(),
    updated_at = NOW()
WHERE id = $1
  AND deleted_at IS NULL;

-- name: HardDeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: CountUsers :one
SELECT COUNT(*) FROM users
WHERE deleted_at IS NULL;

-- name: ExistsUserByUsername :one
SELECT EXISTS(
    SELECT 1 FROM users
    WHERE username = $1
      AND deleted_at IS NULL
);

-- name: ExistsUserByEmail :one
SELECT EXISTS(
    SELECT 1 FROM users
    WHERE email = $1
      AND deleted_at IS NULL
);
