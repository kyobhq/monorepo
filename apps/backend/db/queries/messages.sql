-- name: GetMessage :one
SELECT * FROM messages WHERE id = $1;

-- name: GetMessagesFromChannel :many
WITH base AS (
  SELECT m.*,
    (
      SELECT json_build_object(
        'id', u.id,
        'avatar', u.avatar,
        'display_name', u.display_name,
        'roles', sm.roles,
        'status', CASE 
            WHEN u.id = ANY($5::text[]) THEN 'online'
            ELSE 'offline'
        END
      )
      FROM users u
      LEFT JOIN server_members sm
        ON u.id = sm.user_id
      AND sm.server_id = $2
      AND $2 != 'global'
      WHERE u.id = m.author_id
    ) AS author
  FROM messages m
  WHERE m.channel_id = $1
    AND (
      ($3::text = '') OR
      m.created_at < (SELECT created_at FROM messages WHERE id = $3)
    )
    AND (
      ($4::text = '') OR
      m.created_at > (SELECT created_at FROM messages WHERE id = $4)
    )
  ORDER BY
    CASE WHEN $4::text != '' THEN m.created_at END ASC,
    CASE WHEN $4::text = ''  THEN m.created_at END DESC
  LIMIT 50
)
SELECT *
FROM base
ORDER BY created_at DESC;

-- name: CheckChannelMembership :execresult
SELECT c.id FROM channels c, server_members sm WHERE c.id = $1 and c.server_id = sm.server_id and sm.user_id = $2;

-- name: GetLatestMessagesSent :many
SELECT DISTINCT ON (channel_id) m.id, m.channel_id 
FROM messages m 
WHERE channel_id = ANY($1::text[]) 
ORDER BY channel_id, created_at DESC;

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
WITH sync_data AS (
  SELECT 
    $1 as user_id,
    unnest(@channel_ids::VARCHAR[]) as channel_id,
    unnest(@last_read_message_ids::VARCHAR[]) as last_read_message_id,
    unnest(@unread_mention_ids::JSONB[]) as unread_mention_ids
),
upsert_records AS (
  INSERT INTO user_channel_read_state (user_id, channel_id, last_read_message_id, unread_mention_ids)
  SELECT user_id, channel_id, last_read_message_id, unread_mention_ids
  FROM sync_data
  ON CONFLICT (user_id, channel_id)
  DO UPDATE SET 
    last_read_message_id = EXCLUDED.last_read_message_id,
    unread_mention_ids = EXCLUDED.unread_mention_ids,
    updated_at = NOW()
  RETURNING channel_id
),
channels_to_delete AS (
  SELECT channel_id 
  FROM user_channel_read_state 
  WHERE user_id = $1 
  AND channel_id NOT IN (SELECT channel_id FROM sync_data)
)
DELETE FROM user_channel_read_state ucrs
WHERE ucrs.user_id = $1 
AND ucrs.channel_id IN (SELECT channel_id FROM channels_to_delete);

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
