package middleware

import (
	"net/http"

	"github.com/512/simple_bank/service/token"
	"github.com/512/simple_bank/ultis"
	"github.com/gin-gonic/gin"
)

var (
	jwtService token.JWTService = token.NewJWTService()
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			response := ultis.BuildErrorResponse("Failed to process request", "No token found", nil)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		payload, err := jwtService.VerifyToken(authHeader)
		if err != nil {
			response := ultis.BuildErrorResponse("Token is not valid", err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return

		}
		ctx.Set("payload", payload)
		ctx.Next()
	}

}
