-- name: CreateUserRole :one
INSERT INTO user_roles (
    role_id, user_id
) VALUES (
  ?, ?
)
RETURNING *;

-- name: GetUserRoleByUserId :one
SELECT roles.* FROM roles
INNER JOIN user_roles ON roles.id = user_roles.role_id
WHERE user_roles.user_id = ? LIMIT 1;
