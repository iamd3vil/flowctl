-- Drop indexes
DROP INDEX IF EXISTS idx_cron_schedules_user_active;
DROP INDEX IF EXISTS idx_cron_schedules_is_active;
DROP INDEX IF EXISTS idx_cron_schedules_is_user_created;
DROP INDEX IF EXISTS idx_cron_schedules_created_by;
DROP INDEX IF EXISTS idx_cron_schedules_uuid;

-- Drop columns
ALTER TABLE cron_schedules
    DROP COLUMN IF EXISTS is_active,
    DROP COLUMN IF EXISTS is_user_created,
    DROP COLUMN IF EXISTS created_by,
    DROP COLUMN IF EXISTS inputs,
    DROP COLUMN IF EXISTS uuid;
