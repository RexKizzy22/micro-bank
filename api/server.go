package api

import (
	"fmt"

	db "github.com/Rexkizzy22/simple-bank/db/sqlc"
	"github.com/Rexkizzy22/simple-bank/token"
	"github.com/Rexkizzy22/simple-bank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
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

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validateCurrency)
	}

	server.Routes()
	return server, nil
}

func (server *Server) Routes() {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	
	authRoutes := router.Group("/").Use(authMiddleware(server.token))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.POST("/accounts/:id", server.getAccount)
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
