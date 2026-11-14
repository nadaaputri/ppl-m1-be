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
WITH created AS (
    INSERT INTO divisions (partner_id, code, name)
    VALUES (
        NULL,
        sqlc.arg(code)::varchar(50),
        sqlc.arg(name)::varchar(100)
    )
    RETURNING *
),
logged AS (
    INSERT INTO audit_logs (
        partner_id, user_id, action, resource, resource_id,
        old_value, new_value, ip_address, actor_name, actor_email
    )
    SELECT
        NULL,
        sqlc.arg(actor_id)::int,
        'DIVISION_CREATED',
        'divisions',
        created.id,
        NULL::jsonb,
        jsonb_build_object(
            'code', created.code,
            'name', created.name
        ),
        sqlc.narg(ip_address)::varchar(45),
        actor.name,
        actor.email
    FROM created
    LEFT JOIN users actor ON actor.id = sqlc.arg(actor_id)::int
    RETURNING id
)
SELECT created.*
FROM created
CROSS JOIN (SELECT COUNT(*) FROM logged) AS audit_saved;