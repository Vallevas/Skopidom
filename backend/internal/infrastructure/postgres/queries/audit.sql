-- queries/audit.sql

-- name: CreateAuditEvent :one
INSERT INTO audit_events (item_id, actor_id, action, payload)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at;

-- name: ListAuditEventsByItem :many
SELECT
    ae.id,
    ae.item_id,
    ae.actor_id,
    u.full_name  AS actor_full_name,
    u.email      AS actor_email,
    u.role       AS actor_role,
    ae.action,
    ae.payload,
    ae.tx_hash,
    ae.created_at
FROM audit_events ae
JOIN users u ON u.id = ae.actor_id
WHERE ae.item_id = $1
ORDER BY ae.created_at ASC;

-- name: UpdateAuditEventTxHash :exec
UPDATE audit_events
SET tx_hash = $1
WHERE id = $2;
