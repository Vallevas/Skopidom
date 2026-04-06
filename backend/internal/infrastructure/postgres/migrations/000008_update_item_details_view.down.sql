-- Migration: 000008_update_item_details_view.down.sql
-- Reverts item_details view to previous version without inventory_number

DROP VIEW IF EXISTS item_details;

CREATE VIEW item_details AS
SELECT
    i.id,
    i.barcode,
    i.name,

    i.category_id,
    c.name              AS category_name,

    i.room_id,
    rm.name             AS room_name,
    rm.building_id,
    b.name              AS building_name,
    b.address           AS building_address,

    i.description,
    i.status,
    i.tx_hash,

    i.created_at,
    i.updated_at,

    i.created_by,
    uc.full_name        AS creator_full_name,
    uc.email            AS creator_email,
    uc.role             AS creator_role,
    uc.created_at       AS creator_created_at,
    uc.updated_at       AS creator_updated_at,

    i.last_edited_by,
    ue.full_name        AS editor_full_name,
    ue.email            AS editor_email,
    ue.role             AS editor_role,
    ue.created_at       AS editor_created_at,
    ue.updated_at       AS editor_updated_at

FROM items i
JOIN categories c  ON c.id  = i.category_id
JOIN rooms      rm ON rm.id = i.room_id
JOIN buildings  b  ON b.id  = rm.building_id
JOIN users      uc ON uc.id = i.created_by
JOIN users      ue ON ue.id = i.last_edited_by;
