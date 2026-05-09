-- Migration: 000011_item_relations.down.sql
-- Drops item_relations table and related trigger/function

DROP TRIGGER IF EXISTS trg_check_item_relation ON item_relations;
DROP FUNCTION IF EXISTS check_item_relation();
DROP TABLE IF EXISTS item_relations;
