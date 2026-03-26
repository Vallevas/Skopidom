-- Migration: 000006_repair_audit_actions.up.sql
-- Adds 'sent_to_repair' and 'returned_from_repair' to audit_events.action.

ALTER TABLE audit_events
    DROP CONSTRAINT IF EXISTS audit_events_action_check;

ALTER TABLE audit_events
    ADD CONSTRAINT audit_events_action_check
    CHECK (action IN (
        'created',
        'updated',
        'disposed',
        'moved',
        'sent_to_repair',
        'returned_from_repair'
    ));
