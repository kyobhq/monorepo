-- name: AddFriend :one
INSERT INTO friends (
  id, sender_id, receiver_id
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: AcceptFriend :exec
UPDATE friends SET accepted=true WHERE id=$1 AND receiver_id = $2;

-- name: DeleteFriend :exec
DELETE FROM friends WHERE id=$1 AND receiver_id = $2 OR sender_id = $2;

-- name: GetFriends :many
SELECT u.id, u.display_name, u.avatar, u.banner, u.about_me, f.accepted, f.id AS friendship_id, 
       f.sender_id AS friendship_sender_id, c.id AS channel_id, 'offline' as status
FROM users u
INNER JOIN friends f ON u.id = f.receiver_id
LEFT JOIN channels c ON $1::text = ANY(c.users) AND u.id::text = ANY(c.users)
WHERE f.sender_id = $1

UNION

SELECT u.id, u.display_name, u.avatar, u.banner, u.about_me, f.accepted, f.id AS friendship_id, 
       f.sender_id AS friendship_sender_id, c.id AS channel_id, 'offline' as status
FROM users u
INNER JOIN friends f ON u.id = f.sender_id  
LEFT JOIN channels c ON $1 = ANY(c.users) AND u.id::text = ANY(c.users)
WHERE f.receiver_id = $1;

-- name: GetFriendIDs :many
SELECT u.id
FROM users u
INNER JOIN friends f ON u.id = f.receiver_id
WHERE f.sender_id = $1

UNION

SELECT u.id
FROM users u
INNER JOIN friends f ON u.id = f.sender_id
WHERE f.receiver_id = $1;

-- name: GetExistingChannel :one
-- UPDATE channels SET active = true
-- WHERE friendship_id = $1
-- RETURNING *;
