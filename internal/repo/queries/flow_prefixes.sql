-- name: CreateFlowPrefix :one
INSERT INTO flow_prefixes (namespace_id, name, description)
VALUES ((SELECT namespaces.id FROM namespaces WHERE namespaces.uuid = $1), $2, $3)
RETURNING *;

-- name: UpdateFlowPrefix :one
UPDATE flow_prefixes SET
    name = $3,
    description = $4,
    updated_at = NOW()
WHERE flow_prefixes.uuid = $1
AND flow_prefixes.namespace_id = (SELECT namespaces.id FROM namespaces WHERE namespaces.uuid = $2)
RETURNING *;

-- name: DeleteFlowPrefix :exec
DELETE FROM flow_prefixes
WHERE flow_prefixes.uuid = $1
AND flow_prefixes.namespace_id = (SELECT namespaces.id FROM namespaces WHERE namespaces.uuid = $2);

-- name: GetFlowPrefixByUUID :one
SELECT fp.* FROM flow_prefixes fp
WHERE fp.uuid = $1
AND fp.namespace_id = (SELECT namespaces.id FROM namespaces WHERE namespaces.uuid = $2);

-- name: GetFlowPrefixByName :one
SELECT fp.* FROM flow_prefixes fp
WHERE fp.name = $1
AND fp.namespace_id = (SELECT namespaces.id FROM namespaces WHERE namespaces.uuid = $2);

-- name: ListFlowPrefixes :many
SELECT fp.* FROM flow_prefixes fp
WHERE fp.namespace_id = (SELECT namespaces.id FROM namespaces WHERE namespaces.uuid = $1)
ORDER BY fp.name ASC;

-- name: GetFlowsByPrefixUUID :many
SELECT f.*, n.uuid AS namespace_uuid FROM flows f
JOIN namespaces n ON f.namespace_id = n.id
JOIN flow_prefixes fp ON f.prefix_id = fp.id
WHERE fp.uuid = $1
AND n.uuid = $2
AND f.is_active = TRUE
ORDER BY f.name ASC;
