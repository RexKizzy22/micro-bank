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

	

	server.router = router
	return server
}

