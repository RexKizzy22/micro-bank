package db

import (
	"context"
)

// contains the input parameters of the create user transaction
type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user User) error
}

// result of create user transaction
type CreateUserTxResult struct {
	User User
}

// creates a new user and enqueues an email verification task to the queue in one transaction
func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		// return arg.AfterCreate(result.User)
		return nil
	})

	return result, err
}
