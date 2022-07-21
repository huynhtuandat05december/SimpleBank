package controller

import (
	"net/http"
	"os"
	"time"

	db "github.com/512/simple_bank/db/sqlc"
	"github.com/512/simple_bank/dto"
	"github.com/512/simple_bank/service"
	"github.com/512/simple_bank/service/token"
	"github.com/512/simple_bank/ultis"
	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type authController struct {
	authService service.AuthService
	jwtService  token.JWTService
}

type loginUserResponse struct {
	AccessToken string  `json:"access_token"`
	User        db.User `json:"user"`
}

func NewAuthController(authService service.AuthService, jwtService token.JWTService) AuthController {
	return &authController{
		authService: authService,
		jwtService:  jwtService,
	}
}

func (controller *authController) Register(ctx *gin.Context) {
	var createUserDTO dto.CreateUserDTO
	errDTO := ctx.ShouldBind(&createUserDTO)
	if errDTO != nil {
		res := ultis.BuildErrorResponse("Failed to bind DTO", errDTO.Error(), ultis.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	userResult, err := controller.authService.CreateUser(createUserDTO)
	if err != nil {
		res := ultis.BuildErrorResponse("Failed to process create", err.Error(), ultis.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := ultis.BuildResponse(true, "", userResult)
	ctx.JSON(http.StatusOK, res)
}

func (controller *authController) Login(ctx *gin.Context) {
	var loginDTO dto.LoginDTO
	errDTO := ctx.ShouldBind(&loginDTO)
	if errDTO != nil {
		res := ultis.BuildErrorResponse("Failed to bind DTO", errDTO.Error(), ultis.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	user, err := controller.authService.VerifyAccount(loginDTO)
	if err != nil {
		res := ultis.BuildErrorResponse("Failed to process login", err.Error(), ultis.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	accessTokenDuration := os.Getenv("ACCESS_TOKEN_DURATION")
	duration, _ := time.ParseDuration(accessTokenDuration)
	accessToken, err := controller.jwtService.GenerateToken(
		user.Username,
		duration,
	)
	if err != nil {
		res := ultis.BuildErrorResponse("Failed to process generate access token", err.Error(), ultis.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := ultis.BuildResponse(true, "", loginUserResponse{
		AccessToken: accessToken,
		User:        user,
	})
	ctx.JSON(http.StatusOK, res)

}
