-- Migration: 000011_item_relations.up.sql
-- Creates table for linking items together (e.g., computer + monitor sets)

CREATE TABLE item_relations (
    id BIGSERIAL PRIMARY KEY,
    item_id BIGINT NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    related_item_id BIGINT NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by BIGINT NOT NULL REFERENCES users(id),
    UNIQUE (item_id, related_item_id)
);

-- Index for fast lookups of related items
CREATE INDEX idx_item_relations_item_id ON item_relations(item_id);
CREATE INDEX idx_item_relations_related_item_id ON item_relations(related_item_id);

-- Prevent self-referencing and duplicate bidirectional relations
CREATE OR REPLACE FUNCTION check_item_relation() RETURNS TRIGGER AS $$
BEGIN
    -- Prevent self-referencing
    IF NEW.item_id = NEW.related_item_id THEN
        RAISE EXCEPTION 'Cannot link an item to itself';
    END IF;
    
    -- Prevent duplicate bidirectional relations (if A->B exists, prevent B->A)
    IF EXISTS (
        SELECT 1 FROM item_relations 
        WHERE item_id = NEW.related_item_id AND related_item_id = NEW.item_id
    ) THEN
        RAISE EXCEPTION 'Bidirectional relation already exists between these items';
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_check_item_relation
    BEFORE INSERT ON item_relations
    FOR EACH ROW EXECUTE FUNCTION check_item_relation();
