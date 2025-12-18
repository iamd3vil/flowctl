-- Remove completed_at column from execution_log
ALTER TABLE execution_log DROP COLUMN IF EXISTS completed_at;
