-- Remove indexes
DROP INDEX IF EXISTS idx_credentials_namespace_id;
DROP INDEX IF EXISTS idx_nodes_namespace_id;
DROP INDEX IF EXISTS idx_flows_namespace_id;

-- Remove foreign key constraints
ALTER TABLE credentials DROP CONSTRAINT IF EXISTS fk_credentials_namespace_id;
ALTER TABLE nodes DROP CONSTRAINT IF EXISTS fk_nodes_namespace_id;
ALTER TABLE flows DROP CONSTRAINT IF EXISTS fk_flows_namespace_id;

-- Remove namespace_id columns
ALTER TABLE credentials DROP COLUMN IF EXISTS namespace_id;
ALTER TABLE nodes DROP COLUMN IF EXISTS namespace_id;
ALTER TABLE flows DROP COLUMN IF EXISTS namespace_id;

-- Drop group_namespace_access table
DROP TABLE IF EXISTS group_namespace_access;

-- Remove default namespace (be careful about this in production)
DELETE FROM namespaces WHERE name = 'default';