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

-- name: CountAccounts :one
SELECT COUNT(*)
FROM users u
LEFT JOIN roles r ON r.id = u.role_id
LEFT JOIN divisions d ON d.id = u.division_id
WHERE (sqlc.narg(partner_id)::int IS NULL OR u.partner_id = sqlc.narg(partner_id)::int)
  AND (sqlc.narg(search)::text IS NULL
       OR u.name ILIKE '%' || sqlc.narg(search)::text || '%'
       OR u.email ILIKE '%' || sqlc.narg(search)::text || '%')
  AND (sqlc.narg(role_code)::text IS NULL OR r.code = sqlc.narg(role_code)::text)
  AND (sqlc.narg(division_code)::text IS NULL OR d.code = sqlc.narg(division_code)::text)
  AND (sqlc.narg(unit_usaha_id)::int IS NULL OR u.unit_usaha_id = sqlc.narg(unit_usaha_id)::int)
  AND (sqlc.narg(is_active)::bool IS NULL OR u.is_active = sqlc.narg(is_active)::bool);

-- name: ListAccounts :many
SELECT u.id, u.name, u.email, u.partner_id, u.unit_usaha_id, u.is_active, u.created_at,
       r.code AS role_code,
       d.code AS division_code
FROM users u
LEFT JOIN roles r ON r.id = u.role_id
LEFT JOIN divisions d ON d.id = u.division_id
WHERE (sqlc.narg(partner_id)::int IS NULL OR u.partner_id = sqlc.narg(partner_id)::int)
  AND (sqlc.narg(search)::text IS NULL
       OR u.name ILIKE '%' || sqlc.narg(search)::text || '%'
       OR u.email ILIKE '%' || sqlc.narg(search)::text || '%')
  AND (sqlc.narg(role_code)::text IS NULL OR r.code = sqlc.narg(role_code)::text)
  AND (sqlc.narg(division_code)::text IS NULL OR d.code = sqlc.narg(division_code)::text)
  AND (sqlc.narg(unit_usaha_id)::int IS NULL OR u.unit_usaha_id = sqlc.narg(unit_usaha_id)::int)
  AND (sqlc.narg(is_active)::bool IS NULL OR u.is_active = sqlc.narg(is_active)::bool)
ORDER BY
  CASE WHEN sqlc.arg(sort)::text = 'name' THEN u.name END ASC,
  CASE WHEN sqlc.arg(sort)::text = '-name' THEN u.name END DESC,
  CASE WHEN sqlc.arg(sort)::text = 'email' THEN u.email END ASC,
  CASE WHEN sqlc.arg(sort)::text = '-email' THEN u.email END DESC,
  CASE WHEN sqlc.arg(sort)::text = 'created_at' THEN u.created_at END ASC,
  u.created_at DESC,
  u.id DESC
LIMIT sqlc.arg(page_limit)::int OFFSET sqlc.arg(page_offset)::int;