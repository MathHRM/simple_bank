package api

import (
	"fmt"

	db "github.com/MathHRM/simple_bank/db/sqlc"
	"github.com/MathHRM/simple_bank/token"
	"github.com/MathHRM/simple_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSimetricKey)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar token maker, %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil
}

func (s *Server) setupRouter() {
	router := gin.Default()

	router.GET("/users/:username", s.getUser)
	router.POST("/user/login", s.loginUser)
	router.POST("/users", s.createUser)

	authRouter := router.Group("/").Use( authMiddleware(s.tokenMaker) )

	authRouter.GET("/accounts", s.listAccounts)
	authRouter.GET("/accounts/:id", s.getAccount)
	authRouter.POST("/accounts", s.createAccount)
	authRouter.PUT("/accounts", s.updateAccount)
	authRouter.DELETE("/accounts/:id", s.deleteAccount)

	authRouter.GET("/transfer/:id", s.getTransfer)
	authRouter.POST("/transfer", s.createTransfer)

	s.router = router
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
