package apis

import (
	"github.com/gin-gonic/gin"
	db "github.com/vietthangc1/simple_bank/db/sqlc"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	s := &Server{
		store: store,
	}
	router := gin.Default()

	// TODO: add routes
	router.POST("/account/create", s.createAccount)
	router.GET("/account/get/:id", s.getAccountByID)
	router.GET("/account/list", s.listAccounts)

	s.router = router
	return s
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func errorResponse(err  error) gin.H {
	return gin.H{"error": err.Error()}
}