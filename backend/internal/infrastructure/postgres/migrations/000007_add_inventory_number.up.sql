-- Migration: 000007_add_inventory_number.up.sql
-- Adds inventory_number field to items table

ALTER TABLE items
ADD COLUMN inventory_number VARCHAR(255) NOT NULL DEFAULT '';

-- Add unique constraint
ALTER TABLE items
ADD CONSTRAINT items_inventory_number_unique UNIQUE (inventory_number);

-- Add index for fast lookups
CREATE INDEX idx_items_inventory_number ON items(inventory_number);
