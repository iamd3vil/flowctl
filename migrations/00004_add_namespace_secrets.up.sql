CREATE TABLE IF NOT EXISTS namespace_secrets (
    id SERIAL PRIMARY KEY,
    uuid UUID NOT NULL DEFAULT uuid_generate_v4(),
    key VARCHAR(255) NOT NULL,
    encrypted_value TEXT NOT NULL,
    description TEXT,
    namespace_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (namespace_id) REFERENCES namespaces(id) ON DELETE CASCADE,
    UNIQUE(key, namespace_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_namespace_secrets_uuid ON namespace_secrets(uuid);
CREATE INDEX IF NOT EXISTS idx_namespace_secrets_namespace_id ON namespace_secrets(namespace_id);
