package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/Rexkizzy22/micro-bank/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("unable to load config: ", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("unable to connect to database: ", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}

// TODO: Write tests for entry.sql.go queries
// TODO: Write tests for transfer.sql.go queries
