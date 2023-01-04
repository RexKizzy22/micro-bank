-- name: CreateSession :one
INSERT INTO sessions (
  id,
  username,
  is_blocked,
  refresh_token,
  user_agent,
  client_ip,
  expires_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = $1
LIMIT 1;