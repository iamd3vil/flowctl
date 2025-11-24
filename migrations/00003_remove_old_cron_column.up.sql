-- Remove old cron_schedules column from flows table
-- This column is now replaced by the cron_schedules table
ALTER TABLE flows DROP COLUMN IF EXISTS cron_schedules;
