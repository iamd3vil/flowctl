-- name: GetUserByUUID :one
SELECT * FROM users WHERE uuid = $1;

-- name: DeleteUserByUUID :exec
DELETE FROM users WHERE uuid = $1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: GetUserByUsernameWithGroups :one
SELECT * FROM user_view WHERE username = $1;

-- name: GetAllUsersWithGroups :many
SELECT * FROM user_view;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (
    username,
    password,
    login_type,
    role,
    name
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: SearchUsersWithGroups :many
SELECT * FROM user_view WHERE lower(name) LIKE '%' || lower($1::text) || '%' OR lower(username) LIKE '%' || lower($1::text) || '%';
