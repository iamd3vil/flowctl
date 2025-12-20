-- Add action_retries column to track retry attempts per action
ALTER TABLE execution_log ADD COLUMN action_retries JSONB DEFAULT '{}';
