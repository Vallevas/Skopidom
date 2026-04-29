-- Migration: 000009_pending_disposal.up.sql
-- Adds 'pending_disposal' status and creates disposal_documents table.

-- Extend item status constraint to include pending_disposal.
ALTER TABLE items
    DROP CONSTRAINT IF EXISTS items_status_check;

ALTER TABLE items
    ADD CONSTRAINT items_status_check
    CHECK (status IN ('active', 'disposed', 'in_repair', 'pending_disposal'));

-- Add timestamp tracking for pending_disposal status.
ALTER TABLE items
    ADD COLUMN pending_disposal_at TIMESTAMPTZ,
    ADD COLUMN disposed_at         TIMESTAMPTZ;

-- Create disposal_documents table.
CREATE TABLE disposal_documents (
    id         BIGSERIAL    PRIMARY KEY,
    item_id    BIGINT       NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    filename   VARCHAR(255) NOT NULL,
    url        VARCHAR(512) NOT NULL,
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    uploaded_by BIGINT      NOT NULL REFERENCES users(id)
);

-- Index for quick lookup by item.
CREATE INDEX idx_disposal_documents_item_id ON disposal_documents(item_id);

-- Extend audit action constraint to include disposal actions.
ALTER TABLE audit_events
    DROP CONSTRAINT IF EXISTS audit_events_action_check;

ALTER TABLE audit_events
    ADD CONSTRAINT audit_events_action_check
    CHECK (action IN ('created', 'updated', 'disposed', 'moved', 'sent_to_repair', 'returned_from_repair', 'pending_disposal', 'disposal_finalized'));
