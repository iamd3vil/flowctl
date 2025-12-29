-- Add scheduled_at column to execution_log for scheduled flow executions
ALTER TABLE execution_log ADD COLUMN scheduled_at TIMESTAMP WITH TIME ZONE DEFAULT NULL;

-- Create partial index for scheduled executions (only index non-null values)
CREATE INDEX idx_execution_log_scheduled_at ON execution_log(scheduled_at) WHERE scheduled_at IS NOT NULL;
