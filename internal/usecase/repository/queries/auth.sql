-- name: GetUserByEmail :one
SELECT u.id,
       u.email,
       u.password_hash,
       u.name,
       u.is_active,
       u.is_super_admin,
       u.partner_id,
       u.unit_usaha_id,
       u.role_id,
       r.code AS role_code,
       d.code AS division_code
FROM users u
LEFT JOIN roles r ON r.id = u.role_id
LEFT JOIN divisions d ON d.id = u.division_id
WHERE u.email = $1;

-- name: GetAccountCredentials :one
SELECT email, password_hash, is_active
FROM users
WHERE id = $1;

-- name: UpdatePassword :one
UPDATE users
SET password_hash = sqlc.arg(new_hash)::text,
    updated_at = NOW()
WHERE id = sqlc.arg(id)::int
  AND password_hash = sqlc.arg(old_hash)::text
RETURNING updated_at;