-- name: CreateAuditLog :exec
INSERT INTO audit_logs (
    partner_id, user_id, action, resource, resource_id,
    old_value, new_value, ip_address, actor_name, actor_email
)
VALUES (
    sqlc.narg(partner_id)::int,
    sqlc.narg(user_id)::int,
    sqlc.arg(action)::varchar(80),
    sqlc.arg(resource)::varchar(50),
    sqlc.narg(resource_id)::int,
    sqlc.narg(old_value)::text::jsonb,
    sqlc.narg(new_value)::text::jsonb,
    sqlc.narg(ip_address)::varchar(45),
    sqlc.narg(actor_name)::varchar(150),
    sqlc.narg(actor_email)::varchar(255)
);

-- name: CountAuditLogs :one
SELECT COUNT(*)
FROM audit_logs
WHERE sqlc.arg(include_all)::bool
   OR partner_id = sqlc.narg(partner_id)::int;

-- name: ListAuditLogs :many
SELECT
    id,
    partner_id,
    user_id,
    action,
    resource,
    resource_id,
    old_value::text AS old_value,
    new_value::text AS new_value,
    ip_address,
    actor_name,
    actor_email,
    created_at
FROM audit_logs
WHERE sqlc.arg(include_all)::bool
   OR partner_id = sqlc.narg(partner_id)::int
ORDER BY created_at DESC, id DESC
LIMIT sqlc.arg(page_limit)::int
OFFSET sqlc.arg(page_offset)::int;