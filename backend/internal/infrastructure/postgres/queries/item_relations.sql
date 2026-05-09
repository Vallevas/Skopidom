-- queries/item_relations.sql
-- Item relations queries for linking items together

-- name: ListItemRelations :many
SELECT
    ir.id,
    ir.item_id,
    ir.related_item_id,
    ir.created_at,
    ir.created_by,
    ri.id AS "related_item:id",
    ri.barcode AS "related_item:barcode",
    ri.inventory_number AS "related_item:inventory_number",
    ri.name AS "related_item:name",
    ri.category_id AS "related_item:category_id",
    ri.room_id AS "related_item:room_id",
    ri.description AS "related_item:description",
    ri.status AS "related_item:status",
    ri.tx_hash AS "related_item:tx_hash",
    ri.created_at AS "related_item:created_at",
    ri.updated_at AS "related_item:updated_at",
    ri.created_by AS "related_item:created_by",
    ri.last_edited_by AS "related_item:last_edited_by",
    ri.pending_disposal_at AS "related_item:pending_disposal_at",
    ri.disposed_at AS "related_item:disposed_at"
FROM item_relations ir
JOIN items ri ON ri.id = ir.related_item_id
WHERE ir.item_id = $1;

-- name: CreateItemRelation :one
INSERT INTO item_relations (item_id, related_item_id, created_by)
VALUES ($1, $2, $3)
RETURNING id, created_at;

-- name: DeleteItemRelation :exec
DELETE FROM item_relations
WHERE item_id = $1 AND related_item_id = $2;

-- name: GetItemRelation :one
SELECT id, item_id, related_item_id, created_at, created_by
FROM item_relations
WHERE item_id = $1 AND related_item_id = $2;
