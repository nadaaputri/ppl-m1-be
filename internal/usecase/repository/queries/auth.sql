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
WITH updated AS (
    UPDATE users
    SET password_hash = sqlc.arg(new_hash)::text,
        updated_at = NOW()
    WHERE id = sqlc.arg(id)::int
      AND password_hash = sqlc.arg(old_hash)::text
    RETURNING id, partner_id, name, email, updated_at
),
logged AS (
    INSERT INTO audit_logs (
        partner_id, user_id, action, resource, resource_id,
        old_value, new_value, ip_address, actor_name, actor_email
    )
    SELECT
        updated.partner_id,
        updated.id,
        'PASSWORD_CHANGED',
        'users',
        updated.id,
        NULL::jsonb,
        NULL::jsonb,
        sqlc.narg(ip_address)::varchar(45),
        updated.name,
        updated.email
    FROM updated
    RETURNING id
)
SELECT updated.updated_at
FROM updated
CROSS JOIN (SELECT COUNT(*) FROM logged) AS audit_saved;