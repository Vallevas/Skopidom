-- Migration: 0000011_item_relations.up.sql
-- Creates table for symmetric item-to-item relations (kits, bundles, etc.)

CREATE TABLE item_relations (
    id BIGSERIAL PRIMARY KEY,
    item_id_1 BIGINT NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    item_id_2 BIGINT NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by BIGINT NOT NULL REFERENCES users(id),
    
    -- Ensure item_id_1 < item_id_2 for symmetry (no duplicate pairs in different order)
    CONSTRAINT chk_item_order CHECK (item_id_1 < item_id_2),
    
    -- Unique constraint on the pair (order-independent)
    CONSTRAINT uq_item_pair UNIQUE (item_id_1, item_id_2)
);

-- Index for fast lookup by either item
CREATE INDEX idx_item_relations_item_id_1 ON item_relations(item_id_1);
CREATE INDEX idx_item_relations_item_id_2 ON item_relations(item_id_2);

-- Composite index for efficient queries
CREATE INDEX idx_item_relations_pair ON item_relations(item_id_1, item_id_2);
