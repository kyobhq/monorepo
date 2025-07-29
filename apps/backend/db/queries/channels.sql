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

-- name: GetChannelsFromServers :many
SELECT *
FROM channels
WHERE server_id = ANY($1::text[]) AND active = true;

-- name: CreateChannel :one
INSERT INTO channels (
  id, server_id, name, description, type, e2ee, users, roles
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: UpdateChannelInformations :exec
UPDATE channels SET name = $1, description = $2 WHERE id = $3;

-- name: DeleteChannel :exec
DELETE FROM channels WHERE id = $1;

-- name: DeactivateChannel :one
UPDATE channels SET active = false
WHERE type = 'dm'
  AND array_length(users, 1) = 2
  AND $1::varchar = ANY(users) 
  AND $2::varchar = ANY(users)
RETURNING *;
