-- name: UpsertRole :one
INSERT INTO roles (
  id, position, server_id, name, color, abilities
) VALUES (
  $1, $2, $3, $4, $5, $6
)
ON CONFLICT (id) 
DO UPDATE SET 
  name = EXCLUDED.name,
  color = EXCLUDED.color,
  abilities = EXCLUDED.abilities
RETURNING *;

-- name: GetUserAbilities :one
SELECT COALESCE(array_agg(DISTINCT ability), '{}'::text[])::text[] as abilities
FROM (
  SELECT unnest(r.abilities) as ability
  FROM roles r, server_members sm 
  WHERE sm.server_id = $1 AND sm.user_id = $2 AND r.id = ANY(sm.roles)
  
  UNION
  
  SELECT 'OWNER' as ability
  FROM servers s
  WHERE s.id = $1 AND s.owner_id = $2
) subquery;

-- name: CheckPermission :one
SELECT 1
FROM servers s
WHERE s.id = $1
  AND (
  s.owner_id = $2
  OR EXISTS (
    SELECT 1
    FROM server_members sm
    JOIN roles r ON r.id = ANY(sm.roles)
    WHERE sm.server_id = s.id AND sm.user_id = $2 AND $3::text = ANY(r.abilities)
  )
);

-- name: GetUserRolesFromServers :many
SELECT sm.server_id, sm.roles FROM server_members sm WHERE sm.user_id = $1 AND sm.server_id = ANY($2::text[]);

-- name: GetRoles :many
SELECT r.id, r.position, r.name, r.color, r.abilities, array_agg(sm.user_id) FILTER (WHERE sm.user_id IS NOT NULL) AS members
FROM roles r
LEFT JOIN server_members sm on r.id = ANY(sm.roles)
WHERE r.server_id = $1
GROUP BY r.id;

-- name: GetRole :one
SELECT r.id, r.position, r.name, r.color, r.abilities, r.server_id
FROM roles r
WHERE r.id = $1;

-- name: GiveRole :exec
UPDATE server_members 
SET roles = array_append(roles, $1) -- role_name
WHERE server_id = $2 AND user_id = $3;

-- name: DeleteRole :exec
DELETE FROM roles WHERE id = $1;

-- name: RemoveRoleMember :exec
UPDATE server_members SET roles = array_remove(roles, $1) WHERE server_id = $2 AND user_id = $3;

-- name: RemoveRoleFromAllMembers :exec
UPDATE server_members SET roles = array_remove(roles, $1) WHERE $1 = ANY(roles);

-- name: MoveRole :exec
UPDATE roles SET position = $1 WHERE id = $2;

-- name: UpdateRolePositions :exec
UPDATE roles SET position = position + 1 WHERE position >= $1 AND position < $2;
