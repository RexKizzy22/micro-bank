// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: verify_email.sql

package db

import (
	"context"
)

const createVerifyEmail = `-- name: CreateVerifyEmail :one
INSERT INTO verify_emails (
  username,
  email,
  secret_code
) VALUES (
  $1, $2, $3
) RETURNING id, username, email, secret_code, is_used, created_at, expired_at
`

type CreateVerifyEmailParams struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	SecretCode string `json:"secret_code"`
}

func (q *Queries) CreateVerifyEmail(ctx context.Context, arg CreateVerifyEmailParams) (VerifyEmail, error) {
	row := q.db.QueryRowContext(ctx, createVerifyEmail, arg.Username, arg.Email, arg.SecretCode)
	var i VerifyEmail
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.SecretCode,
		&i.IsUsed,
		&i.CreatedAt,
		&i.ExpiredAt,
	)
	return i, err
}

const updateVerifyEmail = `-- name: UpdateVerifyEmail :one
UPDATE verify_emails
SET 
  is_used = TRUE
WHERE 
  is_used = FALSE
  AND secret_code = $1
  AND expires_at > now()
  AND id = $2 
RETURNING id, username, email, secret_code, is_used, created_at, expired_at
`

type UpdateVerifyEmailParams struct {
	SecretCode string `json:"secret_code"`
	ID         int64  `json:"id"`
}

func (q *Queries) UpdateVerifyEmail(ctx context.Context, arg UpdateVerifyEmailParams) (VerifyEmail, error) {
	row := q.db.QueryRowContext(ctx, updateVerifyEmail, arg.SecretCode, arg.ID)
	var i VerifyEmail
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.SecretCode,
		&i.IsUsed,
		&i.CreatedAt,
		&i.ExpiredAt,
	)
	return i, err
}