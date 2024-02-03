package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/Rexkizzy22/micro-bank/util"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testStore Store

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("unable to load config: ", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("unable to connect to database: ", err)
	}

	testStore = NewStore(connPool)

	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

// TODO: Write tests for entry.sql.go queries
// TODO: Write tests for transfer.sql.go queries
