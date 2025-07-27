-- name: UpdateEmoji :exec
UPDATE emojis SET shortcode = $1 WHERE user_id = $2 AND id = $3;

-- name: DeleteEmoji :exec
DELETE FROM emojis WHERE user_id = $1 AND id = $2;

-- name: GetEmojis :many
SELECT id, url, shortcode FROM emojis WHERE user_id = $1;

-- name: CreateEmoji :copyfrom
INSERT INTO emojis (
  id, user_id, url, shortcode
) VALUES (
  $1, $2, $3, $4
);
