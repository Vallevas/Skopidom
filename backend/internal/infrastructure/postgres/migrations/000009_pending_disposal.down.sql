-- Migration: 000009_pending_disposal.down.sql
-- Reverts pending_disposal status and drops disposal_documents table.

-- Drop disposal_documents table.
DROP TABLE IF EXISTS disposal_documents;

-- Remove timestamp columns.
ALTER TABLE items
    DROP COLUMN IF EXISTS pending_disposal_at,
    DROP COLUMN IF EXISTS disposed_at;

-- Revert item status constraint.
ALTER TABLE items
    DROP CONSTRAINT IF EXISTS items_status_check;

ALTER TABLE items
    ADD CONSTRAINT items_status_check
    CHECK (status IN ('active', 'disposed', 'in_repair'));

-- Revert audit action constraint.
ALTER TABLE audit_events
    DROP CONSTRAINT IF EXISTS audit_events_action_check;

ALTER TABLE audit_events
    ADD CONSTRAINT audit_events_action_check
    CHECK (action IN ('created', 'updated', 'disposed', 'moved', 'sent_to_repair', 'returned_from_repair'));
