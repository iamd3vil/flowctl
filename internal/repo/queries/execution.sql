-- name: AddExecutionLog :one
WITH user_lookup AS (
    SELECT id FROM users WHERE uuid = $5
)
INSERT INTO execution_log (
    exec_id,
    parent_exec_id,
    flow_id,
    input,
    triggered_by
) VALUES (
    $1, $2, $3, $4, (SELECT id FROM user_lookup)
) RETURNING *;

-- name: UpdateExecutionStatus :one
UPDATE execution_log SET status=$1, error=$2, updated_at=$3 WHERE exec_id = $4 RETURNING *;

-- name: GetExecutionsByFlow :many
WITH user_lookup AS (
    SELECT id FROM users WHERE users.uuid = $2
), namespace_lookup AS (
    SELECT id FROM namespaces WHERE namespaces.uuid = $3
)
SELECT el.* FROM execution_log el
INNER JOIN flows f ON el.flow_id = f.id
WHERE f.id = $1 
  AND el.triggered_by = (SELECT id FROM user_lookup)
  AND f.namespace_id = (SELECT id FROM namespace_lookup);

-- name: GetExecutionByExecID :one
SELECT
    el.*,
    u.uuid AS triggered_by_uuid
FROM
    execution_log el
INNER JOIN
    users u ON el.triggered_by = u.id
WHERE
    el.exec_id = $1;

-- name: GetExecutionByExecIDWithNamespace :one
WITH namespace_lookup AS (
    SELECT id FROM namespaces WHERE namespaces.uuid = $2
)
SELECT
    el.*,
    u.uuid AS triggered_by_uuid
FROM
    execution_log el
INNER JOIN
    users u ON el.triggered_by = u.id
INNER JOIN
    flows f ON el.flow_id = f.id
WHERE
    el.exec_id = $1
    AND f.namespace_id = (SELECT id FROM namespace_lookup);

-- name: GetFlowFromExecID :one
WITH exec_log AS (
    SELECT flow_id FROM execution_log WHERE exec_id = $1
)
SELECT * FROM flows inner join exec_log on exec_log.flow_id = flows.id;

-- name: GetFlowFromExecIDWithNamespace :one
WITH exec_log AS (
    SELECT flow_id FROM execution_log WHERE exec_id = $1
), namespace_lookup AS (
    SELECT id FROM namespaces WHERE namespaces.uuid = $2
)
SELECT f.* FROM flows f
INNER JOIN exec_log el ON el.flow_id = f.id
WHERE f.namespace_id = (SELECT id FROM namespace_lookup);

-- name: GetExecutionByID :one
SELECT * FROM execution_log WHERE id = $1;

-- name: GetInputForExecByUUID :one
SELECT input FROM execution_log WHERE exec_id = $1;

-- name: GetChildrenByParentUUID :many
SELECT * FROM execution_log WHERE parent_exec_id = $1;
