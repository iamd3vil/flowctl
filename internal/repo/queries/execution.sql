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
UPDATE execution_log SET status=$1, error=$2 WHERE exec_id = $3 RETURNING *;