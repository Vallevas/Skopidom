-- Migration: 001_initial_schema.down.sql
-- Drops all tables created by the initial schema migration.

DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS rooms;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS buildings;
DROP TABLE IF EXISTS users;
