-- 1. Create flow_prefixes table
CREATE TABLE flow_prefixes (
    id SERIAL PRIMARY KEY,
    uuid UUID NOT NULL DEFAULT uuid_generate_v4(),
    namespace_id INTEGER NOT NULL REFERENCES namespaces(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_prefix_per_namespace UNIQUE(namespace_id, name)
);
CREATE INDEX idx_flow_prefixes_namespace ON flow_prefixes(namespace_id);

-- 2. Add prefix_id FK to flows (replaces old VARCHAR prefix column)
ALTER TABLE flows ADD COLUMN prefix_id INTEGER REFERENCES flow_prefixes(id) ON DELETE SET NULL;
CREATE INDEX idx_flows_prefix_id ON flows(namespace_id, prefix_id);

-- 3. Create prefix_access table (uses prefix_id FK instead of VARCHAR)
CREATE TABLE prefix_access (
    id SERIAL PRIMARY KEY,
    uuid UUID NOT NULL DEFAULT uuid_generate_v4(),
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    group_id INTEGER REFERENCES groups(id) ON DELETE CASCADE,
    namespace_id INTEGER NOT NULL REFERENCES namespaces(id) ON DELETE CASCADE,
    prefix_id INTEGER NOT NULL REFERENCES flow_prefixes(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT check_single_subject CHECK (
        (user_id IS NOT NULL AND group_id IS NULL) OR
        (user_id IS NULL AND group_id IS NOT NULL)
    ),
    CONSTRAINT unique_user_prefix UNIQUE(user_id, namespace_id, prefix_id),
    CONSTRAINT unique_group_prefix UNIQUE(group_id, namespace_id, prefix_id)
);
CREATE INDEX idx_prefix_access_namespace ON prefix_access(namespace_id);
CREATE INDEX idx_prefix_access_user ON prefix_access(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX idx_prefix_access_group ON prefix_access(group_id) WHERE group_id IS NOT NULL;
