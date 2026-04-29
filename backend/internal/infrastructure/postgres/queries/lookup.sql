-- queries/lookup.sql
-- Reference data queries: categories, buildings, rooms.

-- ── Categories ───────────────────────────────────────────────────────────────

-- name: CreateCategory :one
INSERT INTO categories (name) VALUES ($1) RETURNING id;

-- name: GetCategoryByID :one
SELECT id, name FROM categories WHERE id = $1;

-- name: GetCategoryByName :one
SELECT id, name FROM categories WHERE name = $1;

-- name: ListCategories :many
SELECT id, name FROM categories ORDER BY name ASC;

-- name: UpdateCategory :exec
UPDATE categories SET name = $1 WHERE id = $2;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1;

-- name: CountItemsByCategory :one
SELECT COUNT(*) FROM items WHERE category_id = $1;

-- ── Buildings ─────────────────────────────────────────────────────────────────

-- name: CreateBuilding :one
INSERT INTO buildings (name, address) VALUES ($1, $2) RETURNING id;

-- name: GetBuildingByID :one
SELECT id, name, address FROM buildings WHERE id = $1;

-- name: GetBuildingByName :one
SELECT id, name, address FROM buildings WHERE name = $1;

-- name: ListBuildings :many
SELECT id, name, address FROM buildings ORDER BY name ASC;

-- name: UpdateBuilding :exec
UPDATE buildings SET name = $1, address = $2 WHERE id = $3;

-- name: DeleteBuilding :exec
DELETE FROM buildings WHERE id = $1;

-- name: CountRoomsByBuilding :one
SELECT COUNT(*) FROM rooms WHERE building_id = $1;

-- ── Rooms ─────────────────────────────────────────────────────────────────────

-- name: CreateRoom :one
INSERT INTO rooms (name, building_id) VALUES ($1, $2) RETURNING id;

-- name: GetRoomByID :one
SELECT r.id, r.name, r.building_id, b.name AS building_name, b.address AS building_address
FROM rooms r
JOIN buildings b ON b.id = r.building_id
WHERE r.id = $1;

-- name: GetRoomByNameAndBuilding :one
SELECT r.id, r.name, r.building_id, b.name AS building_name, b.address AS building_address
FROM rooms r
JOIN buildings b ON b.id = r.building_id
WHERE r.name = $1 AND r.building_id = $2;

-- name: ListRooms :many
SELECT r.id, r.name, r.building_id, b.name AS building_name, b.address AS building_address
FROM rooms r
JOIN buildings b ON b.id = r.building_id
ORDER BY b.name ASC, r.name ASC;

-- name: ListRoomsByBuilding :many
SELECT r.id, r.name, r.building_id, b.name AS building_name, b.address AS building_address
FROM rooms r
JOIN buildings b ON b.id = r.building_id
WHERE r.building_id = $1
ORDER BY r.name ASC;

-- name: UpdateRoom :exec
UPDATE rooms SET name = $1, building_id = $2 WHERE id = $3;

-- name: DeleteRoom :exec
DELETE FROM rooms WHERE id = $1;

-- name: CountItemsByRoom :one
SELECT COUNT(*) FROM items WHERE room_id = $1;
