package api

import (
	db "github.com/Rexkizzy22/simple-bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

// serves HTTP requests for our banking services
type Server struct {
	store *db.Store
	router *gin.Engine
}

// creates a new HTTP server and set up routing
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.POST("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err}
}



