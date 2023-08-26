package test_db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/go-backend-practice/db"
	"github.com/go-backend-practice/util"
	_ "github.com/lib/pq"
)

var testTx *db.Transaction
var query *db.Queries
var dbConn *sql.DB

func TestMain(m *testing.M) {
	var err error
	config, err := util.Loadconfig("../../", "test")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	dbConn, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testTx = db.NewTransaction(dbConn)
	query = db.New(dbConn)

	os.Exit(m.Run())
}
