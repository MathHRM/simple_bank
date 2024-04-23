package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/MathHRM/simple_bank/db/sqlc"
	"github.com/MathHRM/simple_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Email    string `json:"email" binding:"required,email" `
	Password string `json:"password" binding:"required,min=8" `
	FullName string `json:"full_name" binding:"required" `
}

type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse {
		Username: user.Username,
		Email: user.Email,
		FullName: user.FullName,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt: user.CreatedAt,
	}
}

func (s *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		Email:          req.Email,
		FullName:       req.FullName,
	}

	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
			}
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userResponse := userResponse {
		Username: user.Username,
		Email: user.Email,
		FullName: user.FullName,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt: user.CreatedAt,
	}

	ctx.JSON(http.StatusOK, userResponse)
}

type getUserRequest struct {
	Username string `uri:"username" binding:"required,alphanum"`
}

func (s *Server) getUser(ctx *gin.Context) {
	var req getUserRequest
	
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userResponse := newUserResponse(user)

	ctx.JSON(http.StatusOK, userResponse)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=8" `
}

type loginUserResponse struct {
	AcessToken string `json:"acess_token" binding:"required"`
	User userResponse `json:"user" binding:"required"`
}

func (s *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CompareHashPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	token, err := s.tokenMaker.CreateToken(user.Username, s.config.AcessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	loginResponse := loginUserResponse {
		AcessToken: token,
		User: newUserResponse(user),
	}
	
	ctx.JSON(http.StatusOK, loginResponse)
}