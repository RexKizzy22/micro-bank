package db

import (
	"context"
)

// contains the input parameters of the transfer transaction
type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user User) error
}

// result of transfer transaction
type CreateUserTxResult struct {
	User User
}

// performs money transfer from one account to another
// creates a transfer record, adds account entries and updates accounts' balance in a single transaction
func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		return arg.AfterCreate(result.User)
	})

	return result, err
}
