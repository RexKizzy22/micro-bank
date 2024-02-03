package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// establishes a contract that all store instances must satisfy, facilitates easy mocking for tests
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error)
}

// provides all functions to execute SQL queries and transactions
type SQLStore struct {
	*Queries
	connPool *pgxpool.Pool
}

// creates a new store instance
func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
