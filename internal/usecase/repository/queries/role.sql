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
    INSERT INTO roles (
        partner_id, code, name, description
    )
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
    RETURNING role_id, permission_code
),
logged AS (
    INSERT INTO audit_logs (
        partner_id, user_id, action, resource, resource_id,
        old_value, new_value, ip_address, actor_name, actor_email
    )
    SELECT
        NULL,
        sqlc.arg(actor_id)::int,
        'ROLE_CREATED',
        'roles',
        created.id,
        NULL::jsonb,
        jsonb_build_object(
            'code', created.code,
            'name', created.name,
            'description', created.description,
            'permission_codes',
            COALESCE(
                array_agg(assigned.permission_code)
                FILTER (WHERE assigned.permission_code IS NOT NULL),
                ARRAY[]::text[]
            )
        ),
        sqlc.narg(ip_address)::varchar(45),
        actor.name,
        actor.email
    FROM created
    LEFT JOIN assigned ON assigned.role_id = created.id
    LEFT JOIN users actor ON actor.id = sqlc.arg(actor_id)::int
    GROUP BY created.id, actor.name, actor.email
    RETURNING id
)
SELECT created.*
FROM created
CROSS JOIN (SELECT COUNT(*) FROM logged) AS audit_saved;

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
SELECT rp.permission_code
FROM role_permissions rp
JOIN roles r ON r.id = rp.role_id
WHERE rp.role_id = sqlc.arg(role_id)::int
  AND r.partner_id IS NULL
ORDER BY rp.permission_code;

-- name: SetRolePermissions :exec
WITH target_role AS (
    SELECT id
    FROM roles
    WHERE id = sqlc.arg(role_id)::int
      AND partner_id IS NULL
),
old_permissions AS (
    SELECT COALESCE(
        array_agg(permission_code ORDER BY permission_code),
        ARRAY[]::text[]
    ) AS codes
    FROM role_permissions
    WHERE role_id = sqlc.arg(role_id)::int
),
removed AS (
    DELETE FROM role_permissions
    WHERE role_id = sqlc.arg(role_id)::int
    RETURNING role_id
),
assigned AS (
    INSERT INTO role_permissions (role_id, permission_code)
    SELECT target_role.id, selected.code
    FROM target_role
    CROSS JOIN (SELECT COUNT(*) FROM removed) AS removal_done
    CROSS JOIN LATERAL unnest(
        string_to_array(
            NULLIF(sqlc.arg(permission_codes)::text, ''),
            ','
        )
    ) AS selected(code)
    RETURNING permission_code
)
INSERT INTO audit_logs (
    partner_id, user_id, action, resource, resource_id,
    old_value, new_value, ip_address, actor_name, actor_email
)
SELECT
    NULL,
    sqlc.arg(actor_id)::int,
    'ROLE_PERMISSIONS_UPDATED',
    'roles',
    target_role.id,
    jsonb_build_object(
        'permission_codes', old_permissions.codes
    ),
    jsonb_build_object(
        'permission_codes',
        COALESCE(
            array_agg(assigned.permission_code)
            FILTER (WHERE assigned.permission_code IS NOT NULL),
            ARRAY[]::text[]
        )
    ),
    sqlc.narg(ip_address)::varchar(45),
    actor.name,
    actor.email
FROM target_role
CROSS JOIN old_permissions
LEFT JOIN assigned ON TRUE
LEFT JOIN users actor ON actor.id = sqlc.arg(actor_id)::int
GROUP BY target_role.id, old_permissions.codes, actor.name, actor.email;