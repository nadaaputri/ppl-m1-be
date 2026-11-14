-- name: CountDivisions :one
SELECT COUNT(*)
FROM divisions
WHERE partner_id IS NULL
  AND (sqlc.narg(is_active)::bool IS NULL
       OR is_active = sqlc.narg(is_active)::bool);

-- name: ListDivisions :many
SELECT *
FROM divisions
WHERE partner_id IS NULL
  AND (sqlc.narg(is_active)::bool IS NULL
       OR is_active = sqlc.narg(is_active)::bool)
ORDER BY id
LIMIT sqlc.arg(page_limit)::int
OFFSET sqlc.arg(page_offset)::int;

-- name: CreateDivision :one
INSERT INTO divisions (partner_id, code, name)
VALUES (
    NULL,
    sqlc.arg(code)::varchar(50),
    sqlc.arg(name)::varchar(100)
)
RETURNING *;