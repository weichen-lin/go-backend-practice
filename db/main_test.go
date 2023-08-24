package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/go-backend-practice/util"
	_ "github.com/lib/pq"
)

var testTx *Transaction

var sharedConn *sql.DB

// https://darjun.github.io/2021/08/03/godailylib/testing/
func TestMain(m *testing.M) {
	var connErr error
	config, configErr := util.Loadconfig("../")
	if configErr != nil {
		log.Fatal("cannot load config: ", configErr)
	}

	sharedConn, connErr = sql.Open(config.DBDriver, config.DBSource)

	if connErr != nil {
		log.Fatal("cannot connect to db: ", connErr)
	}

	testTx = NewTransaction(sharedConn)

	os.Exit(m.Run())
}
