-- Migration: 000011_add_item_photos_data.up.sql
-- Creates a separate table for storing Base64-encoded photo data.
-- This keeps heavy binary data out of the main item_photos table.

CREATE TABLE item_photos_data (
    id         BIGSERIAL PRIMARY KEY,
    photo_id   BIGINT       NOT NULL UNIQUE REFERENCES item_photos(id) ON DELETE CASCADE,
    data       BYTEA        NOT NULL,
    mime_type  VARCHAR(64)  NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_item_photos_data_photo_id ON item_photos_data(photo_id);

-- Note: Existing photos will have NULL data until re-uploaded.
-- The application should handle missing data gracefully.
