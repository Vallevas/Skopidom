-- queries/item_relations.sql
-- SQL queries for managing symmetric item-to-item relations

-- name: CreateItemRelation :one
INSERT INTO item_relations (item_id_1, item_id_2, created_by)
VALUES ($1, $2, $3)
RETURNING id, item_id_1, item_id_2, created_at, created_by;

-- name: GetRelationsByItemID :many
SELECT id, item_id_1, item_id_2, created_at, created_by
FROM item_relations
WHERE item_id_1 = $1 OR item_id_2 = $1
ORDER BY created_at DESC;

-- name: DeleteItemRelation :exec
DELETE FROM item_relations
WHERE id = $1;

-- name: RelationExists :one
SELECT EXISTS(
    SELECT 1 FROM item_relations
    WHERE (item_id_1 = $1 AND item_id_2 = $2)
       OR (item_id_1 = $2 AND item_id_2 = $1)
) AS exists;

-- name: GetLinkedItems :many
-- Returns the related item details for a given item ID
SELECT 
    i.id,
    i.barcode,
    i.inventory_number,
    i.name,
    i.category_id,
    c.name AS category_name,
    i.room_id,
    rm.name AS room_name,
    rm.building_id,
    b.name AS building_name,
    b.address AS building_address,
    i.description,
    i.status,
    i.tx_hash,
    i.created_at,
    i.updated_at,
    i.pending_disposal_at,
    i.disposed_at,
    i.created_by,
    uc.full_name AS creator_full_name,
    uc.email AS creator_email,
    uc.role AS creator_role,
    uc.created_at AS creator_created_at,
    uc.updated_at AS creator_updated_at,
    i.last_edited_by,
    ue.full_name AS editor_full_name,
    ue.email AS editor_email,
    ue.role AS editor_role,
    ue.created_at AS editor_created_at,
    ue.updated_at AS editor_updated_at
FROM item_relations ir
JOIN items i ON (ir.item_id_1 = $1 AND ir.item_id_2 = i.id) 
             OR (ir.item_id_2 = $1 AND ir.item_id_1 = i.id)
JOIN categories c ON c.id = i.category_id
JOIN rooms rm ON rm.id = i.room_id
JOIN buildings b ON b.id = rm.building_id
JOIN users uc ON uc.id = i.created_by
JOIN users ue ON ue.id = i.last_edited_by
WHERE ir.item_id_1 = $1 OR ir.item_id_2 = $1
ORDER BY i.created_at DESC;
