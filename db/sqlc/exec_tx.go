package db

import (
	"context"
	"fmt"
)

// executes a function within a transaction
func (s *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, tx rollback: %v", err, rbErr)
		}
	}

	return tx.Commit(ctx)
}
