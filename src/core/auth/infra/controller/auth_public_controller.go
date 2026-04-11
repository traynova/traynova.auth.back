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

// @Summary Registrar Coach o Gimnasio o Cliente
// @Description Crea un nuevo usuario públicamente. Solo permite roles Cliente (1) o Coach (2) o Gym (3).
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body structs_request.RegisterRequest true "Datos del usuario"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /public/auth/register [post]
func (a *AuthPublicController) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		createUserRequest := &structs_request.RegisterRequest{}

		if err := c.ShouldBindJSON(&createUserRequest); err != nil {
			a.logger.Error("Error while binding JSON", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener el ID del usuario"})
			return
		}

		userId, ok := userIDInterface.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ID del usuario en formato incorrecto"})
			return
		}

		response, err := a.authService.RegisterUser(*createUserRequest, userId)
		if err != nil {
			a.logger.Error("Error while creating user", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"response": response})
	}
}

// @Summary Confirmar email de usuario
// @Description Activa al usuario cuando el token de confirmación es válido
// @Tags Auth
// @Accept json
// @Produce json
// @Param token query string true "Token de confirmación"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /public/auth/confirm [get]
func (a *AuthPublicController) ConfirmEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "token es requerido"})
			return
		}

		response, err := a.authService.ActivateUser(token)
		if err != nil {
			a.logger.Error("Error al confirmar email", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

// @Summary Solicitar recuperación de contraseña
// @Description Envía un enlace de recuperación al email del usuario
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body structs_request.PasswordRecoveryRequest true "Email del usuario"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /public/auth/password/recovery [post]
func (a *AuthPublicController) RequestPasswordRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req structs_request.PasswordRecoveryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := a.authService.RequestPasswordRecovery(req.Email); err != nil {
			a.logger.Error("Error al solicitar recuperación de contraseña", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Email de recuperación enviado"})
	}
}

// @Summary Restablecer contraseña
// @Description Cambia la contraseña usando el token de recuperación
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body structs_request.PasswordResetRequest true "Token y nueva contraseña"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /public/auth/password/reset [post]
func (a *AuthPublicController) ResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req structs_request.PasswordResetRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		response, err := a.authService.ResetPassword(req)
		if err != nil {
			a.logger.Error("Error al restablecer contraseña", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
	}
}
