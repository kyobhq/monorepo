-- name: GetServer :one
SELECT * FROM servers WHERE id = $1;

-- name: GetServers :many
SELECT * FROM servers;

-- name: OwnServer :execresult
SELECT * FROM servers WHERE id = $1 AND owner_id = $2;

-- name: IsMember :execresult
SELECT id FROM server_members WHERE server_id = $1 AND user_id = $2;

-- name: GetServersFromUser :many
SELECT DISTINCT s.*, sm.roles, sm.position, (SELECT count(id) FROM server_members smc WHERE smc.server_id=s.id AND ban=false) AS member_count
FROM servers s
INNER JOIN server_members sm ON sm.server_id = s.id AND sm.user_id = $1 AND sm.ban = false
WHERE s.id <> 'global';

-- name: GetServerIDsFromUser :many
SELECT DISTINCT s.id
FROM servers s
INNER JOIN server_members sm ON sm.server_id = s.id AND sm.user_id = $1
WHERE s.id <> 'global';

-- name: GetServerMembers :many
SELECT 
    u.id, 
    u.username, 
    u.display_name, 
    u.avatar,
    u.created_at as joined_kyob,
    sm.roles,
    sm.created_at as joined_server,
    CASE 
        WHEN u.id = ANY($3::text[]) THEN 'online'
        ELSE 'offline'
    END as status,
    COALESCE(MIN(r.position), 999999) as min_role_position
FROM server_members sm
JOIN users u ON u.id = sm.user_id
LEFT JOIN roles r ON r.server_id = sm.server_id AND r.id = ANY(sm.roles)
WHERE sm.server_id = $1 AND sm.ban = false
GROUP BY sm.user_id, u.id, u.username, u.display_name, u.avatar, sm.roles, sm.created_at
ORDER BY 
    CASE WHEN u.id = ANY($3::text[]) THEN 0 ELSE 1 END,
    min_role_position, 
    u.username
LIMIT 50 OFFSET $2;

-- name: GetServerInformations :one
SELECT 
    (
        SELECT json_agg(json_build_object(
            'id', ranked_members.user_id,
            'username', ranked_members.username,
            'display_name', ranked_members.display_name,
            'avatar', ranked_members.avatar,
            'roles', ranked_members.roles,
            'joined_server', ranked_members.joined_server,
            'joined_kyob', ranked_members.joined_kyob,
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
                u.created_at as joined_kyob,
                sm.roles,
                sm.created_at as joined_server,
                COALESCE(MIN(r.position), 999999) as min_role_position
            FROM server_members sm
            JOIN users u ON u.id = sm.user_id
            LEFT JOIN roles r ON r.server_id = sm.server_id AND r.id = ANY(sm.roles)
            WHERE sm.server_id = s.id AND sm.ban = false
            GROUP BY sm.user_id, u.id, u.username, u.display_name, u.avatar, sm.roles, sm.created_at
            ORDER BY 
                CASE WHEN u.id = ANY($2::text[]) THEN 0 ELSE 1 END,
                min_role_position, 
                u.username
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
        SELECT json_agg(json_build_object(
            'id', i.id,
            'creator', json_build_object(
                'display_name', u.display_name,
                'username', u.username,
                'avatar', u.avatar
            ),
            'server_id', i.server_id,
            'invite_id', i.invite_id,
            'expire_at', i.expire_at
        ))
        FROM invites i
        LEFT JOIN users u ON i.creator_id = u.id
        WHERE i.server_id = s.id
    ) as invites,
    (
        SELECT smc.roles
        FROM server_members smc
        WHERE smc.server_id = s.id AND smc.user_id = $3
    ) as user_roles,
    (
        SELECT count(id)
        FROM server_members smc
        WHERE smc.server_id = s.id AND ban = false
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
  VALUES ($1, $2, $3, $4, $5)
  RETURNING *
)
SELECT s.*, ins.roles, ins.position, (SELECT COUNT(*) FROM server_members smc WHERE smc.server_id = ins.server_id AND ban=false) AS member_count
FROM ins
JOIN servers s ON s.id = ins.server_id;

-- name: LeaveServer :exec
DELETE FROM server_members WHERE user_id = $1 AND server_id = $2;

-- name: UpdateServerAvatarNBanner :exec
UPDATE servers
  set 
    avatar = COALESCE($2, avatar),
    banner = COALESCE($3, banner),
    updated_at = now()
WHERE id = $1;

-- name: UpdateServerProfile :exec
UPDATE servers SET name = $1, description = $2, public = $3, updated_at = now() WHERE id = $4;

-- name: DeleteServer :exec
DELETE FROM servers WHERE id = $1 AND owner_id = $2;

-- name: GetServersIDs :many
SELECT id FROM servers;

-- name: GetServersIDFromUser :many
SELECT server_id FROM server_members WHERE user_id = $1 AND ban = false;

-- name: BanUser :exec
UPDATE server_members
  set ban = true, ban_reason = $3
WHERE user_id = $1 AND server_id = $2;

-- name: KickUser :exec
DELETE FROM server_members WHERE user_id = $1 AND server_id = $2;

-- name: CheckBan :one
SELECT ban_reason FROM server_members WHERE user_id = $1 AND server_id = $2 AND ban = true;

-- name: GetBannedMembers :many
SELECT u.id, u.display_name, u.avatar, u.username
FROM server_members sm
INNER JOIN users u ON sm.user_id = u.id
WHERE sm.ban = true AND sm.server_id = $1;
