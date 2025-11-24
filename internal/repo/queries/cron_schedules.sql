-- name: CreateCronSchedule :one
INSERT INTO cron_schedules (flow_id, cron, timezone)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetCronSchedulesByFlowID :many
SELECT * FROM cron_schedules
WHERE flow_id = $1
ORDER BY id;

-- name: DeleteCronSchedulesByFlowID :exec
DELETE FROM cron_schedules
WHERE flow_id = $1;

-- name: GetAllCronSchedules :many
SELECT cs.*, f.slug AS flow_slug, f.name AS flow_name, n.uuid AS namespace_uuid
FROM cron_schedules cs
JOIN flows f ON cs.flow_id = f.id
JOIN namespaces n ON f.namespace_id = n.id
WHERE f.is_active = TRUE
ORDER BY cs.flow_id, cs.id;
