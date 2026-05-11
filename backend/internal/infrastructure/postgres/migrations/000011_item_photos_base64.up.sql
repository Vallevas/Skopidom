-- Migration: 000011_item_photos_base64.up.sql
-- Changes item_photos table to store Base64 encoded image data instead of URL paths.

-- Drop the existing table and recreate with base64_data column
DROP TABLE IF EXISTS item_photos CASCADE;

CREATE TABLE item_photos (
    id         BIGSERIAL    PRIMARY KEY,
    item_id    BIGINT       NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    base64_data TEXT        NOT NULL,
    mime_type  VARCHAR(50)  NOT NULL DEFAULT 'image/jpeg',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_item_photos_item_id ON item_photos(item_id);
