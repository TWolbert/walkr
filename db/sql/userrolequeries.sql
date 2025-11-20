-- name: CreateUserRole :one
INSERT INTO user_roles (
    role_id, user_id
) VALUES (
  ?, ?
)
RETURNING *;
