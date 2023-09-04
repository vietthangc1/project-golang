package apis

import (
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/vietthangc1/simple_bank/db/sqlc"
	"github.com/vietthangc1/simple_bank/pkg/envx"
	"github.com/vietthangc1/simple_bank/pkg/tokenx"
)

var (
	tokenSymmetrickey = envx.String("TOKEN_SECRET_KEY", "12331231231231231231231233123312")
	tokenDuration     = envx.String("TOKEN_DURATION", "24h")
)

type Server struct {
	store  db.Store
	token  tokenx.Token
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	timeParseDuration, err := time.ParseDuration(tokenDuration)
	if err != nil {
		panic(err)
	}
	tokenManager, err := tokenx.NewPasetoImpl(tokenSymmetrickey, timeParseDuration)
	if err != nil {
		panic(err)
	}
	s := &Server{
		store: store,
		token: tokenManager,
	}
	router := gin.Default()

	// TODO: add routes
	router.POST("/account/create", s.createAccount)
	router.GET("/account/get/:id", s.getAccountByID)
	router.GET("/account/list", s.listAccounts)

	router.POST("/transfer", s.transfer)

	router.POST("/user/create", s.createUser)
	router.POST("/login", s.loginUser)

	s.router = router
	return s
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
