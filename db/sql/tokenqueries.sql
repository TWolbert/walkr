-- name: CreateToken :one
INSERT INTO tokens (
  token, user_id, expires_at
) VALUES (
  ?, ?, ?
)
RETURNING *;

-- name: DeleteTokenById :exec
DELETE FROM tokens
WHERE id = ?;

-- name: DeleteTokenByUserId :exec
DELETE FROM tokens
WHERE user_id = ?;

-- name: GetTokenByUserId :one
SELECT * FROM tokens
WHERE user_id = ?;

-- name: GetUserFromToken :one
SELECT users.* FROM users
INNER JOIN tokens ON users.id = tokens.user_id
WHERE tokens.token = ? LIMIT 1;
