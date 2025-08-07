-- name: GetChannel :one
SELECT * FROM channels WHERE id = $1;

-- name: GetFriendChannels :many
SELECT *
FROM channels
WHERE server_id = 'global' AND $1::text = ANY(users) AND active = true;

-- name: GetChannelsFromServer :many
SELECT *
FROM channels
WHERE server_id = $1 AND active = true;

-- name: GetCategoriesFromServer :many
SELECT *
FROM channel_categories
WHERE server_id = $1;

-- name: GetCategoriesFromServers :many
SELECT *
FROM channel_categories
WHERE server_id = ANY($1::text[]);

-- name: GetChannelsFromServers :many
SELECT *
FROM channels
WHERE server_id = ANY($1::text[]) AND active = true;

-- name: CreateChannel :one
INSERT INTO channels (
  id, position, server_id, category_id, name, description, type, e2ee, users, roles
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING *;

-- name: CreateCategory :one
INSERT INTO channel_categories (
  id, position, server_id, name, users, roles, e2ee
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: PinChannel :exec
INSERT INTO channel_pins (
  id, position, server_id, channel_id, user_id
) VALUES (
  $1, $2, $3, $4, $5
);

-- name: UpdateChannelInformations :exec
UPDATE channels SET name = $2, description = $3, users = $4, roles = $5 WHERE id = $1;

-- name: DeleteChannel :exec
DELETE FROM channels WHERE id = $1;

-- name: DeleteCategory :exec
DELETE FROM channel_categories WHERE id = $1;

-- name: DeactivateChannel :one
UPDATE channels SET active = false
WHERE type = 'dm'
  AND array_length(users, 1) = 2
  AND $1::varchar = ANY(users) 
  AND $2::varchar = ANY(users)
RETURNING *;

-- name: GetChannelsIDs :many
SELECT id, server_id FROM channels WHERE id <> 'global';
