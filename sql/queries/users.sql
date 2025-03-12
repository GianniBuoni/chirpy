-- name: CreateUser :one
INSERT INTO users (
  id, created_at, updated_at, email, hashed_password
)
VALUES (
  $1, $2, $3, $4, $5
)
RETURNING id, created_at, updated_at, email;

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserResponse :one
SELECT id, created_at, updated_at, email FROM users
WHERE id = $1;

-- name: DeleteUsers :exec
DELETE FROM users;
