-- name: GetRoleIDByName :one
SELECT id
FROM roles
WHERE partner_id IS NULL
  AND LOWER(name) = LOWER(sqlc.arg(role_name)::text);

-- name: CreateUser :one
INSERT INTO users (email, password_hash, name, role_id)
VALUES ($1, $2, $3, $4)
RETURNING id, email, name;

-- name: GetGlobalRoleIDByCode :one
SELECT id
FROM roles
WHERE partner_id IS NULL
  AND code = sqlc.arg(code)::text
  AND is_active;

-- name: GetGlobalDivisionIDByCode :one
SELECT id
FROM divisions
WHERE partner_id IS NULL
  AND code = sqlc.arg(code)::text
  AND is_active;
  
-- name: CreateAccount :one
WITH created AS (
    INSERT INTO users (
        partner_id, name, email, password_hash,
        role_id, division_id, unit_usaha_id, created_by
    )
    VALUES (
        sqlc.narg(partner_id)::int,
        sqlc.arg(name)::varchar(150),
        sqlc.arg(email)::varchar(255),
        sqlc.arg(password_hash)::text,
        sqlc.narg(role_id)::int,
        sqlc.narg(division_id)::int,
        sqlc.narg(unit_usaha_id)::int,
        sqlc.narg(created_by)::int
    )
    RETURNING *
),
logged AS (
    INSERT INTO audit_logs (
        partner_id, user_id, action, resource, resource_id,
        old_value, new_value, ip_address, actor_name, actor_email
    )
    SELECT
        created.partner_id,
        created.created_by,
        'USER_CREATED',
        'users',
        created.id,
        NULL::jsonb,
        jsonb_build_object(
            'name', created.name,
            'email', created.email,
            'role_id', created.role_id,
            'division_id', created.division_id,
            'unit_usaha_id', created.unit_usaha_id,
            'is_active', created.is_active
        ),
        sqlc.narg(ip_address)::varchar(45),
        actor.name,
        actor.email
    FROM created
    LEFT JOIN users actor ON actor.id = created.created_by
    RETURNING id
)
SELECT created.*
FROM created
CROSS JOIN (SELECT COUNT(*) FROM logged) AS audit_saved;

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

-- name: GetAccountForStatus :one
SELECT u.id,
       u.is_active,
       r.code AS role_code,
       COALESCE(r.is_active, TRUE)::bool AS role_active,
       COALESCE(d.is_active, TRUE)::bool AS division_active
FROM users u
LEFT JOIN roles r ON r.id = u.role_id
LEFT JOIN divisions d ON d.id = u.division_id
WHERE u.id = sqlc.arg(id)::int
  AND (sqlc.narg(partner_id)::int IS NULL OR u.partner_id = sqlc.narg(partner_id)::int);

-- name: SetAccountStatus :one
WITH target AS (
    SELECT
        u.id,
        u.partner_id,
        u.name,
        u.email,
        u.unit_usaha_id,
        u.is_active,
        r.code AS role_code,
        d.code AS division_code
    FROM users u
    LEFT JOIN roles r ON r.id = u.role_id
    LEFT JOIN divisions d ON d.id = u.division_id
    WHERE u.id = sqlc.arg(id)::int
      AND (
          sqlc.narg(partner_id)::int IS NULL
          OR u.partner_id = sqlc.narg(partner_id)::int
      )
    FOR UPDATE OF u
),
updated AS (
    UPDATE users u
    SET is_active = sqlc.arg(is_active)::bool,
        deactivated_at = CASE
            WHEN sqlc.arg(is_active)::bool THEN NULL
            ELSE NOW()
        END,
        updated_at = NOW()
    FROM target
    WHERE u.id = target.id
      AND target.is_active <> sqlc.arg(is_active)::bool
    RETURNING u.id, u.is_active
),
logged AS (
    INSERT INTO audit_logs (
        partner_id, user_id, action, resource, resource_id,
        old_value, new_value, ip_address, actor_name, actor_email
    )
    SELECT
        target.partner_id,
        sqlc.arg(actor_id)::int,
        CASE
            WHEN updated.is_active THEN 'USER_ACTIVATED'
            ELSE 'USER_DEACTIVATED'
        END,
        'users',
        target.id,
        jsonb_build_object(
            'name', target.name,
            'email', target.email,
            'unit_usaha_id', target.unit_usaha_id,
            'role', target.role_code,
            'division', target.division_code,
            'is_active', target.is_active
        ),
        jsonb_build_object(
            'is_active', updated.is_active
        ),
        sqlc.narg(ip_address)::varchar(45),
        actor.name,
        actor.email
    FROM target
    JOIN updated ON updated.id = target.id
    LEFT JOIN users actor ON actor.id = sqlc.arg(actor_id)::int
    RETURNING id
)
SELECT updated.id, updated.is_active
FROM updated
CROSS JOIN (SELECT COUNT(*) FROM logged) AS audit_saved;

-- name: GetAccountByID :one
SELECT
    u.id,
    u.name,
    u.email,
    u.partner_id,
    u.unit_usaha_id,
    u.is_active,
    u.created_at,
    r.code AS role_code,
    d.code AS division_code
FROM users u
LEFT JOIN roles r ON r.id = u.role_id
LEFT JOIN divisions d ON d.id = u.division_id
WHERE u.id = sqlc.arg(id)::int;

-- name: UpdateAccount :one
WITH target AS (
    SELECT
        u.id,
        u.partner_id,
        u.name,
        u.email,
        u.role_id,
        r.code AS role_code,
        u.division_id,
        d.code AS division_code,
        u.unit_usaha_id,
        u.is_active
    FROM users u
    LEFT JOIN roles r ON r.id = u.role_id
    LEFT JOIN divisions d ON d.id = u.division_id
    WHERE u.id = sqlc.arg(id)::int
    FOR UPDATE OF u
),
updated AS (
    UPDATE users u
    SET
        name = COALESCE(sqlc.narg(name)::varchar(150), target.name),
        role_id = COALESCE(sqlc.narg(role_id)::int, target.role_id),
        division_id = COALESCE(sqlc.narg(division_id)::int, target.division_id),
        unit_usaha_id = COALESCE(sqlc.narg(unit_usaha_id)::int, target.unit_usaha_id),
        updated_at = NOW()
    FROM target
    WHERE u.id = target.id
    RETURNING
        u.id, u.name, u.email, u.partner_id,
        u.role_id, u.division_id, u.unit_usaha_id,
        u.is_active, u.created_at
),
logged AS (
    INSERT INTO audit_logs (
        partner_id, user_id, action, resource, resource_id,
        old_value, new_value, ip_address, actor_name, actor_email
    )
    SELECT
        target.partner_id,
        sqlc.arg(actor_id)::int,
        'USER_UPDATED',
        'users',
        target.id,
        jsonb_build_object(
            'name', target.name,
            'email', target.email,
            'role', target.role_code,
            'division', target.division_code,
            'unit_usaha_id', target.unit_usaha_id,
            'is_active', target.is_active
        ),
        jsonb_build_object(
            'name', updated.name,
            'email', updated.email,
            'role', updated_role.code,
            'division', updated_division.code,
            'unit_usaha_id', updated.unit_usaha_id,
            'is_active', updated.is_active
        ),
        sqlc.narg(ip_address)::varchar(45),
        actor.name,
        actor.email
    FROM target
    JOIN updated ON updated.id = target.id
    LEFT JOIN roles updated_role ON updated_role.id = updated.role_id
    LEFT JOIN divisions updated_division ON updated_division.id = updated.division_id
    LEFT JOIN users actor ON actor.id = sqlc.arg(actor_id)::int
    RETURNING id
)
SELECT
    updated.id,
    updated.name,
    updated.email,
    updated.partner_id,
    updated.unit_usaha_id,
    updated.is_active,
    updated.created_at,
    r.code AS role_code,
    d.code AS division_code
FROM updated
LEFT JOIN roles r ON r.id = updated.role_id
LEFT JOIN divisions d ON d.id = updated.division_id
CROSS JOIN (SELECT COUNT(*) FROM logged) AS audit_saved;