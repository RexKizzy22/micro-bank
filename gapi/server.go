package gapi

import (
	"fmt"

	db "github.com/Rexkizzy22/micro-bank/db/sqlc"
	"github.com/Rexkizzy22/micro-bank/pb"
	"github.com/Rexkizzy22/micro-bank/task"
	"github.com/Rexkizzy22/micro-bank/token"
	"github.com/Rexkizzy22/micro-bank/util"
)

// serves gRPC requests for our banking services
type Server struct {
	pb.UnimplementedMicroBankServer
	config          util.Config
	store           db.Store
	token           token.Maker
	taskDistributor task.TaskDistributor
}

// creates a new gRPC server
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config: config,
		store:  store,
		token:  tokenMaker,
		// taskDistributor: taskDistributor,
	}

	return server, nil
}
