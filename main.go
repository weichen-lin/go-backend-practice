package main

import (
	"database/sql"
	"log"

	"github.com/go-backend-practice/api"
	"github.com/go-backend-practice/db"
	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:test_local@localhost:5432/bank?sslmode=disable"
	serverAddress = "0.0.0.0:8000"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		panic(err)
	}

	transaction := db.NewTransaction(conn)
	server := api.NewServer(transaction)

	startServerErr := server.Start(serverAddress)

	if startServerErr != nil {
		log.Fatal("cannot start server: ", startServerErr)
	}
}
