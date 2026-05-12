-- Migration: 000011_item_photos_base64.up.sql
-- Creates a separate table for storing Base64 image data to keep main queries lightweight.

CREATE TABLE item_photos_data (
    id         BIGSERIAL PRIMARY KEY,
    photo_id   BIGINT NOT NULL UNIQUE REFERENCES item_photos(id) ON DELETE CASCADE,
    data       BYTEA  NOT NULL,
    mime_type  VARCHAR(64) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_item_photos_data_photo_id ON item_photos_data(photo_id);

-- Migrate existing photos from local storage to Base64.
-- Note: This migration assumes photos are stored locally and converts them.
-- In production, you may need a separate script to fetch and convert existing files.
-- For now, we leave existing url entries as-is; new uploads will use Base64.
