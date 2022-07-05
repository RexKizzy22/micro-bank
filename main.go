package main

import (
	"database/sql"
	"log"

	"github.com/Rexkizzy22/simple-bank/api"
	db "github.com/Rexkizzy22/simple-bank/db/sqlc"
	"github.com/Rexkizzy22/simple-bank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("unable to load config: ", err)
	}

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
