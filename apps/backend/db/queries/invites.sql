-- name: GetOrCreateInvite :one
WITH existing_invite AS (
  SELECT invite_id 
  FROM invites 
  WHERE invites.creator_id = $2 
    AND invites.server_id = $3 
    AND invites.expire_at > NOW()
  LIMIT 1
),
new_invite AS (
  INSERT INTO invites (id, creator_id, server_id, invite_id, expire_at)
  SELECT $1, $2, $3, $4, $5
  WHERE NOT EXISTS (SELECT 1 FROM existing_invite)
  RETURNING invite_id
)
SELECT COALESCE(
  (SELECT invite_id FROM existing_invite),
  (SELECT invite_id FROM new_invite)
) as invite_id;

-- name: CheckInvite :one
SELECT server_id FROM invites WHERE invite_id = $1 AND expire_at >= NOW();
