package controller

import (
	"net/http"

	"gestrym/src/core/login/app"
	structs_request "gestrym/src/core/login/domain/structs/request"

	"github.com/gin-gonic/gin"
)

type LoginController struct {
	loginService app.ILoginService
}

func NewLoginController(loginService app.ILoginService) *LoginController {
	return &LoginController{loginService: loginService}
}

// @Summary Iniciar sesión
// @Description Inicia sesión con email y contraseña. Devuelve access_token y refresh_token.
// @Tags Login
// @Accept json
// @Produce json
// @Param request body structs_request.LoginRequest true "Credenciales de acceso"
// @Success 200 {object} response.LoginResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /public/login [post]
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
