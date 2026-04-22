-- name: InsertSession :one
INSERT INTO sessions (user_id, refresh_token, ip_address, user_agent, expires_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: FindSessionByToken :one
SELECT * FROM sessions
WHERE refresh_token = $1
  AND revoked_at IS NULL
  AND expires_at > NOW()
LIMIT 1;

-- name: DeleteSession :exec
UPDATE sessions
SET revoked_at = NOW()
WHERE refresh_token = $1;

-- name: DeleteAllUserSessions :exec
UPDATE sessions
SET revoked_at = NOW()
WHERE user_id = $1
  AND revoked_at IS NULL;
