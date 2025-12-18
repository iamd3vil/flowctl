-- Add completed_at column to track when executions actually complete
ALTER TABLE execution_log ADD COLUMN completed_at TIMESTAMP WITH TIME ZONE;

UPDATE execution_log
SET completed_at = updated_at
WHERE status IN ('completed', 'errored', 'cancelled');
