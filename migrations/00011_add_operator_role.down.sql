DELETE FROM namespace_members WHERE role = 'operator';

ALTER TABLE namespace_members
    DROP CONSTRAINT IF EXISTS namespace_members_role_check,
    ADD CONSTRAINT namespace_members_role_check
        CHECK (role IN ('user', 'reviewer', 'admin'));
