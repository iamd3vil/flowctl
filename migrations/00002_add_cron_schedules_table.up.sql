-- Create cron_schedules table to store flow schedules with timezone support
CREATE TABLE cron_schedules (
    id SERIAL PRIMARY KEY,
    flow_id INTEGER NOT NULL REFERENCES flows(id) ON DELETE CASCADE,
    cron VARCHAR(255) NOT NULL,
    timezone VARCHAR(100) NOT NULL DEFAULT 'UTC',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create index for efficient flow schedule lookups
CREATE INDEX idx_cron_schedules_flow_id ON cron_schedules(flow_id);

-- Migrate existing data from flows.cron_schedules array to new table
-- Uses unnest to expand array into rows, defaulting timezone to UTC
INSERT INTO cron_schedules (flow_id, cron, timezone)
SELECT id, unnest(cron_schedules), 'UTC'
FROM flows
WHERE cron_schedules IS NOT NULL AND array_length(cron_schedules, 1) > 0;
