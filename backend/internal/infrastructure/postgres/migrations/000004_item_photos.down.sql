-- Migration: 000004_item_photos.down.sql

ALTER TABLE items ADD COLUMN photo_url VARCHAR(512) NOT NULL DEFAULT '';

UPDATE items i
SET photo_url = (
    SELECT url FROM item_photos p
    WHERE p.item_id = i.id
    ORDER BY p.created_at ASC
    LIMIT 1
);

DROP TABLE IF EXISTS item_photos;
