-- Rollback: Add back the cron_schedules column
ALTER TABLE flows ADD COLUMN IF NOT EXISTS cron_schedules TEXT[];
