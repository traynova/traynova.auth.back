package controller

import (
	"net/http"

	"gestrym/src/common/utils"
	"gestrym/src/core/auth/app"
	structs_request "gestrym/src/core/auth/domain/structs/request"

	"github.com/gin-gonic/gin"
)

type AuthPublicController struct {
	authService app.IAuthService
	logger      utils.ILogger
}

func NewAuthPublicController(as app.IAuthService, logger utils.ILogger) *AuthPublicController {
	return &AuthPublicController{
		authService: as,
		logger:      logger,
	}
}

func (a *AuthPublicController) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		createUserRequest := &structs_request.RegisterRequest{}

		if err := c.ShouldBindJSON(createUserRequest); err != nil {
			a.logger.Error("Error while binding JSON", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		response, err := a.authService.RegisterUser(*createUserRequest, 0)
		if err != nil {
			a.logger.Error("Error while creating user", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"response": response})
	}
}

// @Summary Iniciar Sesión Tradicional
// @Description Realiza el login de un usuario registrado usando email y password, y devuelve un JWT.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body structs_auth.LoginRequest true "Credenciales de login"
// @Success 200 {object} map[string]interface{} "access_token: xxx, refresh_token: yyy"
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /public/auth/login [post]
// func (a *AuthPublicController) Login(ctx *gin.Context) {
// 	var req structs_auth.LoginRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Petición inválida"})
// 		return
// 	}

// 	jwtKey := viper.GetString("JWT_KEY")
// 	accToken, refToken, err := a.authService.Login(ctx.Request.Context(), req.Email, req.Password, []byte(jwtKey))
// 	if err != nil {
// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"access_token": accToken, "refresh_token": refToken})
// }

// @Summary Login con Google
// @Description Inicia sesión utilizando un id_token provisto por Google. Si el usuario no existe, lo registra como Cliente (1).
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body structs_auth.GoogleLoginRequest true "Token de Google"
// @Success 200 {object} map[string]interface{} "access_token: xxx, refresh_token: yyy"
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /public/auth/google [post]
// func (a *AuthPublicController) GoogleLogin(ctx *gin.Context) {
// 	var req structs_auth.GoogleLoginRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Petición inválida: falta id_token"})
// 		return
// 	}

// 	jwtKey := viper.GetString("JWT_KEY")
// 	accToken, refToken, err := a.authService.GoogleLogin(ctx.Request.Context(), req.IDToken, []byte(jwtKey))
// 	if err != nil {
// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"access_token": accToken, "refresh_token": refToken})
// }
