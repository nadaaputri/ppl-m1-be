-- name: CountRoles :one
SELECT COUNT(*)
FROM roles
WHERE partner_id IS NULL
  AND (sqlc.narg(is_active)::bool IS NULL
       OR is_active = sqlc.narg(is_active)::bool);

-- name: ListRoles :many
SELECT *
FROM roles
WHERE partner_id IS NULL
  AND (sqlc.narg(is_active)::bool IS NULL
       OR is_active = sqlc.narg(is_active)::bool)
ORDER BY id
LIMIT sqlc.arg(page_limit)::int
OFFSET sqlc.arg(page_offset)::int;

-- name: CreateRole :one
WITH created AS (
    INSERT INTO roles (partner_id, code, name, description)
    VALUES (
        NULL,
        sqlc.arg(code)::varchar(50),
        sqlc.arg(name)::varchar(100),
        sqlc.narg(description)::text
    )
    RETURNING *
),
assigned AS (
    INSERT INTO role_permissions (role_id, permission_code)
    SELECT created.id, selected.code
    FROM created
    CROSS JOIN LATERAL unnest(
        string_to_array(
            NULLIF(sqlc.arg(permission_codes)::text, ''),
            ','
        )
    ) AS selected(code)
    RETURNING role_id
)
SELECT created.*
FROM created
LEFT JOIN (
    SELECT DISTINCT role_id FROM assigned
) AS assigned_roles ON assigned_roles.role_id = created.id;

-- name: CountPermissions :one
SELECT COUNT(*)
FROM permissions;

-- name: ListPermissions :many
SELECT code, module, description
FROM permissions
ORDER BY module, code
LIMIT sqlc.arg(page_limit)::int
OFFSET sqlc.arg(page_offset)::int;

-- name: CountPermissionsByCodes :one
SELECT COUNT(*)
FROM permissions
WHERE code = ANY(
    string_to_array(
        NULLIF(sqlc.arg(permission_codes)::text, ''),
        ','
    )
);

-- name: GetGlobalRoleByID :one
SELECT *
FROM roles
WHERE id = sqlc.arg(id)::int
  AND partner_id IS NULL;

-- name: GetRolePermissionCodes :many
SELECT permission_code
FROM role_permissions
WHERE role_id = sqlc.arg(role_id)::int
ORDER BY permission_code;

-- name: SetRolePermissions :exec
WITH removed AS (
    DELETE FROM role_permissions
    WHERE role_id = sqlc.arg(role_id)::int
)
INSERT INTO role_permissions (role_id, permission_code)
SELECT sqlc.arg(role_id)::int, selected.code
FROM unnest(
    string_to_array(
        NULLIF(sqlc.arg(permission_codes)::text, ''),
        ','
    )
) AS selected(code);