-- name: GetChannel :one
SELECT * FROM channels WHERE id = $1;

-- name: GetFriendChannels :many
SELECT *
FROM channels
WHERE server_id = 'global' AND $1::text = ANY(users);

-- name: GetChannelsFromServer :many
SELECT *
FROM channels
WHERE server_id = $1;

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
WHERE server_id = ANY($1::text[]);

-- name: CreateChannel :one
INSERT INTO channels (
  id, position, server_id, category_id, friendship_id, name, description, type, e2ee, users, roles
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
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

-- name: UpdateCategoryInformations :exec
UPDATE channel_categories SET name = $2, users = $3, roles = $4 WHERE id = $1;

-- name: DeleteChannel :exec
DELETE FROM channels WHERE id = $1;

-- name: DeleteCategory :exec
DELETE FROM channel_categories WHERE id = $1;

-- -- name: DeactivateChannel :exec
-- UPDATE channels SET active = false
-- WHERE friendship_id = $1;

-- name: GetChannelsIDs :many
SELECT id, server_id, users FROM channels WHERE id <> 'global';
