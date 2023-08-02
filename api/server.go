package api

import (
	db "github.com/SergeyPanov/bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{
		store: store,
	}

	router := gin.Default()

	router.POST("/account", server.createAccount)
	router.GET("/account/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)

	server.router = router

	return server
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
