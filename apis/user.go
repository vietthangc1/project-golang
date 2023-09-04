package apis

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/vietthangc1/simple_bank/db/sqlc"
	"github.com/vietthangc1/simple_bank/pkg/envx"
	"github.com/vietthangc1/simple_bank/pkg/passwordx"
)

type createUserRequest struct {
	Username string `json:"user_name" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type createUserResponse struct {
	Username  string    `json:"user_name"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	passwordManager := passwordx.NewPassword(int(envx.Int("BYCRYPT_COST", 10)))
	hashedPassword, err := passwordManager.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := createUserResponse{
		Username:  user.Username,
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	ctx.JSON(http.StatusOK, resp)
}

var (
	ErrWrongUsernamePassword = errors.New("wrong username or password")
)

type loginUserRequest struct {
	Username string `json:"user_name" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	Username    string    `json:"user_name"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
	AccessToken string    `json:"access_token"`
}

func (s *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, errorResponse(ErrWrongUsernamePassword))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	passwordManager := passwordx.NewPassword(int(envx.Int("BYCRYPT_COST", 10)))
	err = passwordManager.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(ErrWrongUsernamePassword))
		return
	}

	accessToken, err := s.token.CreateToken(user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := loginUserResponse{
		Username:    user.Username,
		FullName:    user.FullName,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt,
		AccessToken: accessToken,
	}

	ctx.JSON(http.StatusOK, resp)
}
