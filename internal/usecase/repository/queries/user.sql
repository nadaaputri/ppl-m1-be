-- name: GetRoleIDByName :one
SELECT id
FROM roles
WHERE LOWER(name) = LOWER(sqlc.arg(role_name)::text);

-- name: CreateUser :one
INSERT INTO users (email, password_hash, name, role_id)
VALUES ($1, $2, $3, $4)
RETURNING id, email, name;

-- name: GetPartnerRoleID :one
SELECT id
FROM roles
WHERE partner_id = sqlc.arg(partner_id)::int
  AND code = sqlc.arg(code)::text
  AND is_active;

-- name: GetPartnerDivisionID :one
SELECT id
FROM divisions
WHERE partner_id = sqlc.arg(partner_id)::int
  AND code = sqlc.arg(code)::text
  AND is_active;

-- name: CreateAccount :one
INSERT INTO users (partner_id, name, email, password_hash, role_id, division_id, unit_usaha_id, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;