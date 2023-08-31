-- name: CreateVerifyEmail :one
INSERT INTO verify_emails (
  username,
  email,
  secret_code
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: UpdateVerifyEmail :one
UPDATE verify_emails
SET 
  is_used = TRUE
WHERE 
  is_used = FALSE
  AND secret_code = @secret_code
  AND expires_at > now()
  AND id = @id 
RETURNING *;