-- name: GetServer :one
SELECT * FROM servers WHERE id = $1;

-- name: GetServers :many
SELECT * FROM servers;

-- name: OwnServer :execresult
SELECT * FROM servers WHERE id = $1 AND owner_id = $2;

-- name: IsMember :execresult
SELECT id FROM server_members WHERE server_id = $1 AND user_id = $2;

-- name: GetServersFromUser :many
SELECT DISTINCT s.*, sm.roles, (SELECT count(id) FROM server_members smc WHERE smc.server_id=s.id) AS member_count
FROM servers s
LEFT JOIN server_members sm ON sm.server_id = s.id AND sm.user_id = $1
WHERE s.id <> 'global' OR sm.user_id IS NOT NULL;

-- name: GetServerMembers :many
SELECT u.id, u.username, u.display_name, u.avatar FROM server_members sm, users u WHERE sm.server_id = $1 AND sm.user_id = u.id;

-- name: GetMembersFromServers :many
SELECT u.id, u.username, u.display_name, u.avatar, u.banner, sm.server_id, sm.roles
FROM server_members sm, users u 
WHERE sm.server_id = ANY($1::text[]) AND sm.user_id = u.id;

-- name: GetRolesFromServers :many
SELECT r.id, r.position, r.name, r.color, r.abilities, r.server_id
FROM roles r
WHERE r.server_id = ANY($1::text[])
ORDER BY r.position;

-- name: CreateServer :one
INSERT INTO servers (
  id, owner_id, name, avatar, description, main_color, public
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: JoinServer :exec
INSERT INTO server_members (
  id, user_id, server_id
) VALUES (
  $1, $2, $3
);

-- name: LeaveServer :exec
DELETE FROM server_members WHERE user_id = $1 AND server_id = $2;

-- name: UpdateServerAvatarNBanner :exec
UPDATE servers SET avatar = $1, banner = $2, main_color = $3 WHERE id = $4 AND owner_id = $5;

-- name: UpdateServerProfile :exec
UPDATE servers SET name = $1, description = $2 WHERE id = $3 AND owner_id = $4;

-- name: DeleteServer :execresult
DELETE FROM servers WHERE id = $1 AND owner_id = $2;
