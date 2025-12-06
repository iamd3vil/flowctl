-- name: CreateNamespaceSecret :one
INSERT INTO namespace_secrets (key, encrypted_value, description, namespace_id)
VALUES ($1, $2, $3, (SELECT id FROM namespaces WHERE namespaces.uuid = $4))
RETURNING *;

-- name: GetNamespaceSecretByUUID :one
SELECT ns.*, n.uuid AS namespace_uuid FROM namespace_secrets ns
JOIN namespaces n ON ns.namespace_id = n.id
WHERE ns.uuid = $1 AND n.uuid = $2;

-- name: ListNamespaceSecrets :many
SELECT ns.*, n.uuid AS namespace_uuid FROM namespace_secrets ns
JOIN namespaces n ON ns.namespace_id = n.id
WHERE n.uuid = $1
ORDER BY ns.created_at DESC;

-- name: UpdateNamespaceSecret :one
UPDATE namespace_secrets SET
    encrypted_value = $3,
    description = $4,
    updated_at = NOW()
WHERE namespace_secrets.uuid = $1 AND namespace_id = (SELECT id FROM namespaces WHERE namespaces.uuid = $2)
RETURNING *;

-- name: DeleteNamespaceSecret :exec
DELETE FROM namespace_secrets
WHERE namespace_secrets.uuid = $1 AND namespace_id = (SELECT id FROM namespaces WHERE namespaces.uuid = $2);

-- name: GetDecryptedNamespaceSecrets :many
-- Used internally for execution - returns all secrets for a namespace
SELECT ns.key, ns.encrypted_value FROM namespace_secrets ns
JOIN namespaces n ON ns.namespace_id = n.id
WHERE n.uuid = $1;
