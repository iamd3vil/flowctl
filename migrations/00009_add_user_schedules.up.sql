-- Add columns to support user-created cron schedules
ALTER TABLE cron_schedules
    ADD COLUMN uuid UUID NOT NULL DEFAULT uuid_generate_v4(),
    ADD COLUMN inputs JSONB DEFAULT '{}'::jsonb,
    ADD COLUMN created_by INTEGER NOT NULL DEFAULT 1 REFERENCES users(id) ON DELETE CASCADE,
    ADD COLUMN is_user_created BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT TRUE;

-- Create indexes for performance
CREATE UNIQUE INDEX idx_cron_schedules_uuid ON cron_schedules(uuid);
CREATE INDEX idx_cron_schedules_created_by ON cron_schedules(created_by);
CREATE INDEX idx_cron_schedules_is_user_created ON cron_schedules(is_user_created);
CREATE INDEX idx_cron_schedules_is_active ON cron_schedules(is_active, flow_id);
CREATE INDEX idx_cron_schedules_user_active ON cron_schedules(created_by, is_active) WHERE is_user_created = TRUE;
