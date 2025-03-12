-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
  token, created_at, updated_at, user_id, expires_at
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT u.id, u.created_at, u.updated_at, u.email, rt.revoked_at
  FROM refresh_tokens rt
  INNER JOIN users u
  ON rt.user_id = u.id
  WHERE rt.token = $1;

-- name: RevokeToken :exec
UPDATE refresh_tokens
  SET updated_at = $2, revoked_at = $3
  WHERE token = $1;
