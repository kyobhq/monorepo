-- name: GetUser :one
SELECT * FROM users WHERE email = $1 OR username = $2;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

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

-- name: UpdateUserAvatarNBanner :exec
UPDATE users
  set avatar = $2, banner = $3, main_color = $4
WHERE id = $1;

-- name: UpdateUserInformations :exec
UPDATE users
set email = $2, username = $3
WHERE id = $1;

-- name: UpdateUserPassword :exec
UPDATE users
  set password = $2
WHERE id = $1;

-- name: UpdateUserProfile :exec
UPDATE users
  set display_name = $2, about_me = $3, links = $4, facts = $5
WHERE id = $1;
