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
	dbSource = "postgresql://root:test_local@localhost:5432/bank?sslmode=disable"
)

var testQuries *Queries
var testTxConn *sql.DB

// https://darjun.github.io/2021/08/03/godailylib/testing/
func TestMain(m *testing.M) {
	var err error
	testTxConn, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testQuries = New(testTxConn)

	os.Exit(m.Run())
}
