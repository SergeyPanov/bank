package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/SergeyPanov/bank/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testStore Store

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../.")
	if err != nil {
		log.Fatal("cannot read config", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("can't connect to db")
	}

	testStore = NewStore(connPool)

	os.Exit(m.Run())
}
