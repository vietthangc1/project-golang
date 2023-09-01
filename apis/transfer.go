package apis

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/vietthangc1/simple_bank/db/sqlc"
)

type transferRequest struct {
	FromAccountID int64 `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64 `json:"to_account_id" binding:"required,min=1"`
	Amount        int64 `json:"amount" binding:"required,min=1"`
}

func (s *Server) transfer(ctx *gin.Context) {
	var req transferRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := s.validAccount(ctx, req.FromAccountID); err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	if err := s.validAccount(ctx, req.ToAccountID); err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	account, err := s.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (s *Server) validAccount(ctx context.Context, accountId int64) error {
	_, err := s.store.GetAccountByID(ctx, accountId)
	if err != nil {
		return err
	}
	return nil
}
