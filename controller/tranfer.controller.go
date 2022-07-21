package controller

import (
	"errors"
	"fmt"
	"net/http"

	db "github.com/512/simple_bank/db/sqlc"
	"github.com/512/simple_bank/dto"
	"github.com/512/simple_bank/service"
	"github.com/512/simple_bank/service/token"
	"github.com/512/simple_bank/ultis"
	"github.com/gin-gonic/gin"
)

type TransferController interface {
	CreateTransfer(ctx *gin.Context)
	ValidAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, error)
}

type transferController struct {
	accountService  service.AccountService
	transferService service.TransferService
}

func NewTransferController(transferService service.TransferService, accountService service.AccountService) TransferController {
	return &transferController{
		transferService: transferService,
		accountService:  accountService,
	}

}

func (controller *transferController) CreateTransfer(ctx *gin.Context) {
	var createTransferDTO dto.CreateTransferDTO
	if errDTO := ctx.ShouldBind(&createTransferDTO); errDTO != nil {
		res := ultis.BuildErrorResponse("Failed to bind DTO", errDTO.Error(), ultis.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	fromAccount, err := controller.ValidAccount(ctx, createTransferDTO.FromAccountID, createTransferDTO.Currency)
	if err != nil {
		res := ultis.BuildErrorResponse("Failed to process get", err.Error(), ultis.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	authPayload := ctx.MustGet("payload").(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("from account doesn't belong to the authenticated user")
		res := ultis.BuildErrorResponse("Failed to process get", err.Error(), ultis.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	_, err = controller.ValidAccount(ctx, createTransferDTO.FromAccountID, createTransferDTO.Currency)
	if err != nil {
		res := ultis.BuildErrorResponse("Failed to process get", err.Error(), ultis.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	resultTransferTx, err := controller.transferService.CreateTransfer(createTransferDTO)
	if err != nil {
		res := ultis.BuildErrorResponse("Failed to process create", err.Error(), ultis.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := ultis.BuildResponse(true, "", resultTransferTx)
	ctx.JSON(http.StatusOK, res)

}

func (controller *transferController) ValidAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, error) {
	account, err := controller.accountService.GetAccountByID(accountID)
	if err != nil {
		return account, err
	}
	if account.Currency != currency {
		messageErr := fmt.Sprintf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		err := errors.New(messageErr)
		return account, err
	}
	return account, nil
}
