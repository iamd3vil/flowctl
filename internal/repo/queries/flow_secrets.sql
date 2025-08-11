-- name: CreateFlowSecret :one
INSERT INTO flow_secrets (flow_id, key, encrypted_value, description, namespace_id)
VALUES ($1, $2, $3, $4, (SELECT id FROM namespaces WHERE namespaces.uuid = $5))
RETURNING *;

-- name: GetFlowSecretByUUID :one
SELECT fs.*, ns.uuid AS namespace_uuid FROM flow_secrets fs
JOIN namespaces ns ON fs.namespace_id = ns.id
WHERE fs.uuid = $1 AND ns.uuid = $2;

-- name: ListFlowSecrets :many
SELECT fs.*, ns.uuid AS namespace_uuid FROM flow_secrets fs
JOIN namespaces ns ON fs.namespace_id = ns.id
WHERE fs.flow_id = $1 AND ns.uuid = $2
ORDER BY fs.created_at DESC;

-- name: UpdateFlowSecret :one
UPDATE flow_secrets SET
    key = $3,
    encrypted_value = $4,
    description = $5,
    updated_at = NOW()
WHERE flow_secrets.uuid = $1 AND namespace_id = (SELECT id FROM namespaces WHERE namespaces.uuid = $2)
RETURNING *;

-- name: DeleteFlowSecret :exec
DELETE FROM flow_secrets
WHERE flow_secrets.uuid = $1 AND namespace_id = (SELECT id FROM namespaces WHERE namespaces.uuid = $2);