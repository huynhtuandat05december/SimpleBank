package main

import (
	"database/sql"

	"github.com/512/simple_bank/config"
	"github.com/512/simple_bank/controller"
	db "github.com/512/simple_bank/db/sqlc"
	"github.com/512/simple_bank/service"
	"github.com/gin-gonic/gin"
)

var (
	connection *sql.DB   = config.ConnectDB()
	store      *db.Store = db.NewStore(connection)
	//service
	accountService service.AccountService = service.NewAccountService(store)
	//controller
	accountController controller.AccountController = controller.NewAccountController(accountService)
)

func main() {
	router := gin.Default()

	accountRouter := router.Group("/api/v1/account")
	{
		accountRouter.POST("/", accountController.Create)
		accountRouter.GET("/:id", accountController.GetByID)
		accountRouter.GET("/", accountController.GetList)
	}

	router.Run()
}
