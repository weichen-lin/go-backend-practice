package main

import (
	"database/sql"
	"log"

	"github.com/go-backend-practice/api"
	"github.com/go-backend-practice/db"
	"github.com/go-backend-practice/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.Loadconfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		panic(err)
	}

	transaction := db.NewTransaction(conn)
	server := api.NewServer(transaction)

	startServerErr := server.Start(config.ServerAddress)

	if startServerErr != nil {
		log.Fatal("cannot start server: ", startServerErr)
	}
}
