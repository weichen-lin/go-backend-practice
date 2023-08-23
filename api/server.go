package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-backend-practice/db"
)

type Server struct {
	transaction *db.Transaction
	router      *gin.Engine
}

func NewServer(transaction *db.Transaction) *Server {

	server := &Server{
		transaction: transaction,
	}

	router := gin.Default()

	router.GET("/account", server.listAccount)
	router.GET("/account/:id", server.getAccount)
	router.POST("/account", server.createAccount)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
