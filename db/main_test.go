package db

import (
	"context"
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
var testTx *Transaction

// https://darjun.github.io/2021/08/03/godailylib/testing/
func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testTx = NewTransaction(conn)

	tx, txerr := testTx.db.BeginTx(context.Background(), nil)
	if txerr != nil {
		log.Fatal("transaction init error: ", txerr)
	}

	testQuries = New(tx)

	os.Exit(m.Run())
}
