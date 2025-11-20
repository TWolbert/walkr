-- name: GetRoleByName :one
SELECT id FROM roles
WHERE role = ?;
