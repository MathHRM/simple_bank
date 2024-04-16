package api

import (
	db "github.com/MathHRM/simple_bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	
	router.GET("/accounts", server.listAccounts)
	router.GET("/accounts/:id", server.getAccount)

	router.PUT("/accounts", server.updateAccount)

	router.DELETE("/accounts/:id", server.deleteAccount)

	server.router = router
	return server
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err}
}
