-- name: CreateCredential :one
INSERT INTO credentials (name, private_key, password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetCredentialByUUID :one
SELECT * FROM credentials WHERE uuid = $1;

-- name: GetCredentialByID :one
SELECT * FROM credentials WHERE id = $1;

-- name: ListCredentials :many
SELECT * FROM credentials;

-- name: UpdateCredential :one
UPDATE credentials
SET name = $2, private_key = $3, password = $4, updated_at = NOW()
WHERE uuid = $1
RETURNING *;

-- name: DeleteCredential :exec
DELETE FROM credentials WHERE uuid = $1;
