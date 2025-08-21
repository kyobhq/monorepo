-- name: GetMessage :one
SELECT * FROM messages WHERE id = $1;

-- name: GetMessagesFromChannel :many
SELECT m.*, (
  SELECT json_build_object(
    'id', u.id, 
    'avatar', u.avatar, 
    'display_name', u.display_name,
    'roles', sm.roles
  ) 
  FROM users u
  LEFT JOIN server_members sm ON (u.id = sm.user_id AND sm.server_id = $2 AND $2 != 'global')
  WHERE u.id = m.author_id
) as author
FROM messages m 
WHERE channel_id = $1
ORDER BY m.created_at
LIMIT 50;

-- name: CheckChannelMembership :execresult
SELECT c.id FROM channels c, server_members sm WHERE c.id = $1 and c.server_id = sm.server_id and sm.user_id = $2;

-- name: GetLatestMessagesSent :many
SELECT m.id, m.channel_id FROM messages m WHERE channel_id = ANY($1::text[]) ORDER BY created_at DESC LIMIT 1;

-- name: GetLatestMessagesRead :many
SELECT channel_id, last_read_message_id, unread_mention_ids FROM user_channel_read_state WHERE user_id = $1;

-- name: CreateMessage :one
INSERT INTO messages (
  id, author_id, server_id, channel_id, content, everyone, mentions_users, mentions_roles, mentions_channels, attachments
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING *;

-- name: SaveUnreadMessagesState :exec
INSERT INTO user_channel_read_state (user_id, channel_id, last_read_message_id, unread_mention_ids)
SELECT $1, unnest(@channel_ids::VARCHAR[]), unnest(@last_read_message_ids::VARCHAR[]), unnest(@unread_mention_ids::JSONB[])
ON CONFLICT (user_id, channel_id)
DO UPDATE SET 
    last_read_message_id = EXCLUDED.last_read_message_id,
    unread_mention_ids = EXCLUDED.unread_mention_ids,
    updated_at = NOW();

-- name: UpdateMessage :exec
UPDATE messages 
SET content = $1, mentions_users = $2, mentions_channels = $3, everyone = $4, updated_at = now()
WHERE id = $5;

-- name: DeleteMessage :exec
DELETE FROM messages WHERE id = $1 AND author_id = $2;

-- name: DeleteServerMessages :exec
DELETE FROM messages WHERE server_id = $1 AND author_id = $2;

-- name: GetMessageAuthor :one
SELECT author_id FROM messages WHERE id = $1;
