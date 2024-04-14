package api

import (
	"fmt"

	db "github.com/Rexkizzy22/micro-bank/db/sqlc"
	"github.com/Rexkizzy22/micro-bank/token"
	"github.com/Rexkizzy22/micro-bank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// serves HTTP requests for our banking services
type Server struct {
	config util.Config
	store  db.Store
	token  token.Maker
	router *gin.Engine
}

// creates a new HTTP server and set up routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config: config,
		store:  store,
		token:  tokenMaker,
	}

	// registers the custom validation for the currency request field
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validateCurrency)
	}

	server.Routes()
	return server, nil
}

func (server *Server) Routes() {
	router := gin.Default()

	if server.config.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	v1 := router.Group("/v1")

	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1.POST("/user", server.createUser)
	v1.POST("/user/login", server.loginUser)
	v1.POST("/tokens/renew_access", server.renewAccessToken)
	// TODO: log out handler
	// router.POST("/logout", server.logout)

	authRoutes := v1.Group("/").Use(authMiddleware(server.token))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccount)

	authRoutes.POST("/transfers", server.createTransfer)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err}
}
