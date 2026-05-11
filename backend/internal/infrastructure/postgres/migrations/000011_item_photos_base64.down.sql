-- Migration: 000011_item_photos_base64.down.sql
-- Reverts item_photos table back to URL-based storage.

DROP TABLE IF EXISTS item_photos CASCADE;

CREATE TABLE item_photos (
    id         BIGSERIAL    PRIMARY KEY,
    item_id    BIGINT       NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    url        VARCHAR(512) NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_item_photos_item_id ON item_photos(item_id);
