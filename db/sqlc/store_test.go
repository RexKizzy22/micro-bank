package db

import (
	"context"
	// "fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	n := 5
	amount := int64(10)
	existed := make(map[int]bool)

	// simulate n concurrent transfer transactions in n goroutines
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()

	}

	// test result or error of each transfer transaction performed in n goroutines
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// fmt.Println(">>>>before transfer", account1.Balance, account2.Balance)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)

		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotEmpty(t, fromEntry.ID)
		require.NotEmpty(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotEmpty(t, toEntry.ID)
		require.NotEmpty(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check accounts' balance
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff2%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ...., n * amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true

		// TODO:check the final updated account
		// updatedAccount1, err := testQueries.GetAccount(context.Background(), fromAccount.ID)
		// require.NoError(t, err)

		// updatedAccount2, err := testQueries.GetAccount(context.Background(), toAccount.ID)
		// require.NoError(t, err)

		// fmt.Println(">>>>after transfer", updatedAccount1.Balance, updatedAccount2.Balance)

		// require.Equal(t, account1.Balance-(int64(n)*amount), updatedAccount1.Balance)
		// require.Equal(t, account2.Balance+(int64(n)*amount), updatedAccount2.Balance)
	}
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// fmt.Println(">>>>before transfer", account1.Balance, account2.Balance)

	errs := make(chan error)

	// run n concurrent transfer transactions
	n := 10
	amount := int64(10)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()

	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		// TODO: check the final updated account
		// updatedAccount1, err := testQueries.GetAccount(context.Background(), fromAccountID.ID)
		// require.NoError(t, err)

		// updatedAccount2, err := testQueries.GetAccount(context.Background(), toAccount.ID)
		// require.NoError(t, err)

		// fmt.Println(">>>>after transfer", updatedAccount1.Balance, updatedAccount2.Balance)

		// require.Equal(t, account1.Balance, updatedAccount1.Balance)
		// require.Equal(t, account2.Balance, updatedAccount2.Balance)
	}
}
