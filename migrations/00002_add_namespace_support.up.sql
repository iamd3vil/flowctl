-- Add namespace foreign key to flows table
ALTER TABLE flows ADD COLUMN namespace_id INTEGER;

-- Add namespace foreign key to nodes table  
ALTER TABLE nodes ADD COLUMN namespace_id INTEGER;

-- Add namespace foreign key to credentials table
ALTER TABLE credentials ADD COLUMN namespace_id INTEGER;

-- Create group_namespace_access table for access control
CREATE TABLE IF NOT EXISTS group_namespace_access (
    id SERIAL PRIMARY KEY,
    group_id INTEGER NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    namespace_id INTEGER NOT NULL REFERENCES namespaces(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(group_id, namespace_id)
);
CREATE INDEX idx_group_namespace_access_group_id ON group_namespace_access(group_id);
CREATE INDEX idx_group_namespace_access_namespace_id ON group_namespace_access(namespace_id);

-- Create default namespace
INSERT INTO namespaces (name, created_at, updated_at) 
VALUES ('default', NOW(), NOW())
ON CONFLICT (name) DO NOTHING;

-- Update all existing resources to use default namespace
UPDATE flows SET namespace_id = (SELECT id FROM namespaces WHERE name = 'default')
WHERE namespace_id IS NULL;

UPDATE nodes SET namespace_id = (SELECT id FROM namespaces WHERE name = 'default')
WHERE namespace_id IS NULL;

UPDATE credentials SET namespace_id = (SELECT id FROM namespaces WHERE name = 'default')
WHERE namespace_id IS NULL;

-- Make namespace_id required
ALTER TABLE flows ALTER COLUMN namespace_id SET NOT NULL;
ALTER TABLE nodes ALTER COLUMN namespace_id SET NOT NULL;
ALTER TABLE credentials ALTER COLUMN namespace_id SET NOT NULL;

-- Add foreign key constraints
ALTER TABLE flows ADD CONSTRAINT fk_flows_namespace_id FOREIGN KEY (namespace_id) REFERENCES namespaces(id) ON DELETE CASCADE;
ALTER TABLE nodes ADD CONSTRAINT fk_nodes_namespace_id FOREIGN KEY (namespace_id) REFERENCES namespaces(id) ON DELETE CASCADE;
ALTER TABLE credentials ADD CONSTRAINT fk_credentials_namespace_id FOREIGN KEY (namespace_id) REFERENCES namespaces(id) ON DELETE CASCADE;

-- Add indexes for better performance
CREATE INDEX idx_flows_namespace_id ON flows(namespace_id);
CREATE INDEX idx_nodes_namespace_id ON nodes(namespace_id);
CREATE INDEX idx_credentials_namespace_id ON credentials(namespace_id);