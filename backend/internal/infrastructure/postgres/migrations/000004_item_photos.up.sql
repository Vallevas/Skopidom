-- Migration: 000004_item_photos.up.sql
-- Replaces the single photo_url field with a dedicated photos table.

CREATE TABLE item_photos (
    id         BIGSERIAL    PRIMARY KEY,
    item_id    BIGINT       NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    url        VARCHAR(512) NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_item_photos_item_id ON item_photos(item_id);

-- Migrate existing photos from items.photo_url.
INSERT INTO item_photos (item_id, url, created_at)
SELECT id, photo_url, created_at
FROM items
WHERE photo_url != '';

ALTER TABLE items DROP COLUMN photo_url;
