-- queries/photos.sql

-- name: AddItemPhoto :one
INSERT INTO item_photos (item_id, base64_data, mime_type)
VALUES ($1, $2, $3)
RETURNING id, item_id, base64_data, mime_type, created_at;

-- name: ListItemPhotos :many
SELECT id, item_id, base64_data, mime_type, created_at
FROM item_photos
WHERE item_id = $1
ORDER BY created_at ASC;

-- name: DeleteItemPhoto :exec
DELETE FROM item_photos
WHERE id = $1 AND item_id = $2;

-- name: DeleteAllItemPhotos :exec
DELETE FROM item_photos WHERE item_id = $1;
