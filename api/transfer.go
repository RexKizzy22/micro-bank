package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/Rexkizzy22/simple-bank/db/sqlc"
	"github.com/Rexkizzy22/simple-bank/token"
	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required, min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required, min=1"`
	Amount        int64  `json:"amount" binding:"required, gt=0"`
	Currency      string `json:"currency" binding:"required, currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	fromAccount, isValid := server.validAccountCurrency(ctx, req.FromAccountID, req.Currency)
	if !isValid {
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Username != fromAccount.Owner {
		err := errors.New("account does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, isValid = server.validAccountCurrency(ctx, req.ToAccountID, req.Currency)
	if !isValid {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// validates that account currency is one of USD, EUR or CAD
func (server *Server) validAccountCurrency(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, nil)
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err = fmt.Errorf("account [%d]: currency mismatch: %s vs %s", accountID, currency, account.Currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}
