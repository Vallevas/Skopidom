-- Migration: 000005_move_and_repair.up.sql
-- Adds 'moved' to audit_events.action and 'in_repair' to items.status.

-- Extend audit action constraint.
ALTER TABLE audit_events
    DROP CONSTRAINT IF EXISTS audit_events_action_check;

ALTER TABLE audit_events
    ADD CONSTRAINT audit_events_action_check
    CHECK (action IN ('created', 'updated', 'disposed', 'moved'));

-- Extend item status constraint.
ALTER TABLE items
    DROP CONSTRAINT IF EXISTS items_status_check;

ALTER TABLE items
    ADD CONSTRAINT items_status_check
    CHECK (status IN ('active', 'disposed', 'in_repair'));
