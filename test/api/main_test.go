package test_api

import (
	"database/sql"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-backend-practice/api"
	"github.com/go-backend-practice/db"
	"github.com/go-backend-practice/util"
	_ "github.com/lib/pq"
)

var tx *db.Transaction
var server *api.Server

func TestMain(m *testing.M) {
	config, err := util.Loadconfig("../../", "test")
	if err != nil {
		panic(err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		panic(err)
	}

	tx = db.NewTransaction(conn)
	server = api.NewServer(tx)

	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}