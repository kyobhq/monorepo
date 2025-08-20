-- name: GetUser :one
SELECT * FROM users WHERE email = $1 OR username = $2;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserProfile :one
SELECT u.id, u.username, u.display_name, u.avatar, u.banner, u.main_color, u.about_me, u.links, u.facts
FROM users u
WHERE u.id = $1;

-- name: GetUsersByIds :many
SELECT id, username, display_name, avatar FROM users WHERE id = ANY($1::text[]);

-- name: CreateUser :one
INSERT INTO users (
  id, email, username, display_name, avatar, banner, main_color, password
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: UpdateUserAvatarNBanner :one
UPDATE users
  set 
    avatar = COALESCE($2, avatar),
    banner = COALESCE($3, banner),
    main_color = COALESCE($4, main_color),
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdateUserEmail :one
UPDATE users
set email = $2
WHERE id = $1
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users
  set password = $2
WHERE id = $1;

-- name: UpdateUserProfile :one
UPDATE users
  set username = $2, display_name = $3, about_me = $4, links = $5, facts = $6
WHERE id = $1
RETURNING *;

-- name: GetUserLinks :many
SELECT links FROM users WHERE id = $1;

-- name: GetUserFacts :many
SELECT facts FROM users WHERE id = $1;

-- name: GetUserPassword :one
SELECT password FROM users WHERE id = $1;