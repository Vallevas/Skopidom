-- Migration: 000006_repair_audit_actions.down.sql

ALTER TABLE audit_events
    DROP CONSTRAINT IF EXISTS audit_events_action_check;

ALTER TABLE audit_events
    ADD CONSTRAINT audit_events_action_check
    CHECK (action IN ('created', 'updated', 'disposed', 'moved'));
