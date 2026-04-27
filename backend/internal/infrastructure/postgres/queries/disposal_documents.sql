-- queries/disposal_documents.sql
-- Disposal document queries.

-- name: CreateDisposalDocument :one
INSERT INTO disposal_documents (
    item_id, filename, url, uploaded_by
)
VALUES ($1, $2, $3, $4)
RETURNING id, uploaded_at;

-- name: GetDisposalDocumentByID :one
SELECT * FROM disposal_documents
WHERE id = $1;

-- name: ListDisposalDocumentsByItemID :many
SELECT * FROM disposal_documents
WHERE item_id = $1
ORDER BY uploaded_at ASC;

-- name: CountDisposalDocumentsByItemID :one
SELECT COUNT(*) FROM disposal_documents
WHERE item_id = $1;

-- name: DeleteDisposalDocument :exec
DELETE FROM disposal_documents
WHERE id = $1 AND item_id = $2;
