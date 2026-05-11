-- queries/photos.sql

-- name: AddItemPhoto :one
INSERT INTO item_photos (item_id, url)
VALUES ($1, $2)
RETURNING id, item_id, url, created_at;

-- name: ListItemPhotos :many
SELECT id, item_id, url, created_at
FROM item_photos
WHERE item_id = $1
ORDER BY created_at ASC;

-- name: DeleteItemPhoto :exec
DELETE FROM item_photos
WHERE id = $1 AND item_id = $2;

-- name: DeleteAllItemPhotos :exec
DELETE FROM item_photos WHERE item_id = $1;

-- name: AddItemPhotoData :one
INSERT INTO item_photos_data (photo_id, data, mime_type)
VALUES ($1, $2, $3)
RETURNING id, photo_id, data, mime_type, created_at;

-- name: GetItemPhotoData :one
SELECT id, photo_id, data, mime_type, created_at
FROM item_photos_data
WHERE photo_id = $1;

-- name: UpdateItemPhotoData :one
UPDATE item_photos_data
SET data = $2, mime_type = $3
WHERE photo_id = $1
RETURNING id, photo_id, data, mime_type, created_at;

-- name: DeleteItemPhotoData :exec
DELETE FROM item_photos_data
WHERE photo_id = $1;
