package api

import (
	"net/http"
	"time"

	db "github.com/Rexkizzy22/simple-bank/db/sqlc"
	"github.com/Rexkizzy22/simple-bank/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username    string `json:"username" binding:"required, alphanum"`
	Email string `json:"email" binding:"required, email"`
	FullName string `json:"full_name" binding:"required"`
	Password string `json:"password" binding:"required, min=6"`
}

type createUserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (route *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	hashPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username: req.Username,
		FullName: req.FullName,
		Email: req.Email,
		HashedPassword: hashPassword,
	}

	user, err := route.store.CreateUser(ctx, arg)
	if err != nil {
		if pErr, ok := err.(*pq.Error); ok {
			switch pErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := createUserResponse{
		Username: user.Username,
		Email: user.Email,
		FullName: user.FullName,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt: user.CreatedAt,
	}

	ctx.JSON(http.StatusOK, res)
}