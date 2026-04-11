package controller

import (
	"net/http"

	"gestrym/src/core/auth/login/app"
	structs_request "gestrym/src/core/auth/login/domain/structs/request"

	"github.com/gin-gonic/gin"
)

type LoginController struct {
	loginService app.ILoginService
}

func NewLoginController(loginService app.ILoginService) *LoginController {
	return &LoginController{loginService: loginService}
}

func (c *LoginController) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req structs_request.LoginRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		response, err := c.loginService.Login(req)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, response)
	}
}
