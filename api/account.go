package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/Rexkizzy22/micro-bank/db/sqlc"
	"github.com/Rexkizzy22/micro-bank/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

// @Summary  creates a new account with a given currency
// @Accepts  json
// @Produce  json
// @Param    currency body     string true "Account Currency"
// @Success  200      {object} db.Account
// @Security ApiKeyAuth
// @Router   /accounts [POST]
func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if pErr, ok := err.(*pq.Error); ok {
			switch pErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// @Summary  gets an existing account for a specified User ID
// @Accepts  json
// @Produce  json
// @Param    id  path     string true "User ID"
// @Success  200 {object} db.Account
// @Security ApiKeyAuth
// @Router   /accounts/:id [GET]
func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Username != account.Owner {
		err := errors.New("account does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
	Page_ID   int32 `form:"page_id" binding:"required,min=1"`
	Page_Size int32 `form:"page_size" binding:"required,min=5,max=10"`
}

// @Summary  lists all accounts for a specific User ID
// @Accepts  json
// @Produce  json
// @Param    currency body    string true "Account Currency"
// @Success  200      {array} []db.Account
// @Security ApiKeyAuth
// @Router   /accounts [GET]
func (server *Server) listAccount(ctx *gin.Context) {
	var req listAccountRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.Page_Size,
		Offset: (req.Page_ID - 1) * req.Page_Size,
	}

	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
