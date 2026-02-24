-- name: AssignUserPrefixAccess :exec
INSERT INTO prefix_access (user_id, namespace_id, prefix_id)
VALUES ((SELECT users.id FROM users WHERE users.uuid = $1),
        (SELECT namespaces.id FROM namespaces WHERE namespaces.uuid = $2),
        (SELECT flow_prefixes.id FROM flow_prefixes WHERE flow_prefixes.name = $3 AND flow_prefixes.namespace_id = (SELECT namespaces.id FROM namespaces WHERE namespaces.uuid = $2)))
ON CONFLICT ON CONSTRAINT unique_user_prefix DO NOTHING;

-- name: AssignGroupPrefixAccess :exec
INSERT INTO prefix_access (group_id, namespace_id, prefix_id)
VALUES ((SELECT groups.id FROM groups WHERE groups.uuid = $1),
        (SELECT namespaces.id FROM namespaces WHERE namespaces.uuid = $2),
        (SELECT flow_prefixes.id FROM flow_prefixes WHERE flow_prefixes.name = $3 AND flow_prefixes.namespace_id = (SELECT namespaces.id FROM namespaces WHERE namespaces.uuid = $2)))
ON CONFLICT ON CONSTRAINT unique_group_prefix DO NOTHING;

-- name: RevokeUserPrefixAccess :exec
DELETE FROM prefix_access
WHERE user_id = (SELECT users.id FROM users WHERE users.uuid = $1)
AND namespace_id = (SELECT namespaces.id FROM namespaces WHERE namespaces.uuid = $2)
AND prefix_id = (SELECT flow_prefixes.id FROM flow_prefixes WHERE flow_prefixes.name = $3 AND flow_prefixes.namespace_id = (SELECT namespaces.id FROM namespaces WHERE namespaces.uuid = $2));

-- name: RevokeGroupPrefixAccess :exec
DELETE FROM prefix_access
WHERE group_id = (SELECT groups.id FROM groups WHERE groups.uuid = $1)
AND namespace_id = (SELECT namespaces.id FROM namespaces WHERE namespaces.uuid = $2)
AND prefix_id = (SELECT flow_prefixes.id FROM flow_prefixes WHERE flow_prefixes.name = $3 AND flow_prefixes.namespace_id = (SELECT namespaces.id FROM namespaces WHERE namespaces.uuid = $2));

-- name: GetAllPrefixAccesses :many
SELECT pa.uuid,
       COALESCE(u.uuid, g.uuid) as subject_uuid,
       CASE WHEN pa.user_id IS NOT NULL THEN 'user' ELSE 'group' END as subject_type,
       n.uuid as namespace_uuid,
       fp.name as prefix,
       pa.created_at
FROM prefix_access pa
JOIN namespaces n ON pa.namespace_id = n.id
JOIN flow_prefixes fp ON pa.prefix_id = fp.id
LEFT JOIN users u ON pa.user_id = u.id
LEFT JOIN groups g ON pa.group_id = g.id;

-- name: GetPrefixMembers :many
SELECT pa.uuid, COALESCE(u.uuid, g.uuid) as subject_uuid,
       COALESCE(u.name, g.name) as subject_name,
       CASE WHEN pa.user_id IS NOT NULL THEN 'user' ELSE 'group' END as subject_type,
       fp.name as prefix, pa.created_at
FROM prefix_access pa
JOIN flow_prefixes fp ON pa.prefix_id = fp.id
LEFT JOIN users u ON pa.user_id = u.id
LEFT JOIN groups g ON pa.group_id = g.id
WHERE pa.namespace_id = (SELECT namespaces.id FROM namespaces WHERE namespaces.uuid = $1)
AND fp.name = $2;

-- name: GetMemberPrefixes :many
SELECT fp.name as prefix, pa.created_at FROM prefix_access pa
JOIN flow_prefixes fp ON pa.prefix_id = fp.id
WHERE pa.namespace_id = (SELECT namespaces.id FROM namespaces WHERE namespaces.uuid = $1)
AND (
    pa.user_id = (SELECT nm.user_id FROM namespace_members nm WHERE nm.uuid = $2)
    OR pa.group_id = (SELECT nm.group_id FROM namespace_members nm WHERE nm.uuid = $2)
);

-- name: RevokeAllMemberPrefixAccess :exec
DELETE FROM prefix_access
WHERE namespace_id = (SELECT namespaces.id FROM namespaces WHERE namespaces.uuid = $1)
AND (
    user_id = (SELECT nm.user_id FROM namespace_members nm WHERE nm.uuid = $2)
    OR group_id = (SELECT nm.group_id FROM namespace_members nm WHERE nm.uuid = $2)
);

-- name: GetUserAccessiblePrefixes :many
SELECT DISTINCT fp.name as prefix FROM prefix_access pa
JOIN flow_prefixes fp ON pa.prefix_id = fp.id
WHERE pa.namespace_id = (SELECT namespaces.id FROM namespaces WHERE namespaces.uuid = $1)
AND (
    pa.user_id = (SELECT users.id FROM users WHERE users.uuid = $2)
    OR pa.group_id IN (SELECT gm.group_id FROM group_memberships gm JOIN users u ON gm.user_id = u.id WHERE u.uuid = $2)
);
