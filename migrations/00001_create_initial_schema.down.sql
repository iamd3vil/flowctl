-- Drop tables in reverse order of creation (respecting foreign key dependencies)
DROP TABLE IF EXISTS group_namespace_access;
DROP TABLE IF EXISTS casbin_rule;
DROP TABLE IF EXISTS namespace_members;
DROP TABLE IF EXISTS nodes;
DROP TABLE IF EXISTS credentials;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS approvals;
DROP TABLE IF EXISTS execution_log;
DROP TABLE IF EXISTS group_memberships;
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS flows;
DROP TABLE IF EXISTS namespaces;

-- Drop views
DROP VIEW IF EXISTS user_view;
DROP VIEW IF EXISTS group_view;

-- Drop custom types
DROP TYPE IF EXISTS authentication_method;
DROP TYPE IF EXISTS approval_status;
DROP TYPE IF EXISTS execution_status;
DROP TYPE IF EXISTS user_role_type;
DROP TYPE IF EXISTS user_login_type;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";
