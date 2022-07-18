package main

import (
	"database/sql"

	"github.com/512/simple_bank/config"
	"github.com/512/simple_bank/controller"
	db "github.com/512/simple_bank/db/sqlc"
	"github.com/512/simple_bank/service"
	"github.com/512/simple_bank/validation"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	connection *sql.DB   = config.ConnectDB()
	store      *db.Store = db.NewStore(connection)
	//service
	accountService  service.AccountService  = service.NewAccountService(store)
	transferService service.TransferService = service.NewTransferService(store)
	//controller
	accountController  controller.AccountController  = controller.NewAccountController(accountService)
	transferController controller.TransferController = controller.NewTransferController(transferService, accountService)
)

func main() {
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validation.ValidCurrency)
	}

	accountRouter := router.Group("/api/v1/account")
	{
		accountRouter.POST("/", accountController.Create)
		accountRouter.GET("/:id", accountController.GetByID)
		accountRouter.GET("/", accountController.GetList)
	}

	transferRouter := router.Group("/api/v1/transfer")
	{
		transferRouter.POST("/", transferController.CreateTransfer)

	}

	router.Run()
}
