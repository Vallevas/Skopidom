-- Migration: 001_initial_schema.up.sql
-- Creates the full initial schema for the inventory system.

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- ── Users ────────────────────────────────────────────────────────────────────

CREATE TABLE users (
    id            BIGSERIAL    PRIMARY KEY,
    full_name     VARCHAR(255) NOT NULL,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role          VARCHAR(50)  NOT NULL DEFAULT 'editor'
                               CHECK (role IN ('admin', 'editor')),
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- ── Buildings ────────────────────────────────────────────────────────────────

CREATE TABLE buildings (
    id      BIGSERIAL    PRIMARY KEY,
    name    VARCHAR(255) NOT NULL UNIQUE,
    address VARCHAR(512) NOT NULL DEFAULT ''
);

-- ── Categories ───────────────────────────────────────────────────────────────

CREATE TABLE categories (
    id   BIGSERIAL    PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

-- ── Rooms ────────────────────────────────────────────────────────────────────

CREATE TABLE rooms (
    id          BIGSERIAL    PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    building_id BIGINT       NOT NULL REFERENCES buildings(id),
    UNIQUE (name, building_id)
);

-- ── Items ────────────────────────────────────────────────────────────────────

CREATE TABLE items (
    id             BIGSERIAL    PRIMARY KEY,
    barcode        VARCHAR(255) NOT NULL UNIQUE,
    name           VARCHAR(255) NOT NULL,
    category_id    BIGINT       NOT NULL REFERENCES categories(id),
    room_id        BIGINT       NOT NULL REFERENCES rooms(id),
    description    TEXT         NOT NULL DEFAULT '',
    photo_url      VARCHAR(512) NOT NULL DEFAULT '',
    status         VARCHAR(50)  NOT NULL DEFAULT 'active'
                                CHECK (status IN ('active', 'disposed')),
    -- tx_hash is populated once blockchain integration is added.
    tx_hash        VARCHAR(66)  NOT NULL DEFAULT '',
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    created_by     BIGINT       NOT NULL REFERENCES users(id),
    last_edited_by BIGINT       NOT NULL REFERENCES users(id)
);

-- Indexes for the most common list-query filters.
CREATE INDEX idx_items_category_id ON items(category_id);
CREATE INDEX idx_items_room_id     ON items(room_id);
CREATE INDEX idx_items_status      ON items(status);
CREATE INDEX idx_items_created_at  ON items(created_at);
