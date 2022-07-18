package controller

import (
	"fmt"
	"net/http"

	"github.com/512/simple_bank/dto"
	"github.com/512/simple_bank/service"
	"github.com/512/simple_bank/ultis"
	"github.com/gin-gonic/gin"
)

type TransferController interface {
	CreateTransfer(ctx *gin.Context)
	SameAccountCurrency(ctx *gin.Context, accountID int64, currency string) bool
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

	if !controller.SameAccountCurrency(ctx, createTransferDTO.FromAccountID, createTransferDTO.Currency) {
		return
	}

	if !controller.SameAccountCurrency(ctx, createTransferDTO.ToAccountID, createTransferDTO.Currency) {
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

func (controller *transferController) SameAccountCurrency(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := controller.accountService.GetAccountByID(accountID)
	if err != nil {
		res := ultis.BuildErrorResponse("Failed to process get", err.Error(), ultis.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return false
	}

	if account.Currency != currency {
		err := fmt.Sprintf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		res := ultis.BuildErrorResponse("Failed to process transfer", err, ultis.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return false
	}

	return true
}
