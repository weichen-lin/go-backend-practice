package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = ""
)

var testTx *Transaction

var sharedConn *sql.DB

// https://darjun.github.io/2021/08/03/godailylib/testing/
func TestMain(m *testing.M) {
	var connErr error
	sharedConn, connErr = sql.Open(dbDriver, dbSource)

	if connErr != nil {
		log.Fatal("cannot connect to db: ", connErr)
	}

	testTx = NewTransaction(sharedConn)

	os.Exit(m.Run())
}
