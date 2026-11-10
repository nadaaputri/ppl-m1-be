-- name: CountDivisions :one
SELECT COUNT(*)
FROM divisions
WHERE (sqlc.narg(partner_id)::int IS NULL OR partner_id = sqlc.narg(partner_id)::int)
  AND (sqlc.narg(is_active)::bool IS NULL OR is_active = sqlc.narg(is_active)::bool);

-- name: ListDivisions :many
SELECT *
FROM divisions
WHERE (sqlc.narg(partner_id)::int IS NULL OR partner_id = sqlc.narg(partner_id)::int)
  AND (sqlc.narg(is_active)::bool IS NULL OR is_active = sqlc.narg(is_active)::bool)
ORDER BY id
LIMIT sqlc.arg(page_limit)::int OFFSET sqlc.arg(page_offset)::int;

-- name: CreateDivision :one
INSERT INTO divisions (partner_id, code, name)
VALUES ($1, $2, $3)
RETURNING *;