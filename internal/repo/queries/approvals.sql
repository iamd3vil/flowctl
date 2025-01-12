-- name: AddApprovalRequest :one
WITH inserted_approval AS (
    INSERT INTO approvals (
        exec_log_id,
        approvers,
        action_id
    ) VALUES (
        $1, $2, $3
    ) RETURNING *
)
SELECT
    a.*,
    u.name as requested_by
FROM inserted_approval a
JOIN execution_log el ON a.exec_log_id = el.id
JOIN users u ON el.triggered_by = u.id;

-- name: ApproveRequestByUUID :one
WITH updated AS (
    UPDATE approvals SET status = 'approved', decided_by = $2, updated_at = NOW()
    WHERE approvals.uuid = $1
    RETURNING *
)
SELECT
    a.*,
    u.name as requested_by
FROM updated a
JOIN execution_log el ON a.exec_log_id = el.id
JOIN users u ON el.triggered_by = u.id;

-- name: RejectRequestByUUID :one
WITH updated AS (
    UPDATE approvals SET status = 'rejected', decided_by = $2, updated_at = NOW()
    WHERE approvals.uuid = $1
    RETURNING *
)
SELECT
    a.*,
    u.name as requested_by
FROM updated a
JOIN execution_log el ON a.exec_log_id = el.id
JOIN users u ON el.triggered_by = u.id;

-- name: UpdateApprovalStatusByUUID :one
WITH updated AS (
    UPDATE approvals SET status = $1, decided_by = $2, updated_at = NOW()
    WHERE uuid = $1
    RETURNING *
)
SELECT
    a.*,
    u.name as requested_by
FROM updated a
JOIN execution_log el ON a.exec_log_id = el.id
JOIN users u ON el.triggered_by = u.id;

-- name: GetApprovalByUUID :one
SELECT
    a.*,
    u.name as requested_by
FROM approvals a
JOIN execution_log el ON a.exec_log_id = el.id
JOIN users u ON el.triggered_by = u.id
WHERE a.uuid = $1;

-- name: GetApprovalRequestForActionAndExec :one
WITH exec_lookup AS (
    SELECT id FROM execution_log WHERE exec_id = $1
)
SELECT * FROM approvals WHERE exec_log_id = (SELECT id FROM exec_lookup) AND action_id = $2;

-- name: GetPendingApprovalRequestForExec :one
WITH exec_lookup AS (
    SELECT id FROM execution_log WHERE execution_log.exec_id = $1
)
SELECT
    a.*,
    u.name as requested_by
FROM approvals a
JOIN execution_log el ON a.exec_log_id = el.id
JOIN users u ON el.triggered_by = u.id
WHERE a.exec_log_id = (SELECT id FROM exec_lookup) AND a.status = 'pending';
