-- name: VerifyToken :one
SELECT * FROM users WHERE id = (SELECT user_id FROM tokens WHERE token = $1);

-- name: DeleteRememberMeToken :exec
DELETE FROM tokens WHERE user_id = $1 AND type = 'REMEMBER_ME_TOKEN';

-- name: CreateToken :one
INSERT INTO tokens (
  id, user_id, token, expire_at, type
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;
