-- name: AddExecutionLog :one
INSERT INTO execution_log (
    exec_id,
    flow_id,
    input,
    triggered_by
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: UpdateExecutionStatus :one
UPDATE execution_log SET status=$1, error=$2, updated_at=$3 WHERE exec_id = $4 RETURNING *;

-- name: GetExecutionsByFlow :many
SELECT * FROM execution_log WHERE flow_id = $1 and triggered_by = $2;

-- name: GetExecutionByExecID :one
SELECT * FROM execution_log WHERE exec_id = $1;

-- name: GetFlowFromExecID :one
WITH exec_log AS (
    SELECT flow_id FROM execution_log WHERE exec_id = $1
)
SELECT * FROM flows inner join exec_log on exec_log.flow_id = flows.id;