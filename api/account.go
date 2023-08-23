package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-backend-practice/db"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateAccountParams struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var params CreateAccountParams
	if paramsErr := ctx.ShouldBindJSON(&params); paramsErr != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(paramsErr))
		return
	}

	arg := db.CreateAccountParams{
		Owner:    params.Owner,
		Balance:  decimal.NewFromFloat(0),
		Currency: params.Currency,
	}

	var account db.Account

	createAccountErr := server.transaction.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		account, err = q.CreateAccount(ctx, arg)
		return err
	}, false)

	if createAccountErr != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(createAccountErr))
	}

	ctx.JSON(http.StatusOK, account)
}

type GetAccountParams struct {
	ID string `uri:"id" binding:"required"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var params GetAccountParams

	if paramsErr := ctx.ShouldBindUri(&params); paramsErr != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(paramsErr))
		return
	}

	var account db.Account

	getAccountErr := server.transaction.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		uuid, parseUUIDErr := uuid.Parse(params.ID)
		if parseUUIDErr != nil {
			return parseUUIDErr
		}
		account, err = q.GetAccount(ctx, uuid)
		return err
	}, false)

	if getAccountErr != nil {
		if getAccountErr == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(getAccountErr))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(getAccountErr))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type listAccountParams struct {
	Start    int32 `form:"start" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccount(ctx *gin.Context) {
	var params listAccountParams

	if paramsErr := ctx.ShouldBindQuery(&params); paramsErr != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(paramsErr))
		return
	}

	var accounts []db.Account

	arg := db.ListAccountsParams{
		Limit:  params.PageSize,
		Offset: params.Start,
	}

	listAccountErr := server.transaction.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		accounts, err = q.ListAccounts(ctx, arg)
		return err
	}, false)

	if listAccountErr != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(listAccountErr))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
