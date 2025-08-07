-- name: GetServer :one
SELECT * FROM servers WHERE id = $1;

-- name: GetServers :many
SELECT * FROM servers;

-- name: OwnServer :execresult
SELECT * FROM servers WHERE id = $1 AND owner_id = $2;

-- name: IsMember :execresult
SELECT id FROM server_members WHERE server_id = $1 AND user_id = $2;

-- name: GetServersFromUser :many
SELECT DISTINCT s.*, sm.roles, sm.position, (SELECT count(id) FROM server_members smc WHERE smc.server_id=s.id) AS member_count
FROM servers s
INNER JOIN server_members sm ON sm.server_id = s.id AND sm.user_id = $1
WHERE s.id <> 'global';

-- name: GetServerMembers :many
SELECT 
    u.id, 
    u.username, 
    u.display_name, 
    u.avatar,
    sm.roles,
    COALESCE(MIN(r.position), 999999) as min_role_position
FROM server_members sm
JOIN users u ON u.id = sm.user_id
LEFT JOIN roles r ON r.server_id = sm.server_id AND r.id = ANY(sm.roles)
WHERE sm.server_id = $1
GROUP BY sm.user_id, u.id, u.username, u.display_name, u.avatar, sm.roles
ORDER BY min_role_position, u.username
LIMIT $2 OFFSET $3;

-- name: GetServerInformations :one
SELECT 
    (
        SELECT json_agg(json_build_object(
            'id', ranked_members.user_id,
            'username', ranked_members.username,
            'display_name', ranked_members.display_name,
            'avatar', ranked_members.avatar,
            'roles', ranked_members.roles,
            'status', CASE 
                WHEN ranked_members.user_id = ANY($2::text[]) THEN 'online'
                ELSE 'offline'
            END
        ))
        FROM (
            SELECT 
                u.id as user_id,
                u.username,
                u.display_name,
                u.avatar,
                sm.roles,
                COALESCE(MIN(r.position), 999999) as min_role_position
            FROM server_members sm
            JOIN users u ON u.id = sm.user_id
            LEFT JOIN roles r ON r.server_id = sm.server_id AND r.id = ANY(sm.roles)
            WHERE sm.server_id = s.id
            GROUP BY sm.user_id, u.id, u.username, u.display_name, u.avatar, sm.roles
            ORDER BY min_role_position, u.username
            LIMIT 50
        ) ranked_members
    ) as members,
    (
        SELECT json_agg(json_build_object(
            'id', r.id,
            'name', r.name,
            'color', r.color,
            'position', r.position,
            'abilities', r.abilities
        ) ORDER BY r.position)
        FROM roles r
        WHERE r.server_id = s.id
    ) as roles,
    (
        SELECT count(id)
        FROM server_members smc
        WHERE smc.server_id = s.id
    ) as member_count
FROM servers s
WHERE s.id = $1;

-- name: GetMembersFromServers :many
SELECT u.id, u.username, u.display_name, u.avatar, u.banner, sm.server_id, sm.roles
FROM server_members sm, users u 
WHERE sm.server_id = ANY($1::text[]) AND sm.user_id = u.id;

-- name: GetRolesFromServers :many
SELECT r.id, r.position, r.name, r.color, r.abilities, r.server_id
FROM roles r
WHERE r.server_id = ANY($1::text[])
ORDER BY r.position;

-- name: GetRolesFromServer :many
SELECT r.id, r.position, r.name, r.color, r.abilities, r.server_id
FROM roles r
WHERE r.server_id = $1
ORDER BY r.position;

-- name: CreateServer :one
INSERT INTO servers (
  id, owner_id, name, avatar, description, main_color, public
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: JoinServer :one
WITH ins AS (
  INSERT INTO server_members (id, user_id, server_id, position, roles)
  VALUES ($1, $2, $3, $4, '{}'::varchar[])
  RETURNING *
)
SELECT s.*, ins.roles, ins.position, (SELECT COUNT(*) FROM server_members smc WHERE smc.server_id = ins.server_id) AS member_count
FROM ins
JOIN servers s ON s.id = ins.server_id;

-- name: LeaveServer :exec
DELETE FROM server_members WHERE user_id = $1 AND server_id = $2;

-- name: UpdateServerAvatarNBanner :exec
UPDATE servers SET avatar = $1, banner = $2, main_color = $3 WHERE id = $4 AND owner_id = $5;

-- name: UpdateServerProfile :exec
UPDATE servers SET name = $1, description = $2 WHERE id = $3 AND owner_id = $4;

-- name: DeleteServer :execresult
DELETE FROM servers WHERE id = $1 AND owner_id = $2;

-- name: GetServersIDs :many
SELECT id FROM servers WHERE id <> 'global';

-- name: GetServersIDFromUser :many
SELECT server_id FROM server_members WHERE user_id = $1;