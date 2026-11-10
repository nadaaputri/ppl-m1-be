-- name: CountRoles :one
SELECT COUNT(*)
FROM roles
WHERE (sqlc.narg(partner_id)::int IS NULL OR partner_id = sqlc.narg(partner_id)::int)
  AND (sqlc.narg(is_active)::bool IS NULL OR is_active = sqlc.narg(is_active)::bool);

-- name: ListRoles :many
SELECT *
FROM roles
WHERE (sqlc.narg(partner_id)::int IS NULL OR partner_id = sqlc.narg(partner_id)::int)
  AND (sqlc.narg(is_active)::bool IS NULL OR is_active = sqlc.narg(is_active)::bool)
ORDER BY id
LIMIT sqlc.arg(page_limit)::int OFFSET sqlc.arg(page_offset)::int;

-- name: CreateRole :one
INSERT INTO roles (partner_id, code, name, description)
VALUES ($1, $2, $3, $4)
RETURNING *;