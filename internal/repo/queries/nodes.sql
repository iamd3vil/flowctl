-- name: CreateNode :one
INSERT INTO nodes (name, hostname, port, username, os_family, tags, auth_method, credential_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetNodeByUUID :one
SELECT * FROM nodes WHERE uuid = $1;

-- name: ListNodes :many
WITH filtered AS (
    SELECT * FROM nodes
),
total AS (
    SELECT COUNT(*) AS total_count FROM filtered
),
paged AS (
    SELECT * FROM filtered
    LIMIT $1 OFFSET $2
),
page_count AS (
    SELECT COUNT(*) AS page_count FROM paged
)
SELECT 
    p.*,
    pc.page_count,
    t.total_count
FROM paged p, page_count pc, total t;

-- name: UpdateNode :one
UPDATE nodes
SET name = $2, hostname = $3, port = $4, username = $5, os_family = $6, tags = $7, auth_method = $8, credential_id = $9, updated_at = NOW()
WHERE uuid = $1
RETURNING *;

-- name: DeleteNode :exec
DELETE FROM nodes WHERE uuid = $1;
