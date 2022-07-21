package main

import (
	"database/sql"
	"log"

	"github.com/Rexkizzy22/simple-bank/api"
	db "github.com/Rexkizzy22/simple-bank/db/sqlc"
	"github.com/Rexkizzy22/simple-bank/docs"
	"github.com/Rexkizzy22/simple-bank/util"
	_ "github.com/lib/pq"
)

// @securitydefinitions.apiKey  ApiAuthKey
// @in                          header
// @name                        Authorization
func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("unable to load config: ", err)
	}

	// programmatically setting general swagger info
	docs.SwaggerInfo.Title = "Simple Bank API"
	docs.SwaggerInfo.Description = "A production-grade Go API that provides money transfer services between accounts of registered users"
	docs.SwaggerInfo.Version = "1.0.0"
	docs.SwaggerInfo.Host = config.ServerAddress
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("unable to connect to database: ", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("unable to start server: ", err)
	}
}
