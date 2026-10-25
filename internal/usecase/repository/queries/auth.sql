-- name: GetUserByEmail :one
SELECT u.id,
       u.email,
       u.password_hash,
       u.full_name,
       u.is_active,
       r.name AS role_name,
       COALESCE(d.id::text, '')::text AS division_id,
       COALESCE(d.name, '') AS division_name
FROM users u
JOIN roles r ON r.id = u.role_id
LEFT JOIN divisions d ON d.id = u.division_id
WHERE u.email = $1;