-- Migration: 000003_audit_events.up.sql
-- Immutable audit log for all item lifecycle events.
-- tx_hash is populated once blockchain integration is added.

CREATE TABLE audit_events (
    id         BIGSERIAL    PRIMARY KEY,
    item_id    BIGINT       NOT NULL REFERENCES items(id),
    actor_id   BIGINT       NOT NULL REFERENCES users(id),
    action     VARCHAR(50)  NOT NULL
                            CHECK (action IN ('created', 'updated', 'disposed')),
    -- payload holds a JSON snapshot of the item state at the time of the event.
    payload    JSONB        NOT NULL DEFAULT '{}',
    -- tx_hash is empty until blockchain integration is enabled.
    tx_hash    VARCHAR(66)  NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Partial index — most queries filter by item_id.
CREATE INDEX idx_audit_events_item_id ON audit_events(item_id);
