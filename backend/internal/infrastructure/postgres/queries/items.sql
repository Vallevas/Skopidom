-- queries/items.sql
-- Item queries — all SELECTs go through the item_details view.
-- The view owns the JOIN logic; queries here stay intentionally simple.

-- name: GetItemByID :one
SELECT * FROM item_details
WHERE id = $1;

-- name: GetItemByBarcode :one
SELECT * FROM item_details
WHERE barcode = $1;

-- name: ListItems :many
SELECT * FROM item_details
WHERE
    (sqlc.narg('category_id')::bigint IS NULL OR category_id = sqlc.narg('category_id')) AND
    (sqlc.narg('room_id')::bigint     IS NULL OR room_id     = sqlc.narg('room_id'))     AND
    (sqlc.narg('status')::text        IS NULL OR status      = sqlc.narg('status'))      AND
    (sqlc.narg('date_from')::timestamptz IS NULL OR created_at >= sqlc.narg('date_from')) AND
    (sqlc.narg('date_to')::timestamptz   IS NULL OR created_at <= sqlc.narg('date_to'))
ORDER BY created_at DESC;

-- name: CreateItem :one
INSERT INTO items (
    barcode, inventory_number, name, category_id, room_id,
    description, status,
    created_by, last_edited_by
)
VALUES ($1, $2, $3, $4, $5, $6, 'active', $7, $7)
RETURNING id, created_at, updated_at;

-- name: UpdateItem :exec
UPDATE items
SET
    description    = $1,
    last_edited_by = $2,
    updated_at     = NOW()
WHERE id = $3;

-- name: UpdateItemStatus :exec
UPDATE items
SET
    status         = $1,
    last_edited_by = $2,
    updated_at     = NOW()
WHERE id = $3;

-- name: UpdateItemTxHash :exec
UPDATE items
SET tx_hash = $1
WHERE id = $2;

-- name: MoveItemToRoom :exec
UPDATE items
SET
    room_id        = $1,
    last_edited_by = $2,
    updated_at     = NOW()
WHERE id = $3;

-- name: BarcodeExists :one
SELECT EXISTS(
    SELECT 1 FROM items WHERE barcode = $1
) AS exists;

-- name: InventoryNumberExists :one
SELECT EXISTS(
    SELECT 1 FROM items WHERE inventory_number = $1
) AS exists;
