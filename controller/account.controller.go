package controller

import (
	"net/http"
	"strconv"

	"github.com/512/simple_bank/dto"
	"github.com/512/simple_bank/service"

	"github.com/512/simple_bank/ultis"
	"github.com/gin-gonic/gin"
)

type AccountController interface {
	Create(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	GetList(ctx *gin.Context)
}

type accountController struct {
	accountService service.AccountService
}

func NewAccountController(accountService service.AccountService) AccountController {
	return &accountController{
		accountService: accountService,
	}
}

func (controller *accountController) Create(ctx *gin.Context) {
	var createAccountDTO dto.CreateAccountDTO
	errDTO := ctx.ShouldBind(&createAccountDTO)
	if errDTO != nil {
		res := ultis.BuildErrorResponse("Failed to bind DTO", errDTO.Error(), ultis.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	resultAccount, err := controller.accountService.CreateAccount(createAccountDTO)
	if err != nil {
		res := ultis.BuildErrorResponse("Failed to process create", err.Error(), ultis.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := ultis.BuildResponse(true, "", resultAccount)
	ctx.JSON(http.StatusOK, res)
}

func (controller *accountController) GetByID(ctx *gin.Context) {
	accountID, err := strconv.ParseInt(ctx.Param("id"), 0, 0)
	if err != nil {
		res := ultis.BuildErrorResponse("No param id was found", err.Error(), ultis.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	resultAccount, err := controller.accountService.GetAccountByID(accountID)
	if err != nil {
		res := ultis.BuildErrorResponse("Failed to process get", err.Error(), ultis.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := ultis.BuildResponse(true, "", resultAccount)
	ctx.JSON(http.StatusOK, res)
}

func (controller *accountController) GetList(ctx *gin.Context) {
	var getListAccountDTO dto.GetListAccountDTO
	errDTO := ctx.ShouldBindQuery(&getListAccountDTO)
	if errDTO != nil {
		res := ultis.BuildErrorResponse("Failed to bind DTO", errDTO.Error(), ultis.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	resultAccounts, err := controller.accountService.GetListAccount(getListAccountDTO)
	if err != nil {
		res := ultis.BuildErrorResponse("Failed to process get", err.Error(), ultis.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := ultis.BuildResponse(true, "", resultAccounts)
	ctx.JSON(http.StatusOK, res)

}
