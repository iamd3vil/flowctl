-- Remove scheduled_at column from execution_log
DROP INDEX IF EXISTS idx_execution_log_scheduled_at;
ALTER TABLE execution_log DROP COLUMN IF EXISTS scheduled_at;
