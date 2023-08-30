package db

import (
	"context"
	"database/sql"
)

// contains the input parameters of the verify email transaction
type VerifyEmailTxParams struct {
	EmailId int64
	SecretCode string
}

// result of verify email transaction
type VerifyEmailTxResult struct {
	User User
	VerifyEmail VerifyEmail
}

// performs money transfer from one account to another
// creates a transfer record, adds account entries and updates accounts' balance in a single transaction
func (store *SQLStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.VerifyEmail, err = q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID: arg.EmailId,
			SecretCode: arg.SecretCode,
		})
		if err != nil {
			return err
		}

		result.User, err = q.UpdateUser(ctx, UpdateUserParams{
			IsEmailVerified: sql.NullBool{
				Bool: true,
				Valid: true,
			},
		})

		return err
	})

	return result, err
}
