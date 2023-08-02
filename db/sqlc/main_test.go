package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/SergeyPanov/bank/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../.")
	if err != nil {
		log.Fatal("cannot read config", err)
	}

	testDb, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can't connect to db")
	}

	testQueries = New(testDb)

	os.Exit(m.Run())
}
