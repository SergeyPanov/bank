package api

import (
	"fmt"

	db "github.com/SergeyPanov/bank/db/sqlc"
	"github.com/SergeyPanov/bank/token"
	"github.com/SergeyPanov/bank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store      db.Store
	config     util.Config
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		store:      store,
		config:     config,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil
}

func (s *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", s.createUser)
	router.POST("/users/login", s.loginUser)

	authRouts := router.Group("/").Use(authMiddleware(s.tokenMaker))

	authRouts.POST("/accounts", s.createAccount)
	authRouts.GET("/accounts/:id", s.getAccount)
	authRouts.GET("/accounts", s.listAccount)

	authRouts.POST("/transfers", s.createTransfer)

	s.router = router
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
