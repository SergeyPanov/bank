package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/SergeyPanov/bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (s *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	if !s.validAccount(ctx, req.FromAccountID, req.Currency) {
		return
	}

	if !s.validAccount(ctx, req.ToAccountID, req.Currency) {
		return
	}

	arg := db.TransferTxParams{
		FromAccountId: req.FromAccountID,
		ToAccountId:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := s.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (s *Server) validAccount(ctx *gin.Context, accoundId int64, currency string) bool {
	acc, err := s.store.GetAccount(ctx, accoundId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return false
		}

		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return false
	}

	if acc.Currency.String != currency {
		err := fmt.Errorf("acc [%d] currency mismatch: %s vs %s", acc.ID, acc.Currency.String, currency)
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return false
	}

	return true
}
