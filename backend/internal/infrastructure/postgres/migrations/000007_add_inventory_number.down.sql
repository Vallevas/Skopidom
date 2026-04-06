-- Migration: 000007_add_inventory_number.down.sql
-- Removes inventory_number field from items table

DROP INDEX IF EXISTS idx_items_inventory_number;

ALTER TABLE items
DROP CONSTRAINT IF EXISTS items_inventory_number_unique;

ALTER TABLE items
DROP COLUMN IF EXISTS inventory_number;
